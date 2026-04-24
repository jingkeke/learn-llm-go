package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Key            string
	BaseURL        string
	Model          string
	Prompt         string
	Image          string // specific to vision
	SearchProvider string // specific to search CLI
	SearchAPIKey   string // specific to search CLI (for Brave/Tavily)
}

func ParseFlags() (*Config, []string) {
	cfg := &Config{}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.StringVar(&cfg.Key, "key", "", "API Key (or set OPENAI_API_KEY)")
	fs.StringVar(&cfg.BaseURL, "base-url", "", "Base URL for the API (or set OPENAI_BASE_URL)")
	fs.StringVar(&cfg.Model, "model", "", "Model to use (or set OPENAI_MODEL)")
	fs.StringVar(&cfg.Prompt, "prompt", "", "Text prompt to send")
	fs.StringVar(&cfg.Image, "image", "", "Path to image file or URL (for vision)")
	fs.StringVar(&cfg.SearchProvider, "search-provider", "duckduckgo", "Search provider: duckduckgo, brave, tavily (for search)")
	fs.StringVar(&cfg.SearchAPIKey, "search-api-key", "", "API Key for search provider (or set SEARCH_API_KEY)")

	fs.Parse(os.Args[1:])

	if cfg.Key == "" {
		cfg.Key = os.Getenv("OPENAI_API_KEY")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
	}

	if cfg.Model == "" {
		cfg.Model = os.Getenv("OPENAI_MODEL")
	}

	if cfg.SearchAPIKey == "" {
		cfg.SearchAPIKey = os.Getenv("SEARCH_API_KEY")
	}

	if os.Getenv("SEARCH_PROVIDER") != "" && cfg.SearchProvider == "duckduckgo" {
		cfg.SearchProvider = os.Getenv("SEARCH_PROVIDER")
	}

	// If prompt is not passed via flag, take the remaining arguments as prompt
	if cfg.Prompt == "" && fs.NArg() > 0 {
		cfg.Prompt = strings.Join(fs.Args(), " ")
	}

	return cfg, fs.Args()
}

func (c *Config) ValidateTextPrompt() error {
	if c.Key == "" {
		return fmt.Errorf("API key is required. Use --key or set OPENAI_API_KEY environment variable")
	}
	if c.Prompt == "" {
		return fmt.Errorf("prompt is required. Use --prompt or pass it as an argument")
	}
	return nil
}
