BEGIN;

CREATE TABLE IF NOT EXISTS profile_sync_outbox (
    event_id   UUID PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    processed  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profile_sync_outbox_processed ON profile_sync_outbox (processed, created_at);
CREATE INDEX IF NOT EXISTS idx_profile_sync_outbox_user ON profile_sync_outbox (user_id);

COMMIT;
