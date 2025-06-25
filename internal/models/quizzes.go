package models

import (
	"database/sql"
	"time"
)

type QuizModelInterface interface {
	Latest() ([]QuizMetadata, error)
}

type Quiz struct {
	ID          int
	Profile     Profile
	Title       string
	Description string
	Questions   []Question
	Published   *time.Time
	Unpublished *time.Time
	Created     time.Time
	Update      time.Time
	Deleted     *time.Time
}

type QuizMetadata struct {
	ID            int
	Profile       PublicProfile
	Title         string
	Description   string
	QuestionCount int
	Published     time.Time
}

type QuizModel struct {
	DB *sql.DB
}

func (m *QuizModel) Latest() ([]QuizMetadata, error) {
	var quizzes []QuizMetadata

	stmt := `SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, (SELECT count(id) FROM question WHERE quiz_id = q.id) AS question_count, q.published
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id
	WHERE q.published IS NOT NULL AND (q.unpublished IS NULL OR q.published > q.unpublished) AND q.deleted IS NULL
	ORDER BY q.published DESC
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return quizzes, nil
	}

	defer rows.Close()

	for rows.Next() {
		var p PublicProfile
		var q QuizMetadata

		err = rows.Scan(&q.ID, &p.ID, &p.FirstName, &p.LastName, &p.Deleted, &q.Title, &q.Description, &q.QuestionCount, &q.Published)
		if err != nil {
			return []QuizMetadata{}, err
		}

		q.Profile = p

		quizzes = append(quizzes, q)
	}

	return quizzes, nil
}
