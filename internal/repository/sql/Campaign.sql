-- name: CreateCampaign :one
INSERT INTO campaigns (
    creator_id,
    title,
    description, 
    target_amount,
    end_date 
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCurrentAmount :one
SELECT current_amount 
FROM campaigns
WHERE id = $1;

-- name: GetCampaignByID :one
SELECT *
FROM campaigns
WHERE id = $1;

-- name: DeleteCampaign :exec
DELETE FROM campaigns 
WHERE id = $1;

-- name: DecreaseCampaignAmount :one
UPDATE campaigns 
SET current_amount = current_amount - $2
WHERE id = $1AND current_amount >= $2
RETURNING *;

-- name: IncreaseCampaignAmount :one
UPDATE campaigns 
SET current_amount = current_amount + $2
WHERE id = $1
RETURNING *;
