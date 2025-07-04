package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tony-montemuro/elenchus/internal/models"
)

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

var QuizResponseSchema = generateSchema[models.QuizJSONSchema]()
var ErrGenerationRefusal = errors.New("Quiz creation unavailable: The submitted content doesn't meet our safety standards for educational content.")

func (app *application) generateQuiz(notes string, ctx context.Context) (models.QuizJSONSchema, error) {
	max_attempts := 3
	var quiz models.QuizJSONSchema
	var err error

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "quiz",
		Description: openai.String("Quiz generated strictly on the notes provided by the user. DO NOT EVER ENUMERATE THE QUESTIONS BY LETTER!"),
		Schema:      QuizResponseSchema,
		Strict:      openai.Bool(true),
	}

	for range max_attempts {
		var response *http.Response
		chat, err := app.openAIClient.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(notes),
				},
				ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
				},
				Model: openai.ChatModelO4Mini2025_04_16,
			},
			option.WithResponseInto(&response),
		)

		if err != nil {
			status := response.StatusCode

			// to be consistent with openAI sdk: https://github.com/openai/openai-go?tab=readme-ov-file#retries
			// for now, not worrying about "connection errors"
			if status == http.StatusRequestTimeout || status == http.StatusConflict || status == http.StatusTooManyRequests || status >= 500 {
				continue
			}

			return quiz, err
		}

		if chat.Choices[0].Message.Refusal != "" {
			app.logger.Warn("quiz generation refused", slog.String("message", chat.Choices[0].Message.Refusal))
			return quiz, ErrGenerationRefusal
		}

		// if we can successfully unmarshal the response, we can say the request was successful -- break out of loop
		err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &quiz)
		if err == nil {
			break
		}
	}

	return quiz, err
}
