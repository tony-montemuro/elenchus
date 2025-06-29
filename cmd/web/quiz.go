package main

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

type Quiz struct {
	Title       string     `json:"title" jsonschema:"The ideal name of the quiz, based on the user input"`
	Description string     `json:"description" jsonschema:"A description of what the quiz is trying to teach, between 140 and 280 characters"`
	Questions   []Question `json:"questions" jsonschema:"Up to 5 questions based on the input provided by the user -- can be less if user input is short"`
}

type Question struct {
	Content string   `json:"content" jsonschema:"The question that the user is requested to answer"`
	Answers []Answer `json:"answers" jsonschema:"Between 2 and 4 answers to the question, where exactly 1 is marked as correct"`
}

type Answer struct {
	Content string `json:"content" jsonschema:"One of N answers the user has to select from, where N = (2, 3, 4)"`
	Correct bool   `json:"correct" jsonschema:"True if the answer is correct, false otherwise. This flag should be true on exactly one answer per question."`
}

func generateSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var QuizResponseSchema = generateSchema[Quiz]()

func (app *application) generateQuiz(notes string, ctx context.Context) (Quiz, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "quiz",
		Description: openai.String("Quiz generated strictly on the notes provided by the user"),
		Schema:      QuizResponseSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := app.openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(notes),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		Model: openai.ChatModelGPT4o2024_08_06,
	})

	var quiz Quiz
	if err != nil {
		return quiz, err
	}

	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &quiz)
	return quiz, err
}
