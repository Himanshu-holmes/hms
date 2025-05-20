-- name: CreatePatient :one
INSERT INTO patients (
    first_name, last_name, date_of_birth, gender,
    contact_phone, contact_email, address, medical_history,
    registered_by_user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPatientByID :one
SELECT * FROM patients
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListPatients :many
SELECT * FROM patients
WHERE deleted_at IS NULL
ORDER BY last_name, first_name
LIMIT $1
OFFSET $2;

-- name: UpdatePatient :one
UPDATE patients
SET
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    date_of_birth = COALESCE(sqlc.narg(date_of_birth), date_of_birth),
    gender = COALESCE(sqlc.narg(gender), gender),
    contact_phone = COALESCE(sqlc.narg(contact_phone), contact_phone),
    contact_email = COALESCE(sqlc.narg(contact_email), contact_email),
    address = COALESCE(sqlc.narg(address), address),
    medical_history = COALESCE(sqlc.narg(medical_history), medical_history),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeletePatient :one
UPDATE patients
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: HardDeletePatient :exec
DELETE FROM patients
WHERE id = $1;

-- Specific update for doctor (e.g., only medical_history)
-- name: UpdatePatientMedicalInfo :one
UPDATE patients
SET
    medical_history = $2,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CountPatients :one
SELECT COUNT(*) FROM patients
WHERE deleted_at IS NULL;