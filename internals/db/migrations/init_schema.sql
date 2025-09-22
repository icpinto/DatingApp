-- Schema initialization script for the DatingApp service
-- Drops existing objects (if any) and recreates the full schema expected by the backend.

BEGIN;

-- ========================================================
-- 1. Tear down existing dependent objects (if present)
-- ========================================================
DROP TRIGGER IF EXISTS trg_profiles_set_updated_at ON profiles;
DROP FUNCTION IF EXISTS set_profiles_updated_at();

DROP TABLE IF EXISTS conversation_outbox CASCADE;
DROP TABLE IF EXISTS user_answers CASCADE;
DROP TABLE IF EXISTS questions CASCADE;
DROP TABLE IF EXISTS friend_requests CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS friend_request_status_type;
DROP TYPE IF EXISTS employment_status_type;
DROP TYPE IF EXISTS education_level_type;
DROP TYPE IF EXISTS habit_freq_type;
DROP TYPE IF EXISTS dietary_pref_type;
DROP TYPE IF EXISTS civil_status_type;

-- ========================================================
-- 2. Enumerated types used across the schema
-- ========================================================
CREATE TYPE civil_status_type AS ENUM (
    'single', 'married', 'divorced', 'widowed', 'separated'
);

CREATE TYPE dietary_pref_type AS ENUM (
    'veg', 'non_veg', 'vegan', 'eggetarian', 'other'
);

CREATE TYPE habit_freq_type AS ENUM (
    'no', 'occasional', 'yes'
);

CREATE TYPE education_level_type AS ENUM (
    'secondary', 'diploma', 'bachelor', 'master', 'phd', 'professional', 'other'
);

CREATE TYPE employment_status_type AS ENUM (
    'student', 'employed', 'self_employed', 'unemployed', 'retired', 'other'
);

CREATE TYPE friend_request_status_type AS ENUM (
    'pending', 'accepted', 'rejected'
);

-- ========================================================
-- 3. Core user data
-- ========================================================
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50)  NOT NULL UNIQUE,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ========================================================
-- 4. Extended profile data
-- ========================================================
CREATE TABLE profiles (
    id                         SERIAL PRIMARY KEY,
    user_id                    INT                NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio                        TEXT,
    gender                     VARCHAR(20),
    date_of_birth              DATE,
    location_legacy            VARCHAR(255),
    interests                  TEXT[]             NOT NULL DEFAULT ARRAY[]::text[],
    civil_status               civil_status_type,
    religion                   TEXT,
    religion_detail            TEXT,
    caste                      TEXT,
    height_cm                  SMALLINT CHECK (height_cm BETWEEN 80 AND 250),
    weight_kg                  SMALLINT CHECK (weight_kg BETWEEN 30 AND 300),
    dietary_preference         dietary_pref_type,
    smoking                    habit_freq_type,
    alcohol                    habit_freq_type,
    languages                  TEXT[]             NOT NULL DEFAULT ARRAY[]::text[],
    phone_number               VARCHAR(25),
    contact_verified           BOOLEAN            NOT NULL DEFAULT FALSE,
    identity_verified          BOOLEAN            NOT NULL DEFAULT FALSE,
    country_code               CHAR(2)            NOT NULL DEFAULT 'LK',
    province                   VARCHAR(100),
    district                   VARCHAR(100),
    city                       VARCHAR(100),
    postal_code                VARCHAR(20),
    highest_education          education_level_type,
    field_of_study             VARCHAR(255),
    institution                VARCHAR(255),
    employment_status          employment_status_type,
    occupation                 VARCHAR(255),
    father_occupation          VARCHAR(255),
    mother_occupation          VARCHAR(255),
    siblings_count             SMALLINT CHECK (siblings_count >= 0),
    siblings                   JSONB              NOT NULL DEFAULT '{}'::jsonb,
    horoscope_available        BOOLEAN            NOT NULL DEFAULT FALSE,
    birth_time                 TIME,
    birth_place                VARCHAR(255),
    sinhala_raasi              VARCHAR(50),
    nakshatra                  VARCHAR(50),
    horoscope                  JSONB              NOT NULL DEFAULT '{}'::jsonb,
    profile_image_url          VARCHAR(512),
    profile_image_thumb_url    VARCHAR(512),
    verified                   BOOLEAN            NOT NULL DEFAULT FALSE,
    moderation_status          VARCHAR(20)        NOT NULL DEFAULT 'clean',
    last_active_at             TIMESTAMP,
    metadata                   JSONB              NOT NULL DEFAULT '{}'::jsonb,
    created_at                 TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    CONSTRAINT profiles_country_code_format_chk CHECK (country_code ~ '^[A-Z]{2}$'),
    CONSTRAINT profiles_dob_past_chk CHECK (date_of_birth IS NULL OR date_of_birth <= CURRENT_DATE)
);

-- Maintain updated_at automatically on profile updates.
CREATE OR REPLACE FUNCTION set_profiles_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_profiles_set_updated_at
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION set_profiles_updated_at();

-- Helpful indexes for profile queries
CREATE INDEX idx_profiles_interests_gin ON profiles USING GIN (interests);
CREATE INDEX idx_profiles_languages_gin ON profiles USING GIN (languages);
CREATE INDEX idx_profiles_metadata_gin ON profiles USING GIN (metadata);
CREATE INDEX idx_profiles_horoscope_gin ON profiles USING GIN (horoscope);
CREATE INDEX idx_profiles_residency ON profiles (country_code, province, district, city);
CREATE INDEX idx_profiles_occupation ON profiles (occupation);

