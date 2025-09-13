ALTER TABLE friend_requests
    ALTER COLUMN conversation_id DROP NOT NULL,
    ALTER COLUMN conversation_id TYPE UUID USING NULL::uuid;

ALTER TABLE conversation_outbox
    ALTER COLUMN conversation_id DROP NOT NULL,
    ALTER COLUMN conversation_id TYPE UUID USING NULL::uuid;
