BEGIN;

ALTER TABLE user_lifecycle_outbox
    DROP CONSTRAINT IF EXISTS user_lifecycle_outbox_event_type_chk;

ALTER TABLE user_lifecycle_outbox
    ADD CONSTRAINT user_lifecycle_outbox_event_type_chk CHECK (event_type IN ('deactivated', 'deleted', 'reactivated'));

COMMIT;
