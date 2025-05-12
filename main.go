package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	// Check if Ollama is installed
	if !CheckOllamaInstalled() {
		fmt.Println("Ollama is not installed. Install it with:")
		fmt.Println("curl -fsSL https://ollama.com/install.sh | sh")
		return
	}

	// Check if Ollama API is reachable
	if !CheckOllamaAPIReachable() {
		fmt.Println("Ollama API is not reachable. Make sure Ollama is running.")
		fmt.Println("You can start it by running 'ollama serve' in a terminal.")
		return
	}

	// Check if required model is available
	modelName := "qwen2.5-coder:3b"
	if !CheckModelAvailable(modelName) {
		fmt.Println("Required model not found. Pull it with:")
		fmt.Printf("ollama pull %s\n", modelName)

		// Ask if the user wants to pull the model now
		fmt.Print("Do you want to pull the model now? (y/n): ")
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			fmt.Printf("Pulling %s model...\n", modelName)
			cmd := exec.Command("ollama", "pull", modelName)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error pulling model: %s\n", err)
				return
			}
			fmt.Println("Model successfully pulled!")
		} else {
			return
		}
	}

	var userPrompt string
	fmt.Println("Enter your prompt:")
	reader := bufio.NewReader(os.Stdin)
	userPrompt, _ = reader.ReadString('\n')

	command, err := GenerateCommand(modelName, userPrompt)
	if err != nil {
		fmt.Printf("Error generating command: %s\n", err)
		return
	}

	fmt.Print("🔧 Command: ")
	fmt.Println(command)
	clipboard.WriteAll(command)
	fmt.Println("Model Used: ", modelName)
	fmt.Println("✅ Command copied to clipboard. Paste it with Ctrl+Shift+V or Cmd+V.")

}

func CheckOllamaInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

// checkOllamaAPIReachable verifies if the Ollama API is reachable
func CheckOllamaAPIReachable() bool {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ModelInfo represents information about an Ollama model
type ModelInfo struct {
	Name string `json:"name"`
}

// ListModelsResponse represents the response from Ollama's list models API
type ListModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

// checkModelAvailable checks if the specified model is available
func CheckModelAvailable(modelName string) bool {

	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Simple check: if the model name appears in the output
	return strings.Contains(string(output), modelName)
}

// generateCommand sends a request to Ollama to generate a shell command
func GenerateCommand(modelName, userPrompt string) (string, error) {
	llm, err := ollama.New(ollama.WithModel("qwen2.5-coder:3b"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	prompt := fmt.Sprintf(`You are a strict command-line assistant for Fedora 42 user and not root.
		Do not include explanations, descriptions, or markdown formatting.
		Return ONLY the raw Bash command, without any quotes, comments, or formatting.
		NO triple backticks, NO "bash" labels, and NO explanations. Only one valid Bash command per response.
		Take care of commands that require sudo privileges.

		Human: %s
		Assistant:`, userPrompt)

	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.2),
	)
	if err != nil {
		return "", err
	}

	_ = completion
	validated, err := validateOutput(completion)
	if err != nil {
		return string(completion), err

	}
	return string(validated), nil
}
func validateOutput(output string) (string, error) {
	// Trim whitespace
	output = strings.TrimSpace(output)

	// Remove any leading/trailing markdown
	output = strings.TrimPrefix(output, "```bash")
	output = strings.TrimPrefix(output, "```")
	output = strings.TrimSuffix(output, "```")
	output = strings.TrimSpace(output)

	// Fail if output contains markdown or comments
	if strings.Contains(output, "```") || strings.Contains(output, "#") {
		return "", errors.New("output contains markdown or comments")
	}

	// Fail if multiple lines (can adjust if multiline commands are valid for you)
	if strings.Count(output, "\n") > 0 {
		return "", errors.New("output contains multiple lines")
	}

	// Optional: Regex for a basic shell command structure (very loose)
	commandRegex := regexp.MustCompile(`^[a-zA-Z0-9/\-_.:]+(?:\s+[^\s]+)*$`)
	if !commandRegex.MatchString(output) {
		return "", errors.New("output does not look like a valid shell command")
	}

	return output, nil
}
