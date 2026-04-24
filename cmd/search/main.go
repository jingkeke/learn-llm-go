package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"laconic-cli-tools/internal/config"

	"github.com/smhanov/laconic"
	"github.com/smhanov/laconic/search"
	"github.com/smhanov/llmhub"
	_ "github.com/smhanov/llmhub/providers/openai"
)

func main() {
	cfg, _ := config.ParseFlags()

	if err := cfg.ValidateTextPrompt(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var opts []llmhub.Option
	if cfg.BaseURL != "" {
		opts = append(opts, llmhub.WithBaseURL(cfg.BaseURL))
	}
	if cfg.Model != "" {
		opts = append(opts, llmhub.WithModel(cfg.Model))
	}

	client, err := llmhub.New("openai", cfg.Key, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating llmhub client: %v\n", err)
		os.Exit(1)
	}

	adapter := newLLMAdapter(client)

	// Add current date to the prompt to give context to the AI
	currentDate := time.Now().Format("2006-01-02")
	contextPrompt := fmt.Sprintf("Today's date is %s. %s", currentDate, cfg.Prompt)

	var searchProvider laconic.SearchProvider
	switch cfg.SearchProvider {
	case "brave":
		if cfg.SearchAPIKey == "" {
			fmt.Fprintln(os.Stderr, "Error: --search-api-key is required for Brave search")
			os.Exit(1)
		}
		searchProvider = search.NewBrave(cfg.SearchAPIKey)
	case "tavily":
		if cfg.SearchAPIKey == "" {
			fmt.Fprintln(os.Stderr, "Error: --search-api-key is required for Tavily search")
			os.Exit(1)
		}
		searchProvider = search.NewTavily(cfg.SearchAPIKey, "basic")
	default:
		searchProvider = search.NewDuckDuckGo()
	}

	agent := laconic.New(
		laconic.WithPlannerModel(adapter),
		laconic.WithSynthesizerModel(adapter),
		laconic.WithFinalizerModel(adapter),
		laconic.WithSearchProvider(searchProvider),
		laconic.WithStrategyName("scratchpad"), // use default strategy
		laconic.WithDebug(true),                // useful for CLI search tools
	)

	ctx := context.Background()
	result, err := agent.Answer(ctx, contextPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error answering query: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Final Answer ---")
	fmt.Println(result.Answer)
	if result.Cost > 0 {
		fmt.Printf("\nCost: $%.6f\n", result.Cost)
	}
}
