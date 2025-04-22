-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateAuthors :copyfrom
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
);


-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- name: CreateNewUser :one
WITH new_user AS (
  INSERT INTO "public"."user"
  ("name", "email", "email_verified", "image")
       VALUES ($1, $2, $3, $4)
       RETURNING id AS user_id
), new_account AS (
  INSERT INTO "public"."account"
  ("account_id", "provider_id", "access_token", "refresh_token", "id_token", "scope","user_id")
       SELECT $5, $6, $7, $8, $9, $10, user_id
         FROM new_user
       RETURNING id AS account_id, user_id
), new_session AS (
  INSERT INTO "public"."session"
  ("expires_at", "token", "ip_address", "user_agent", "user_id")
       SELECT $11,
              $12,
              $13,
              $14,
              a.user_id
         FROM new_account a
       RETURNING id AS session_id, user_id
 )
SELECT
*
FROM
new_session;
