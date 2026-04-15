CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{"STUDENT"}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);

CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_teachers_name ON teachers(name);

CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS idx_subjects_name ON subjects(name);

CREATE TABLE IF NOT EXISTS academic_periods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL DEFAULT 'Semester',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    mode VARCHAR(20) NOT NULL DEFAULT 'PARITY',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT chk_period_mode CHECK (mode IN ('PARITY', 'OTHER')),
    CONSTRAINT chk_period_dates CHECK (end_date >= start_date)
);
