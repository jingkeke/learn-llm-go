package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"laconic-cli-tools/internal/config"

	"github.com/smhanov/llmhub"
	_ "github.com/smhanov/llmhub/providers/openai"
)

func resolveImage(path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "data:") {
		return path, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(data)
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func main() {
	cfg, _ := config.ParseFlags()

	if err := cfg.ValidateTextPrompt(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg.Image == "" {
		fmt.Fprintln(os.Stderr, "Error: --image is required for vision CLI")
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

	imgURL, err := resolveImage(cfg.Image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing image %s: %v\n", cfg.Image, err)
		os.Exit(1)
	}

	ctx := context.Background()
	msgs := []*llmhub.Message{
		llmhub.NewUserMessage(
			llmhub.Text(cfg.Prompt),
			llmhub.Image(imgURL),
		),
	}

	resp, err := client.Generate(ctx, msgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resp.Text())
}
