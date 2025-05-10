# Jarvis

This cli tool is used to quickly generate shell commands by using offline available llm (qwen2.5-coder:0.5b).

# Pre-requisites
1. Ollama
2. LLM Model (qwen2.5-coder:0.5b   --- Small which makes it very fast all well)
3. Go to compile or use pre-compiled binaries

# Use

1. Clone repo and have prequisites installed
2. go build
3. Place binary in bin folder
4. Type jarvis, press enter and then fill the prompt. It will give the command which is automatically copied to clipboard
5. Add an alias to your shell that can use `xclip -selection clipboard` to automatically output the command on terminal input
