-- 2. Profiles
CREATE TABLE profiles (
    id            SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users(id) ON DELETE CASCADE,
    bio           TEXT,
    gender        VARCHAR(20),
    date_of_birth DATE,
    location      VARCHAR(255),
    interests     TEXT[],
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_user_id UNIQUE (user_id)
);

ALTER TABLE profiles
    ADD COLUMN profile_image VARCHAR(255);


BEGIN;

-- =========================================================
-- 1) ENUM TYPES (create only if missing)
-- =========================================================
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'civil_status_type') THEN
    CREATE TYPE civil_status_type AS ENUM ('single','married','divorced','widowed','separated');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dietary_pref_type') THEN
    CREATE TYPE dietary_pref_type AS ENUM ('veg','non_veg','vegan','eggetarian','other');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'habit_freq_type') THEN
    CREATE TYPE habit_freq_type AS ENUM ('no','occasional','yes');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'education_level_type') THEN
    CREATE TYPE education_level_type AS ENUM ('secondary','diploma','bachelor','master','phd','professional','other');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'employment_status_type') THEN
    CREATE TYPE employment_status_type AS ENUM ('student','employed','self_employed','unemployed','retired','other');
  END IF;
END $$;

-- =========================================================
-- 2) OPTIONAL: preserve legacy free-text location
-- =========================================================
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name='profiles' AND column_name='location'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name='profiles' AND column_name='location_legacy'
  ) THEN
    ALTER TABLE profiles RENAME COLUMN location TO location_legacy;
    COMMENT ON COLUMN profiles.location_legacy IS 'Deprecated free-text location retained for backward compatibility.';
  END IF;
END $$;

-- =========================================================
-- 3) ADD/ENSURE COLUMNS (only if missing)
-- =========================================================
DO $$
BEGIN
  -- Personal info
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='civil_status') THEN
    ALTER TABLE profiles ADD COLUMN civil_status civil_status_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='religion') THEN
    ALTER TABLE profiles ADD COLUMN religion TEXT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='religion_detail') THEN
    ALTER TABLE profiles ADD COLUMN religion_detail TEXT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='caste') THEN
    ALTER TABLE profiles ADD COLUMN caste TEXT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='height_cm') THEN
    ALTER TABLE profiles ADD COLUMN height_cm SMALLINT CHECK (height_cm BETWEEN 80 AND 250);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='weight_kg') THEN
    ALTER TABLE profiles ADD COLUMN weight_kg SMALLINT CHECK (weight_kg BETWEEN 30 AND 300);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='dietary_preference') THEN
    ALTER TABLE profiles ADD COLUMN dietary_preference dietary_pref_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='smoking') THEN
    ALTER TABLE profiles ADD COLUMN smoking habit_freq_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='alcohol') THEN
    ALTER TABLE profiles ADD COLUMN alcohol habit_freq_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='languages') THEN
    ALTER TABLE profiles ADD COLUMN languages TEXT[];
  END IF;

  -- Residency (Sri Lanka oriented)
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='country_code') THEN
    ALTER TABLE profiles ADD COLUMN country_code CHAR(2) DEFAULT 'LK';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='province') THEN
    ALTER TABLE profiles ADD COLUMN province VARCHAR(100);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='district') THEN
    ALTER TABLE profiles ADD COLUMN district VARCHAR(100);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='city') THEN
    ALTER TABLE profiles ADD COLUMN city VARCHAR(100);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='postal_code') THEN
    ALTER TABLE profiles ADD COLUMN postal_code VARCHAR(20);
  END IF;

  -- Education & work (no salary fields)
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='highest_education') THEN
    ALTER TABLE profiles ADD COLUMN highest_education education_level_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='field_of_study') THEN
    ALTER TABLE profiles ADD COLUMN field_of_study VARCHAR(255);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='institution') THEN
    ALTER TABLE profiles ADD COLUMN institution VARCHAR(255);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='employment_status') THEN
    ALTER TABLE profiles ADD COLUMN employment_status employment_status_type;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='occupation') THEN
    ALTER TABLE profiles ADD COLUMN occupation VARCHAR(255);
  END IF;

  -- Family
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='father_occupation') THEN
    ALTER TABLE profiles ADD COLUMN father_occupation VARCHAR(255);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='mother_occupation') THEN
    ALTER TABLE profiles ADD COLUMN mother_occupation VARCHAR(255);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='siblings_count') THEN
    ALTER TABLE profiles ADD COLUMN siblings_count SMALLINT CHECK (siblings_count >= 0);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='siblings') THEN
    ALTER TABLE profiles ADD COLUMN siblings JSONB;
  END IF;

  -- Horoscope
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='horoscope_available') THEN
    ALTER TABLE profiles ADD COLUMN horoscope_available BOOLEAN DEFAULT FALSE;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='birth_time') THEN
    ALTER TABLE profiles ADD COLUMN birth_time TIME;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='birth_place') THEN
    ALTER TABLE profiles ADD COLUMN birth_place VARCHAR(255);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='sinhala_raasi') THEN
    ALTER TABLE profiles ADD COLUMN sinhala_raasi VARCHAR(50);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='nakshatra') THEN
    ALTER TABLE profiles ADD COLUMN nakshatra VARCHAR(50);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='horoscope') THEN
    ALTER TABLE profiles ADD COLUMN horoscope JSONB;
  END IF;

  -- Images (store URLs/keys; keep legacy profile_image if present)
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='profile_image_url') THEN
    ALTER TABLE profiles ADD COLUMN profile_image_url VARCHAR(512);
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='profile_image_thumb_url') THEN
    ALTER TABLE profiles ADD COLUMN profile_image_thumb_url VARCHAR(512);
  END IF;

  -- Admin/meta
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='verified') THEN
    ALTER TABLE profiles ADD COLUMN verified BOOLEAN DEFAULT FALSE;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='moderation_status') THEN
    ALTER TABLE profiles ADD COLUMN moderation_status VARCHAR(20) DEFAULT 'clean';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='last_active_at') THEN
    ALTER TABLE profiles ADD COLUMN last_active_at TIMESTAMP;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='metadata') THEN
    ALTER TABLE profiles ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
  END IF;
