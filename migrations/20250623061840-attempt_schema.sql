
-- +migrate Up
ALTER TABLE quiz
ADD published DATETIME;

ALTER TABLE quiz
ADD unpublished DATETIME;

ALTER TABLE question
ADD points INTEGER UNSIGNED NOT NULL;

ALTER TABLE question_type
ADD default_points INTEGER UNSIGNED NOT NULL;

CREATE TABLE attempt (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    profile_id INTEGER NOT NULL,
    quiz_id INTEGER NOT NULL,
    points_earned INTEGER UNSIGNED NOT NULL,
    created DATETIME NOT NULL,
    deleted DATETIME 
);

ALTER TABLE attempt 
ADD CONSTRAINT fk_attempt_profile
FOREIGN KEY (profile_id) REFERENCES profile(id);

ALTER TABLE attempt
ADD CONSTRAINT fk_attempt_quiz
FOREIGN KEY (quiz_id) REFERENCES quiz(id);

CREATE INDEX idx_attempt_profile_quiz ON attempt(profile_id, quiz_id);

CREATE TABLE multiple_choice_attempt (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    attempt_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer_id INTEGER NOT NULL,
    created DATETIME NOT NULL,
    deleted DATETIME
);

ALTER TABLE multiple_choice_attempt
ADD CONSTRAINT fk_multiple_choice_attempt_attempt
FOREIGN KEY (attempt_id) REFERENCES attempt(id);

ALTER TABLE multiple_choice_attempt
ADD CONSTRAINT fk_multiple_choice_attempt_question
FOREIGN KEY (question_id) REFERENCES question(id);

ALTER TABLE multiple_choice_attempt
ADD CONSTRAINT fk_multiple_choice_attempt_answer
FOREIGN KEY (answer_id) REFERENCES answer(id);

CREATE INDEX idx_multiple_choice_attempt_attempt_id ON multiple_choice_attempt(attempt_id);

CREATE TABLE free_response_attempt (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    attempt_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer VARCHAR(4096) NOT NULL,
    points INTEGER UNSIGNED NOT NULL,
    feedback VARCHAR(2048) NOT NULL,
    created DATETIME NOT NULL,
    deleted DATETIME
);

ALTER TABLE free_response_attempt 
ADD CONSTRAINT fk_free_response_attempt_attempt
FOREIGN KEY (attempt_id) REFERENCES attempt(id);

ALTER TABLE free_response_attempt 
ADD CONSTRAINT fk_free_response_attempt_question
FOREIGN KEY (question_id) REFERENCES question(id);

CREATE INDEX idx_free_response_attempt_attempt_id ON free_response_attempt(attempt_id);

-- +migrate Down
ALTER TABLE free_response_attempt DROP FOREIGN KEY fk_free_response_attempt_attempt;
ALTER TABLE free_response_attempt DROP FOREIGN KEY fk_free_response_attempt_question;
DROP INDEX idx_free_response_attempt_attempt_id ON free_response_attempt;
DROP TABLE free_response_attempt;
ALTER TABLE multiple_choice_attempt DROP FOREIGN KEY fk_multiple_choice_attempt_attempt;
ALTER TABLE multiple_choice_attempt DROP FOREIGN KEY fk_multiple_choice_attempt_question;
ALTER TABLE multiple_choice_attempt DROP FOREIGN KEY fk_multiple_choice_attempt_answer;
DROP INDEX idx_multiple_choice_attempt_attempt_id ON multiple_choice_attempt;
DROP TABLE multiple_choice_attempt;
ALTER TABLE attempt DROP FOREIGN KEY fk_attempt_profile;
ALTER TABLE attempt DROP FOREIGN KEY fk_attempt_quiz;
DROP INDEX idx_attempt_profile_quiz ON attempt;
DROP TABLE attempt;
ALTER TABLE question_type DROP COLUMN default_points;
ALTER TABLE question DROP COLUMN points;
ALTER TABLE quiz DROP COLUMN published;
ALTER TABLE quiz DROP COLUMN unpublished;
