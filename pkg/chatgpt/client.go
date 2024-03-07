package chatgpt

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type GPTClient struct {
	APIKey string
	client *openai.Client
}

type SEOPromptResult struct {
	Title       string
	Description string
	Tags        []string
}

const PROMPT = `
You are an SEO expert making sure that a YouTube channel posting short-form content is performing well. The content is usually taken from a forum, where short stories are found and transformed into a script that is then read out to the user viewing the video. Your goal is to write a Title, a Description and tags that will make sure the video performs very well with the intent of monetizing views on the long-term. When provided with a json input, you generate an output in this format, and include absolutely no other text, not even an explanation of the format used. Your title, description and tags should avoid talking about anything other than the content provided at the end of this message. Here is the format to follow:
---
Title: [insert SEO-optimized video title here]
Description: [Insert SEO-optimized video description here]
Tags: [Insert SEO-Optimized video tags here]
---
Here is the content of the video you will be working with:

%s
`

func NewGPTClient(apiKey string) *GPTClient {
	return &GPTClient{
		APIKey: apiKey,
		client: openai.NewClient(apiKey),
	}
}

func (g *GPTClient) Prompt(prompt string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	answer := resp.Choices[0].Message.Content
	return answer, nil
}

func (g *GPTClient) PromptFromContent(content string) (SEOPromptResult, error) {
	fullPrompt := fmt.Sprintf(PROMPT, content)
	answer, err := g.Prompt(fullPrompt)
	if err != nil {
		return SEOPromptResult{}, fmt.Errorf("Error prompting GPT4: %v", err)
	}

	lines := strings.Split(answer, "\n")

	title := strings.TrimPrefix(lines[0], "Title: ")
	description := strings.TrimPrefix(lines[1], "Description: ")
	tags := strings.Split(strings.TrimPrefix(lines[2], "Tags: "), ",")

	return SEOPromptResult{
		Title:       title,
		Description: description,
		Tags:        tags,
	}, nil
}