END $$;

-- If legacy 'profile_image' exists, mark it deprecated (comment only)
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='profiles' AND column_name='profile_image') THEN
    COMMENT ON COLUMN profiles.profile_image IS 'Deprecated; prefer profile_image_url stored in object storage (S3/R2).';
  END IF;
END $$;

-- =========================================================
-- 4) CONSTRAINTS & BACKFILLS (safe)
-- =========================================================

-- Backfill country_code to 'LK' where NULL, then set NOT NULL
UPDATE profiles SET country_code = 'LK' WHERE country_code IS NULL;

ALTER TABLE profiles
  ALTER COLUMN country_code SET NOT NULL;

-- Country code format constraint (NOT VALID to avoid blocking on bad legacy data)
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_country_code_format_chk;
ALTER TABLE profiles
  ADD CONSTRAINT profiles_country_code_format_chk
  CHECK (country_code ~ '^[A-Z]{2}$')
  NOT VALID;

-- Date of birth must be in the past (NOT VALID to avoid blocking on bad legacy data)
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_dob_past_chk;
ALTER TABLE profiles
  ADD CONSTRAINT profiles_dob_past_chk
  CHECK (date_of_birth IS NULL OR date_of_birth <= CURRENT_DATE)
  NOT VALID;

-- =========================================================
-- 5) INDEXES (create if missing)
-- =========================================================
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_interests_gin') THEN
    CREATE INDEX idx_profiles_interests_gin ON profiles USING GIN (interests);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_languages_gin') THEN
    CREATE INDEX idx_profiles_languages_gin ON profiles USING GIN (languages);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_metadata_gin') THEN
    CREATE INDEX idx_profiles_metadata_gin ON profiles USING GIN (metadata);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_horoscope_gin') THEN
    CREATE INDEX idx_profiles_horoscope_gin ON profiles USING GIN (horoscope);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_residency') THEN
    CREATE INDEX idx_profiles_residency ON profiles (country_code, province, district, city);
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='idx_profiles_occupation') THEN
    CREATE INDEX idx_profiles_occupation ON profiles (occupation);
  END IF;
END $$;

-- =========================================================
-- 6) updated_at AUTO-UPDATE TRIGGER
-- =========================================================
CREATE OR REPLACE FUNCTION set_profiles_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END $$;

DROP TRIGGER IF EXISTS trg_profiles_set_updated_at ON profiles;
CREATE TRIGGER trg_profiles_set_updated_at
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION set_profiles_updated_at();

COMMIT;

-- (Optional, run later after cleaning data)
-- ALTER TABLE profiles VALIDATE CONSTRAINT profiles_country_code_format_chk;
-- ALTER TABLE profiles VALIDATE CONSTRAINT profiles_dob_past_chk;
