
-- +migrate Up
CREATE TABLE profile (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    deleted DATETIME NULL
);

ALTER TABLE profile ADD CONSTRAINT profile_uc_email UNIQUE(email);

CREATE TABLE quiz (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    profile_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1024) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    deleted DATETIME NULL,
    FOREIGN KEY(profile_id)
	REFERENCES profile(id)
);

CREATE TABLE question_type (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(25) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE TABLE question (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    quiz_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    content VARCHAR(2048) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    deleted DATETIME NULL,
    FOREIGN KEY(quiz_id)
	REFERENCES quiz(id),
    FOREIGN KEY(type_id)
	REFERENCES question_type(id)
);

CREATE INDEX idx_question_quiz_id ON question(quiz_id);
CREATE INDEX idx_question_type_id ON question(type_id);

CREATE TABLE answer (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    question_id INTEGER NOT NULL,
    content VARCHAR(2048) NOT NULL,
    correct BOOLEAN NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    deleted DATETIME NULL,
    FOREIGN KEY(question_id)
	REFERENCES question(id)
);

CREATE INDEX idx_answer_question_id ON answer(question_id);


-- +migrate Down
DROP TABLE question_type;
DROP INDEX idx_answer_question_id ON answer;
DROP TABLE answer;
DROP INDEX idx_question_quiz_id ON question;
DROP INDEX idx_question_type_id ON question;
DROP TABLE question;
DROP TABLE quiz;
DROP INDEX profile_uc_email ON profile;
DROP TABLE profile;
