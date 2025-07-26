package main

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

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
var ErrNotEditable = errors.New("Quiz edit unavailable: user attempted to edit a quiz that cannot be edited")

func (app *application) generateQuizByForm(form createForm, ctx context.Context) (models.QuizJSONSchema, error) {
	var message openai.ChatCompletionMessageParamUnion

	if form.Type == "text" {
		message = openai.UserMessage(form.Text)
	} else {
		message = openai.UserMessage(
			[]openai.ChatCompletionContentPartUnionParam{
				{
					OfFile: &openai.ChatCompletionContentPartFileParam{
						File: openai.ChatCompletionContentPartFileFileParam{
							Filename: openai.String("input.pdf"),
							FileData: openai.String("data:application/pdf;base64," + string(form.File)),
						},
					},
				},
			},
		)
	}

	return app.generateQuiz(message, ctx)
}

func (app *application) generateQuiz(message openai.ChatCompletionMessageParamUnion, ctx context.Context) (models.QuizJSONSchema, error) {
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
					message,
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

func (app *application) getEditableQuizById(quizID int, r *http.Request) (models.QuizPublic, error) {
	quiz, err := app.getQuizByID(quizID, r)
	if err != nil {
		return quiz, err
	}

	if !quiz.Editable {
		return quiz, ErrNotEditable
	}

	return quiz, nil
}

func (app *application) getQuizByID(quizID int, r *http.Request) (models.QuizPublic, error) {
	profileID, _ := app.getProfileID(r)
	quiz, err := app.quizzesService.GetQuizByID(quizID, profileID)
	return quiz, err
}

func (app *application) buildNewQuizPublic(oldQuiz models.QuizPublic, editForm editForm) (models.QuizPublic, error) {
	quiz := models.QuizPublic{
		ID:          oldQuiz.ID,
		Profile:     oldQuiz.Profile,
		Title:       editForm.Title,
		Description: editForm.Description,
		Questions:   []models.QuestionPublic{},
		Published:   nil,
		Editable:    true,
	}

	mcType, err := app.questionTypes.GetMultipleChoice()
	if err != nil {
		return quiz, err
	}

	for _, id := range sortedKeys(editForm.Questions) {
		q := editForm.Questions[id]
		question := models.QuestionPublic{
			ID:      id,
			Type:    mcType,
			Content: q.Content,
			Points:  q.Points,
			Answers: []models.AnswerPublic{},
		}

		correctAnswerID := q.Correct
		answers := slices.SortedStableFunc(slices.Values(q.Answers), func(x, y answerEdit) int {
			return cmp.Compare(x.ID, y.ID)
		})

		for _, answer := range answers {
			question.Answers = append(question.Answers, models.AnswerPublic{
				ID:      answer.ID,
				Content: answer.Content,
				Correct: correctAnswerID == answer.ID,
			})
		}

		quiz.Questions = append(quiz.Questions, question)

	}

	return quiz, nil
}

func (app *application) getQuizIDParam(w http.ResponseWriter, r *http.Request) (int, error) {
	quizID, err := strconv.Atoi(r.PathValue("quizID"))
	if err != nil {
		if r.Method == "POST" {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.redirectNotFound(w, r, "user attempted to access a quiz that does not exist", err)
		}

		return 0, err
	}

	return quizID, nil
}
