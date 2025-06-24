
-- +migrate Up
INSERT INTO question_type (name, default_points, created, updated)
VALUES
    ('multiple_choice', 1, NOW(), NOW()),
    ('free_response', 5, NOW(), NOW());

-- +migrate Down
DELETE FROM question_type
WHERE name IN ('multiple_choice', 'free_response');
