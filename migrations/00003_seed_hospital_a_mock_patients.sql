-- +goose Up
-- +goose StatementBegin
INSERT INTO hospitals (name, hospital_url)
VALUES ('Hospital A', 'hospital-a')
ON CONFLICT (hospital_url) DO NOTHING;

INSERT INTO patients (
    hospital_id,
    patient_hn,
    national_id,
    passport_id,
    first_name_en,
    last_name_en,
    date_of_birth,
    gender,
    phone_number,
    email
)
SELECT h.id, 'HN-A-0001', 'NIDA000000000001', 'PA000000000001', 'John', 'Doe', '1990-01-01', 'M', '0811111111', 'john.doe@hospital-a.mock'
FROM hospitals h
WHERE h.hospital_url = 'hospital-a'
ON CONFLICT (hospital_id, patient_hn) DO NOTHING;

INSERT INTO patients (
    hospital_id,
    patient_hn,
    national_id,
    passport_id,
    first_name_en,
    last_name_en,
    date_of_birth,
    gender,
    phone_number,
    email
)
SELECT h.id, 'HN-A-0002', 'NIDA000000000002', 'PA000000000002', 'Jane', 'Smith', '1988-07-14', 'F', '0822222222', 'jane.smith@hospital-a.mock'
FROM hospitals h
WHERE h.hospital_url = 'hospital-a'
ON CONFLICT (hospital_id, patient_hn) DO NOTHING;

INSERT INTO patients (
    hospital_id,
    patient_hn,
    national_id,
    passport_id,
    first_name_en,
    last_name_en,
    date_of_birth,
    gender,
    phone_number,
    email
)
SELECT h.id, 'HN-A-0003', 'NIDA000000000003', 'PA000000000003', 'Alex', 'Brown', '2001-11-23', 'M', '0833333333', 'alex.brown@hospital-a.mock'
FROM hospitals h
WHERE h.hospital_url = 'hospital-a'
ON CONFLICT (hospital_id, patient_hn) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM patients
WHERE hospital_id IN (
    SELECT id FROM hospitals WHERE hospital_url = 'hospital-a'
)
AND patient_hn IN ('HN-A-0001', 'HN-A-0002', 'HN-A-0003');

DELETE FROM hospitals
WHERE hospital_url = 'hospital-a';
-- +goose StatementEnd
