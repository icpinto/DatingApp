BEGIN;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS deactivated_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS user_lifecycle_outbox (
    event_id     UUID PRIMARY KEY,
    user_id      INT          NOT NULL,
    event_type   VARCHAR(20)  NOT NULL,
    payload      JSONB        NOT NULL DEFAULT '{}'::jsonb,
    processed    BOOLEAN      NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT user_lifecycle_outbox_event_type_chk CHECK (event_type IN ('deactivated', 'deleted'))
);

CREATE INDEX IF NOT EXISTS idx_user_lifecycle_outbox_processed ON user_lifecycle_outbox (processed, created_at);
CREATE INDEX IF NOT EXISTS idx_user_lifecycle_outbox_user ON user_lifecycle_outbox (user_id);

COMMIT;
