package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"strconv"
	"strings"
)

func CallOpenAI(apiKey, original string, variants []string) (int, error) {
	prompt := fmt.Sprintf(
		"Select the address variant that best matches the original address based on semantic similarity. The original address is: '%s'. The variants are:\n",
		original,
	)
	for i, variant := range variants {
		prompt += fmt.Sprintf("%d. '%s'\n", i+1, variant)
	}
	prompt += "Respond only with the number of the most similar variant, in int64 format."

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return 0, errors.New("failed to get a successful response from OpenAI")
	}

	// Обрабатываем ответ
	responseContent := resp.Choices[0].Message.Content
	responseContent = strings.TrimSpace(responseContent)
	selectedVariant, err := strconv.Atoi(responseContent)
	if err != nil {
		return 0, errors.New("failed to convert response to integer")
	}

	if selectedVariant > len(variants) || selectedVariant < 1 {
		return 0, errors.New("result is out of range")
	}

	return selectedVariant - 1, nil
}
