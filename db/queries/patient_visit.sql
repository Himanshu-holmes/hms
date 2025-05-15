-- name: CreatePatientVisit :one
INSERT INTO patient_visits (
    patient_id, doctor_id, visit_date, symptoms, diagnosis, prescription, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetPatientVisitByID :one
SELECT * FROM patient_visits
WHERE id = $1
LIMIT 1;

-- name: ListPatientVisitsByPatientID :many
SELECT pv.*, u.first_name as doctor_first_name, u.last_name as doctor_last_name
FROM patient_visits pv
JOIN users u ON pv.doctor_id = u.id -- Join to get doctor's name
WHERE pv.patient_id = $1
ORDER BY pv.visit_date DESC
LIMIT $2
OFFSET $3;

-- name: ListPatientVisitsByDoctorID :many
SELECT pv.*, p.first_name as patient_first_name, p.last_name as patient_last_name
FROM patient_visits pv
JOIN patients p ON pv.patient_id = p.id -- Join to get patient's name
WHERE pv.doctor_id = $1
ORDER BY pv.visit_date DESC
LIMIT $2
OFFSET $3;

-- name: UpdatePatientVisit :one
UPDATE patient_visits
SET
    visit_date = COALESCE(sqlc.narg(visit_date), visit_date),
    symptoms = COALESCE(sqlc.narg(symptoms), symptoms),
    diagnosis = COALESCE(sqlc.narg(diagnosis), diagnosis),
    prescription = COALESCE(sqlc.narg(prescription), prescription),
    notes = COALESCE(sqlc.narg(notes), notes),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeletePatientVisit :exec
DELETE FROM patient_visits
WHERE id = $1;