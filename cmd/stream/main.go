package main

import (
	"context"
	"fmt"
	"os"

	"laconic-cli-tools/internal/config"

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
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	msgs := []*llmhub.Message{
		llmhub.NewUserMessage(llmhub.Text(cfg.Prompt)),
	}

	chunks, err := client.Stream(ctx, msgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting stream: %v\n", err)
		os.Exit(1)
	}

	for chunk := range chunks {
		if chunk.Err != nil {
			fmt.Fprintf(os.Stderr, "\nStream error: %v\n", chunk.Err)
			os.Exit(1)
		}
		if chunk.ReasoningDelta != "" {
			fmt.Fprintf(os.Stderr, "[reasoning] %s", chunk.ReasoningDelta)
		}
		fmt.Print(chunk.Delta)
		if chunk.Done {
			break
		}
	}
	fmt.Println()
}
