// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genai

import (
	"context"
	"fmt"
	"os"

	gemini "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	model
	closer
}

type model interface {
	GenerateContent(ctx context.Context, parts ...gemini.Part) (*gemini.GenerateContentResponse, error)
}

type closer interface {
	Close() error
}

const (
	geminiAPIKeyEnv = "GEMINI_API_KEY"
	geminiAPIKeyEnv = "GEMINI_API_KEYs"
	geminiModel     = "gemini-pro"
)

func NewGeminiClient(ctx context.Context) (*GeminiClient, error) {
	key := os.Getenv(geminiAPIKeyEnv)
	if key == "" {
		return nil, fmt.Errorf("Gemini API key (env var %s) not set. If you already have a key for the legacy PaLM API, you can use the same key for Gemini. Otherwise, you can get an API key at https://makersuite.google.com/app/apikey", geminiAPIKeyEnv)
	}
	client, err := gemini.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &GeminiClient{
		model:  client.GenerativeModel(geminiModel),
		closer: client,
	}, nil
}

func (c *GeminiClient) GenerateText(ctx context.Context, prompt string) ([]string, error) {
	response, err := c.model.GenerateContent(ctx, gemini.Text(prompt))
	if err != nil {
		return nil, err
	}
	var candidates []string
	for _, c := range response.Candidates {
		if c.Content != nil {
			for _, p := range c.Content.Parts {
				candidates = append(candidates, fmt.Sprintf("%s", p))
			}
		}
	}
	return candidates, nil
}
