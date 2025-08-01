package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type QuizModelInterface interface {
	Latest() ([]QuizMetadata, error)
	GetQuizByID(int, *int) (QuizPublic, error)
	InsertQuiz(QuizJSONSchema, int, *sql.Tx) (int, error)
	GetPublishedQuizzesByProfile(*int) ([]QuizMetadata, error)
	GetUnpublishedQuizzesByProfile(*int) ([]QuizMetadata, error)
	UnpublishQuizByID(int, *sql.Tx) error
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
	Profile       ProfilePublic
	Title         string
	Description   string
	QuestionCount int
	PointsCount   int
	Published     *time.Time
	Unpublished   *time.Time
	Editable      bool
}

type QuizPublic struct {
	ID          int
	Profile     ProfilePublic
	Title       string
	Description string
	Questions   []QuestionPublic
	PointsCount int
	Published   *time.Time
	Editable    bool
}

type QuizJSONSchema struct {
	Title       string               `json:"title" jsonschema:"The ideal name of the quiz, based on the user input"`
	Description string               `json:"description" jsonschema:"A description of what the quiz is trying to teach, between 140 and 280 characters"`
	Questions   []QuestionJSONSchema `json:"questions" jsonschema:"Between 5 to 15 questions based on the input provided by the user -- can be less if user input is short"`
}

type QuizModel struct {
	DB *sql.DB
}

func (q *QuizMetadata) isNotPublished(unpublished *time.Time) bool {
	return isNotPublished(q.Published, unpublished)
}

func (q *QuizMetadata) isOwnedByProfile(profileID *int) bool {
	return isOwnedByProfile(profileID, q.Profile.ID)
}

func (q *QuizPublic) isNotPublished(unpublished *time.Time) bool {
	return isNotPublished(q.Published, unpublished)
}

func (q *QuizPublic) isOwnedByProfile(profileID *int) bool {
	return isOwnedByProfile(profileID, q.Profile.ID)
}

func (q *QuizPublic) Grade(answers QuestionAnswer) (AttemptPublic, error) {
	attempt := AttemptPublic{
		Quiz:    *q,
		Answers: answers,
	}

	for _, question := range q.Questions {
		correct, err := question.IsCorrect(answers[question.ID])
		if err != nil {
			return attempt, err
		}

		points := question.Points
		if correct {
			attempt.PointsEarned += points
		}
	}

	return attempt, nil
}

func (q *QuizPublic) IsSavable(profileID *int) bool {
	return profileID != nil && !q.Editable
}

func (m *QuizModel) Latest() ([]QuizMetadata, error) {
	var quizzes []QuizMetadata

	stmt := `SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, (SELECT count(id) FROM question WHERE quiz_id = q.id) AS question_count, (SELECT SUM(points) FROM question WHERE quiz_id = q.id) AS points_count, q.published
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id
	WHERE q.published IS NOT NULL AND (q.unpublished IS NULL OR q.published > q.unpublished) AND q.deleted IS NULL
	ORDER BY q.published DESC
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return quizzes, err
	}

	defer rows.Close()

	for rows.Next() {
		var p ProfilePublic
		var q QuizMetadata

		err = rows.Scan(&q.ID, &p.ID, &p.FirstName, &p.LastName, &p.Deleted, &q.Title, &q.Description, &q.QuestionCount, &q.PointsCount, &q.Published)
		if err != nil {
			return quizzes, err
		}

		q.Profile = p

		quizzes = append(quizzes, q)
	}

	return quizzes, nil
}

func (m *QuizModel) GetQuizByID(id int, profileID *int) (QuizPublic, error) {
	var quiz QuizPublic
	var profile ProfilePublic
	var deleted sql.NullTime
	var published sql.NullTime
	var unpublished *time.Time

	stmt := `SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, (SELECT SUM(points) FROM question WHERE quiz_id = q.id) AS points_count, q.published, q.unpublished
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id 
	WHERE q.id = ? AND q.deleted IS NULL`

	err := m.DB.QueryRow(stmt, id).Scan(&quiz.ID, &profile.ID, &profile.FirstName, &profile.LastName, &deleted, &quiz.Title, &quiz.Description, &quiz.PointsCount, &published, &unpublished)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quiz, ErrNoRecord
		}

		return quiz, err
	}

	if deleted.Valid {
		profile.Deleted = &deleted.Time
	} else {
		profile.Deleted = nil
	}

	quiz.Profile = profile

	if published.Valid {
		quiz.Published = &published.Time
	} else {
		quiz.Published = nil
	}

	if quiz.isNotPublished(unpublished) {
		// If the quiz is not published, and is also NOT owned by profile, then access is not allowed -- treat as missing record
		if !quiz.isOwnedByProfile(profileID) {
			return quiz, ErrNoRecord
		}

		// however, if the profile *does* own the unpublished quiz, they not only have access, but can *edit* the quiz
		quiz.Editable = true
	}

	return quiz, nil
}

func (m *QuizModel) InsertQuiz(quiz QuizJSONSchema, profileID int, tx *sql.Tx) (int, error) {
	stmt, err := tx.Prepare(`INSERT INTO quiz (profile_id, title, description, created, updated)
	VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(profileID, quiz.Title, quiz.Description)
	if err != nil {
		return 0, err
	}

	quizID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(quizID), nil
}

