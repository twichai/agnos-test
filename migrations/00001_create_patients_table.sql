-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS hospitals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),

    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),

    date_of_birth DATE,

    national_id VARCHAR(20),
    passport_id VARCHAR(20),

    phone_number VARCHAR(50),
    email VARCHAR(255),

    gender VARCHAR(1) CHECK (gender IN ('M','F')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(national_id),
    UNIQUE(passport_id)
);

CREATE TABLE IF NOT EXISTS patient_hospitals (
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    patient_hn VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (patient_id, hospital_id),
    UNIQUE(hospital_id, patient_hn)
);

CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(username, hospital_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS staffs;
DROP TABLE IF EXISTS patient_hospitals;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS hospitals;
-- +goose StatementEnd
