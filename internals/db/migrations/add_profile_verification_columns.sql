ALTER TABLE profiles
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR(32),
    ADD COLUMN IF NOT EXISTS contact_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS identity_verified BOOLEAN DEFAULT FALSE;

-- Ensure existing rows have explicit false values for the new verification flags.
UPDATE profiles
SET contact_verified = COALESCE(contact_verified, FALSE),
    identity_verified = COALESCE(identity_verified, FALSE)
WHERE contact_verified IS NULL OR identity_verified IS NULL;
