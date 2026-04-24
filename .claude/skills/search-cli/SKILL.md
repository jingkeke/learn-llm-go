---
name: search-cli
description: Agentic web search tool using laconic. Finds grounded, up-to-date information on the web.
---

# Agentic Web Search Skill

You are a researcher that leverages the `search-cli` tool. When users ask questions requiring up-to-date, real-world, or highly specific information that an LLM might hallucinate or not know, use this tool. It uses the `laconic` agent to search the web (e.g. via DuckDuckGo) and synthesize an answer grounded in search facts.

## Tool Execution
You can execute it via the compiled binary or directly with Go:
`go run ./cmd/search [flags] "Your Question Here"`

## Available Flags
- `--key` (or `OPENAI_API_KEY`): Required API key for the LLM.
- `--base-url` (or `OPENAI_BASE_URL`): Optional API base URL.
- `--model` (or `OPENAI_MODEL`): Optional model to use (default depends on provider).
- `--search-provider` (or `SEARCH_PROVIDER`): Web search provider (`duckduckgo`, `brave`, `tavily`). Defaults to `duckduckgo`.
- `--search-api-key` (or `SEARCH_API_KEY`): Required if using Brave or Tavily.

## Usage Strategy
Because `search-cli` runs an iterative agent loop that calls the LLM multiple times to search and synthesize, it will take several seconds to complete. The output will include the final answer grounded in facts.

## Examples

To ask about real-time events or niche facts:
```bash
go run ./cmd/search "What is the current stock price of Apple?"
```

To configure a different search provider:
```bash
go run ./cmd/search --search-provider "tavily" --search-api-key "tvly-..." "Latest news on AI models"
```