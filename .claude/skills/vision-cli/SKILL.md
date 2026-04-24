---
name: vision-cli
description: Multi-modal image analysis tool using llmhub. Describe, analyze, or answer questions about images.
---

# Vision Analysis Skill

You are a computer vision expert. You can use the `vision-cli` tool to "see" images and answer questions about them. The tool supports local file paths, web URLs, or Base64 data URIs.

## Tool Execution
You can execute it via the compiled binary or directly with Go:
`go run ./cmd/vision --image <IMAGE_PATH> [flags] "Your Prompt Here"`

## Available Flags
- `--image`: **(Required)** Path to a local image file, a web URL (`http://...`), or a Base64 string.
- `--key` (or `OPENAI_API_KEY`): Required API key.
- `--base-url` (or `OPENAI_BASE_URL`): Optional. Set to the base URL of the provider.
- `--model` (or `OPENAI_MODEL`): Optional. Set the specific vision-capable model to use.

## Examples

To analyze a local image file:
```bash
go run ./cmd/vision --image "cat.jpg" "Describe this image in detail."
```

To extract text (OCR) from a web image:
```bash
go run ./cmd/vision --image "https://example.com/sign.png" "What text is written on this sign?"
```