-- +goose Up
-- +goose StatementBegin
INSERT INTO hospitals (name, hospital_url)
VALUES ('Hospital B', 'hospital-b')
ON CONFLICT (hospital_url) DO NOTHING;

-- Thai Patients (2)
INSERT INTO patients (
    national_id, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, gender, phone_number, email
) VALUES
('1200000000001', NULL, 'กิตติ', 'แสนดี', 'Kitti', 'Saendee', '1993-03-18', 'M', '0861111111', 'kitti@hospital-b.mock'),
('1200000000002', NULL, 'อรทัย', 'บุญมี', 'Orathai', 'Boonmee', '1995-12-02', 'F', '0862222222', 'orathai@hospital-b.mock')
ON CONFLICT (national_id) DO NOTHING;

-- Foreign Patients (1)
INSERT INTO patients (
    national_id, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, gender, phone_number, email
) VALUES
(NULL, 'JP7654321', NULL, NULL, 'Aiko', 'Tanaka', '1991-08-27', 'F', '0863333333', 'aiko.tanaka@hospital-b.mock')
ON CONFLICT (passport_id) DO NOTHING;

INSERT INTO patient_hospitals (patient_id, hospital_id, patient_hn)
SELECT p.id, h.id, mapping.patient_hn
FROM (
    VALUES
        ('NID', '1200000000001', 'HN-B-0001'),
        ('NID', '1200000000002', 'HN-B-0002'),
        ('PPT', 'JP7654321', 'HN-B-0003')
) AS mapping(id_type, ext_id, patient_hn)
JOIN patients p ON (mapping.id_type = 'NID' AND p.national_id = mapping.ext_id) OR (mapping.id_type = 'PPT' AND p.passport_id = mapping.ext_id)
JOIN hospitals h ON h.hospital_url = 'hospital-b'
ON CONFLICT (patient_id, hospital_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM patients
WHERE national_id IN ('1200000000001', '1200000000002')
   OR passport_id IN ('JP7654321');

DELETE FROM hospitals
WHERE hospital_url = 'hospital-b';
-- +goose StatementEnd
