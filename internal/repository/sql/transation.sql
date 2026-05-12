-- name: addTX :one 
INSERT INTO transactions
(user_id, 
donation_id, 
transaction_type, 
amount, reated_at)
VALUES (
    $1, $2, $3, $4, NOW()
)
RETURNING * ;

-- name: HistoryTX :many 
SELECT * 
FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;

