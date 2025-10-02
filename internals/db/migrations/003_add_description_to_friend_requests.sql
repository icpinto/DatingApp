ALTER TABLE friend_requests
    ADD COLUMN description TEXT;

UPDATE friend_requests
SET description = ''
WHERE description IS NULL;

ALTER TABLE friend_requests
    ALTER COLUMN description SET DEFAULT '',
    ALTER COLUMN description SET NOT NULL;
