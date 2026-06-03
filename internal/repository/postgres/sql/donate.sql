-- name: GetListDonations :many
SELECT * 
FROM donations
WHERE campaign_id = $1;

-- name: GetDonationByID :one
SELECT *
FROM donations
WHERE id = $1;

-- name: UpdateDonationStatus :one
UPDATE donations
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DonationRefunded :one
UPDATE donations
SET status = 'refund'
WHERE id = $1
RETURNING *;

-- name: CreateDonation :one
INSERT INTO donations (
    user_id,
    campaign_id,
    amount
) VALUES (
    $1, $2, $3
)
RETURNING *;
