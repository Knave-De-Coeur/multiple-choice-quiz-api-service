INSERT INTO games (created_at, updated_at, deleted_at, score) VALUES
                                                                     (now(), now(), NULL, 5);

INSERT INTO questions (created_at, updated_at, deleted_at, game_id, description) VALUES
(NOW(), NOW(), NULL, 1, 'What is the result of 240 / 12?'),
(NOW(), NOW(), NULL, 1, 'Who was the first man on the moon?'),
(NOW(), NOW(), NULL, 1, 'Who wrote the hit single, Yellow Submarine?'),
(NOW(), NOW(), NULL, 1, 'In what year did Malta gain it\'s independence?'),
(NOW(), NOW(), NULL, 1, 'When was Go launched?');
