ALTER TABLE conversations
    ADD COLUMN user1_username VARCHAR(255),
    ADD COLUMN user2_username VARCHAR(255);

ALTER TABLE friend_requests
    ADD COLUMN sender_username VARCHAR(255),
    ADD COLUMN receiver_username VARCHAR(255);

UPDATE conversations c
SET user1_username = u1.username,
    user2_username = u2.username
FROM users u1, users u2
WHERE c.user1_id = u1.id AND c.user2_id = u2.id;

UPDATE friend_requests fr
SET sender_username = su.username,
    receiver_username = ru.username
FROM users su, users ru
WHERE fr.sender_id = su.id AND fr.receiver_id = ru.id;

ALTER TABLE conversations
    ALTER COLUMN user1_username SET NOT NULL,
    ALTER COLUMN user2_username SET NOT NULL;

ALTER TABLE friend_requests
    ALTER COLUMN sender_username SET NOT NULL,
    ALTER COLUMN receiver_username SET NOT NULL;
