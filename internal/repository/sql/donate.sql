-- name: GetListDonations :many
SELECT * 
FROM donations
WHERE campaign_id = $1;

-- name: CreateDonation :one
INSERT INTO donations (
    user_id,
    campaign_id,
    amount
) VALUES (
    $1, $2, $3
)
RETURNING *;
