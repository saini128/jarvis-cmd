# Jarvis

This cli tool is used to quickly generate shell commands by using offline available llm (qwen2.5-coder:0.5b).

# Pre-requisites
1. Ollama
2. LLM Model (qwen2.5-coder:3b   --- Medium Sized Model with good accuracy)
3. Go to compile or use pre-compiled binaries

# Use

1. Clone repo and have prequisites installed
2. go build
3. Place binary in bin folder
4. Type jarvis, press enter and then fill the prompt. It will give the command which is automatically copied to clipboard
5. Add an alias to your shell that can use `xclip -selection clipboard` to automatically output the command on terminal input


# 🧪 Model Comparison for Shell Command Generation

**Test Environment**:
- **CPU**: AMD Ryzen 7 5800H
- **RAM**: 16 GB
- **GPU**: NVIDIA GTX 1650 4GB
- **Storage**: SSD
- **Use Case**: Bash command generation from natural language (Fedora 42, non-root context)

| Model Name              | Size  | 🕒 Speed        | 🎯 Accuracy                            | Notes                                                                                                                               |
|-------------------------|-------|-----------------|----------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| `qwen2.5-coder:0.5b`    | 0.5B  | ⚡ Very Fast     | ❌ Poor – Ignores strict formatting     | Often returns commands with markdown or explanations despite prompts                                                                 |
| `qwen2.5-coder:7b`      | 7B    | 🟡 Slow          | ✅✅ Very Good – Follows format well    | Accurate but noticeably slower on local 4GB GPU                                                                                      |
| `phi4-mini`             | 3.8B     | 🐌 Very Slow      | ❌ Poor                                | Slower than `qwen2.5-coder:7b` and less accurate; for prompt `i want to see all docker containers` it returned `sudo systemctl status docker && sudo lsof -i :2375` |
| `qwen2.5-coder:3b`      | 3B    | ⚡️ Fast         | ✅✅✅ Excellent                        | Best combination of speed and accuracy observed.                                                                                    |

> ✅ Recommendation: For strict Bash command output without explanations, try lighter **instruction-tuned models** like `phi2`, `phi3`, or `qwen1.5:1.8b`. They offer better speed without a significant drop in quality for this specific task. `qwen2.5-coder:3b` also shows promising results for a balance of speed and accuracy.

# Auto Paste Generated Command to Terminal Cursor
The following setup describes how to automatically paste a generated command to the terminal cursor in a Linux environment (Fedora 42), using xclip and a Fish shell function.

**Prerequisites:**

- `xclip`: This command-line utility is required to interact with the system's clipboard.  Ensure it is installed:

```bash
sudo dnf install xclip
```
**Setup in Fedora 42 (Fish Shell):**

The following Fish shell function, js, will execute a command (assumed to be jarvis in this case) and then insert its output into the current command line buffer.
```bash
function js
    jarvis #  Replace 'jarvis' with your command generation tool
    commandline -i (xclip -selection clipboard -o)
end
```