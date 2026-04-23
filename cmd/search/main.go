package main

import (
	"context"
	"fmt"
	"os"

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

	agent := laconic.New(
		laconic.WithPlannerModel(adapter),
		laconic.WithSynthesizerModel(adapter),
		laconic.WithFinalizerModel(adapter),
		laconic.WithSearchProvider(search.NewDuckDuckGo()),
		laconic.WithStrategyName("scratchpad"), // use default strategy
		laconic.WithDebug(true), // useful for CLI search tools
	)

	ctx := context.Background()
	result, err := agent.Answer(ctx, cfg.Prompt)
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