-- ========================================================
-- 5. Friend requests & conversation outbox
-- ========================================================
CREATE TABLE friend_requests (
    id                  SERIAL PRIMARY KEY,
    sender_id           INT                        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_username     VARCHAR(255)               NOT NULL,
    receiver_id         INT                        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_username   VARCHAR(255)               NOT NULL,
    status              friend_request_status_type NOT NULL DEFAULT 'pending',
    description         TEXT                       NOT NULL DEFAULT '',
    conversation_id     UUID,
    created_at          TIMESTAMPTZ                NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ                NOT NULL DEFAULT NOW(),
    CONSTRAINT friend_requests_sender_receiver_chk CHECK (sender_id <> receiver_id),
    CONSTRAINT friend_requests_sender_receiver_key UNIQUE (sender_id, receiver_id)
);

CREATE INDEX idx_friend_requests_receiver_status ON friend_requests (receiver_id, status);
CREATE INDEX idx_friend_requests_sender ON friend_requests (sender_id);

CREATE TABLE conversation_outbox (
    event_id        UUID PRIMARY KEY,
    user1_id        INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id        INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id UUID,
    processed       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT conversation_outbox_users_chk CHECK (user1_id <> user2_id)
);

CREATE INDEX idx_conversation_outbox_processed ON conversation_outbox (processed, created_at);

-- ========================================================
-- 6. Questionnaire data
-- ========================================================
CREATE TABLE questions (
    id             SERIAL PRIMARY KEY,
    question_text  TEXT        NOT NULL,
    question_type  VARCHAR(50) NOT NULL CHECK (question_type IN ('multiple_choice', 'scale', 'open_text')),
    options        TEXT[]      NOT NULL DEFAULT ARRAY[]::text[],
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_answers (
    id           SERIAL PRIMARY KEY,
    user_id      INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id  INT         NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer_text  TEXT,
    answer_value INT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_answers_user_question_key UNIQUE (user_id, question_id)
);

-- ========================================================
-- 7. Seed data for default questions
-- ========================================================
INSERT INTO questions (question_text, question_type, options) VALUES
('How do you prefer to spend your free time?', 'multiple_choice', ARRAY[
    'Socializing with friends',
    'Reading or watching movies',
    'Engaging in sports or physical activities',
    'Learning something new'
]),
('When making decisions, do you rely more on logic or emotions?', 'multiple_choice', ARRAY[
    'Completely on logic',
    'Mostly on logic but some emotion',
    'Mostly on emotion',
    'Completely on emotion'
]),
('Are you more introverted or extroverted?', 'multiple_choice', ARRAY[
    'Strongly introverted',
    'Somewhat introverted',
    'Somewhat extroverted',
    'Strongly extroverted'
]),
('Do you enjoy outdoor activities like hiking or camping?', 'multiple_choice', ARRAY[
    'Yes, I love them',
    'I enjoy them occasionally',
    'Not really my thing',
    'I prefer staying indoors'
]),
('What is your stance on travel?', 'multiple_choice', ARRAY[
    'I love traveling often and exploring new places',
    'I like traveling but don’t do it often',
    'I prefer short trips or staycations',
    'I’m not into travel much'
]),
('How important is fitness in your daily routine?', 'multiple_choice', ARRAY[
    'Very important, I exercise regularly',
    'Somewhat important, I try to stay active',
    'Neutral, I don’t prioritize it',
    'Not important at all'
]),
('How important is family to you?', 'multiple_choice', ARRAY[
    'Extremely important, family comes first',
    'Important, but I value my independence too',
    'Somewhat important, I prioritize other aspects of life',
    'Not that important to me'
]),
('What is your perspective on long-term relationships?', 'multiple_choice', ARRAY[
    'I’m looking for a long-term committed relationship',
    'I’m open to a long-term relationship but not actively looking',
    'I prefer casual relationships for now',
    'I’m not interested in relationships currently'
]),
('How do you handle disagreements in a relationship?', 'multiple_choice', ARRAY[
    'I like to resolve them through calm discussion',
    'I avoid conflict as much as possible',
    'I express my feelings openly and expect the same',
    'I tend to let things cool down before addressing them'
]),
('How often do you read books?', 'multiple_choice', ARRAY[
    'Daily',
    'Weekly',
    'Occasionally',
    'Rarely or never'
]),
('Which types of movies do you prefer?', 'multiple_choice', ARRAY[
    'Action/Adventure',
    'Comedy',
    'Drama/Romance',
    'Sci-Fi/Fantasy'
]),
('What type of music do you enjoy most?', 'multiple_choice', ARRAY[
    'Pop/Top 40',
    'Rock/Alternative',
    'Classical/Jazz',
    'Hip-hop/Rap'
]),
('Would you be okay with dating someone who smokes?', 'multiple_choice', ARRAY[
    'Yes, it doesn’t bother me',
    'I don’t mind it occasionally',
    'I would prefer not to date someone who smokes',
    'It’s a dealbreaker for me'
]),
('How do you feel about pets in a relationship?', 'multiple_choice', ARRAY[
    'I love pets and have/want one',
    'I like pets but don’t have one',
    'I’m indifferent to pets',
    'I’m not a fan of pets'
]),
('What’s your favorite way to relax after a long day?', 'multiple_choice', ARRAY[
    'Watching TV or movies',
    'Reading a book',
    'Spending time with friends',
    'Exercising or going for a walk'
]);

COMMIT;
