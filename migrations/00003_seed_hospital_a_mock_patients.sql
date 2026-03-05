-- +goose Up
-- +goose StatementBegin
INSERT INTO hospitals (name, hospital_url)
VALUES ('Hospital A', 'hospital-a')
ON CONFLICT (hospital_url) DO NOTHING;

-- Thai Patients (3)
INSERT INTO patients (
    national_id, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, gender, phone_number, email
) VALUES
('1100000000001', NULL, 'สมชาย', 'ใจดี', 'Somchai', 'Jaidee', '1990-01-01', 'M', '0811111111', 'somchai@hospital-a.mock'),
('1100000000002', NULL, 'สมศรี', 'รักดี', 'Somsri', 'Rakdee', '1988-07-14', 'F', '0822222222', 'somsri@hospital-a.mock'),
('1100000000003', NULL, 'มานะ', 'ยืนยง', 'Mana', 'Yuenyong', '2001-11-23', 'M', '0833333333', 'mana@hospital-a.mock')
ON CONFLICT (national_id) DO NOTHING;

-- Foreign Patients (2)
INSERT INTO patients (
    national_id, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, gender, phone_number, email
) VALUES
(NULL, 'US1234567', NULL, NULL, 'John', 'Smith', '1985-05-05', 'M', '0844444444', 'john.smith@hospital-a.mock'),
(NULL, 'UK9876543', NULL, NULL, 'Emma', 'Watson', '1992-09-09', 'F', '0855555555', 'emma.watson@hospital-a.mock')
ON CONFLICT (passport_id) DO NOTHING;

INSERT INTO patient_hospitals (patient_id, hospital_id, patient_hn)
SELECT p.id, h.id, mapping.patient_hn
FROM (
    VALUES
        ('NID', '1100000000001', 'HN-A-0001'),
        ('NID', '1100000000002', 'HN-A-0002'),
        ('NID', '1100000000003', 'HN-A-0003'),
        ('PPT', 'US1234567', 'HN-A-0004'),
        ('PPT', 'UK9876543', 'HN-A-0005')
) AS mapping(id_type, ext_id, patient_hn)
JOIN patients p ON (mapping.id_type = 'NID' AND p.national_id = mapping.ext_id) OR (mapping.id_type = 'PPT' AND p.passport_id = mapping.ext_id)
JOIN hospitals h ON h.hospital_url = 'hospital-a'
ON CONFLICT (patient_id, hospital_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM patients
WHERE national_id IN ('1100000000001', '1100000000002', '1100000000003')
   OR passport_id IN ('US1234567', 'UK9876543');

DELETE FROM hospitals
WHERE hospital_url = 'hospital-a';
-- +goose StatementEnd
