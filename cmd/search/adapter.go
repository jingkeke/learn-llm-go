package main

import (
	"context"

	"github.com/smhanov/laconic"
	"github.com/smhanov/llmhub"
)

// llmAdapter adapts an llmhub.Client to the laconic.LLMProvider interface.
type llmAdapter struct {
	client *llmhub.Client
}

func newLLMAdapter(client *llmhub.Client) *llmAdapter {
	return &llmAdapter{client: client}
}

func (a *llmAdapter) Generate(ctx context.Context, systemPrompt, userPrompt string) (laconic.LLMResponse, error) {
	msgs := []*llmhub.Message{}
	if systemPrompt != "" {
		msgs = append(msgs, llmhub.NewSystemMessage(llmhub.Text(systemPrompt)))
	}
	if userPrompt != "" {
		msgs = append(msgs, llmhub.NewUserMessage(llmhub.Text(userPrompt)))
	}

	resp, err := a.client.Generate(ctx, msgs)
	if err != nil {
		return laconic.LLMResponse{}, err
	}

	return laconic.LLMResponse{
		Text:      resp.Text(),
		Cost:      resp.Usage.Cost,
		Reasoning: resp.ReasoningText(),
	}, nil
}