func (m *QuizModel) GetPublishedQuizzesByProfile(profileID *int) ([]QuizMetadata, error) {
	return m.getProfileQuizzes(profileID, true)
}

func (m *QuizModel) GetUnpublishedQuizzesByProfile(profileID *int) ([]QuizMetadata, error) {
	return m.getProfileQuizzes(profileID, false)
}

func (m *QuizModel) getProfileQuizzes(profileID *int, isPublished bool) ([]QuizMetadata, error) {
	var quizzes []QuizMetadata
	whereClause := `q.published IS NOT NULL AND q.published > COALESCE(q.unpublished, 0)`
	if !isPublished {
		whereClause = `NOT (` + whereClause + `)`
	}

	orderClause := `q.published DESC`
	if !isPublished {
		orderClause = `q.unpublished DESC, ` + orderClause
	}

	stmt := fmt.Sprintf(`SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, (SELECT count(id) FROM question WHERE quiz_id = q.id) AS question_count, (SELECT SUM(points) FROM question WHERE quiz_id = q.id) AS points_count, q.published, q.unpublished
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id
	WHERE %s AND p.id = ? 
	ORDER BY %s`, whereClause, orderClause)

	rows, err := m.DB.Query(stmt, *profileID)
	if err != nil {
		return quizzes, nil
	}

	defer rows.Close()

	for rows.Next() {
		var p ProfilePublic
		var q QuizMetadata

		err = rows.Scan(&q.ID, &p.ID, &p.FirstName, &p.LastName, &p.Deleted, &q.Title, &q.Description, &q.QuestionCount, &q.PointsCount, &q.Published, &q.Unpublished)
		if err != nil {
			return []QuizMetadata{}, err
		}

		q.Profile = p
		if q.isNotPublished(q.Unpublished) && q.isOwnedByProfile(profileID) {
			q.Editable = true
		}
		quizzes = append(quizzes, q)
	}

	return quizzes, nil
}

func (m *QuizModel) UpdateQuiz(quiz QuizPublic, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE quiz q
	SET q.title = ?, q.description = ?, q.updated = UTC_TIMESTAMP()
	WHERE q.id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(quiz.Title, quiz.Description, quiz.ID)
	return err
}

func (m *QuizModel) PublishQuizById(id int, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE quiz q
	SET q.published = UTC_TIMESTAMP(), q.updated = UTC_TIMESTAMP()
	WHERE q.id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (m *QuizModel) UnpublishQuizByID(quizID int, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE quiz q 
	SET q.unpublished = UTC_TIMESTAMP(), q.updated = UTC_TIMESTAMP()
	WHERE q.id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(quizID)
	return err
}

func isNotPublished(published, unpublished *time.Time) bool {
	return published == nil || (unpublished != nil && unpublished.After(*published))
}

func isOwnedByProfile(profileID *int, quizProfileId int) bool {
	return profileID != nil && quizProfileId == *profileID
}
