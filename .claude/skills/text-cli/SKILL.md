---
name: text-cli
description: Standard text generation tool using llmhub. Send a prompt to an LLM and get a response.
---

# Text Generation Skill

You are a helper that can generate text using the local `text-cli` tool.

## Tool Execution
Use this tool to ask questions or generate content using external models.
You can execute it via the compiled binary or directly with Go:
`go run ./cmd/text [flags] "Your Prompt Here"`

## Available Flags (or Environment Variables)
- `--key` (or `OPENAI_API_KEY`): Required API key.
- `--base-url` (or `OPENAI_BASE_URL`): Optional. Set to the base URL of the provider (e.g. `http://localhost:11434/v1` for Ollama).
- `--model` (or `OPENAI_MODEL`): Optional. Set the specific model to use.

## Examples

To generate a simple response:
```bash
go run ./cmd/text "Explain quantum physics in simple terms"
```

To use a local Ollama server:
```bash
go run ./cmd/text --base-url "http://localhost:11434/v1" --model "llama3" "What is Go?"
```