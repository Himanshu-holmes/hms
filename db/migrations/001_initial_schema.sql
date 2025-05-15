-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Define ENUM types
CREATE TYPE user_role AS ENUM ('receptionist', 'doctor');
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other', 'prefer_not_to_say');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

-- Patients table
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender gender_enum,
    contact_phone VARCHAR(20) UNIQUE,
    contact_email VARCHAR(255) UNIQUE,
    address TEXT,
    medical_history TEXT,
    registered_by_user_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ, -- For soft deletes

    CONSTRAINT fk_patients_registered_by_user
        FOREIGN KEY(registered_by_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_patients_last_name ON patients(last_name);
CREATE INDEX idx_patients_contact_phone ON patients(contact_phone) WHERE contact_phone IS NOT NULL;
CREATE INDEX idx_patients_contact_email ON patients(contact_email) WHERE contact_email IS NOT NULL;
CREATE INDEX idx_patients_deleted_at ON patients(deleted_at);

-- Patient Visits table
CREATE TABLE patient_visits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    visit_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    symptoms TEXT,
    diagnosis TEXT,
    prescription TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_patient_visits_patient
        FOREIGN KEY(patient_id)
        REFERENCES patients(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_patient_visits_doctor
        FOREIGN KEY(doctor_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_patient_visits_patient_id ON patient_visits(patient_id);
CREATE INDEX idx_patient_visits_doctor_id ON patient_visits(doctor_id);
CREATE INDEX idx_patient_visits_visit_date ON patient_visits(visit_date);

-- Trigger function for updated_at
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Apply trigger to tables
CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_patients_updated_at
BEFORE UPDATE ON patients
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_patient_visits_updated_at
BEFORE UPDATE ON patient_visits
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TRIGGER IF EXISTS set_patient_visits_updated_at ON patient_visits;
DROP TRIGGER IF EXISTS set_patients_updated_at ON patients;
DROP TRIGGER IF EXISTS set_users_updated_at ON users;

DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP TABLE IF EXISTS patient_visits;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS gender_enum;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS "uuid-ossp";