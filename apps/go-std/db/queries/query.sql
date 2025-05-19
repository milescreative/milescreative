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
  ON CONFLICT(email) DO UPDATE SET name = $1
  RETURNING id AS user_id
), new_account AS (
  INSERT INTO "public"."account"
  ("account_id", "provider_id", "access_token", "refresh_token", "id_token", "scope","user_id", "access_token_expires_at")
       SELECT $5, $6, $7, $8, $9, $10, user_id, $15
         FROM new_user
  ON CONFLICT(account_id) DO UPDATE SET
  provider_id = $6,
  access_token = $7,
  refresh_token = $8,
  id_token = $9,
  scope = $10,
  access_token_expires_at = $15

  RETURNING id AS account_id, user_id
), new_session AS (
  INSERT INTO "public"."session"
  ("expires_at", "token", "ip_address", "user_agent", "user_id", "account_id")
       SELECT $11,
              $12,
              $13,
              $14,
              a.user_id,
              a.account_id
         FROM new_account a
       RETURNING id AS session_id, user_id, token as session_token
 )
SELECT
*
FROM
new_session;



-- name: GetSessionByToken :one
SELECT session.*, account.refresh_token FROM "public"."session" session
INNER JOIN public.user  ON session.user_id = "user".id
INNER JOIN public.account account ON session.account_id = account.id
WHERE session.token = $1
LIMIT 1;


-- name: DeleteSession :one
DELETE FROM "public"."session"
WHERE token = $1
RETURNING id;

-- name: DeleteSessionsForUser :one
DELETE FROM "public"."session"
WHERE "user_id" = $1
RETURNING id;


-- name: UpdateAccount :exec
UPDATE "public"."account"
SET "access_token" = COALESCE(sqlc.narg('access_token'), access_token),
    "refresh_token" = COALESCE(sqlc.narg('refresh_token'), refresh_token),
    "id_token" = COALESCE(sqlc.narg('id_token'), id_token),
    "access_token_expires_at" = COALESCE(sqlc.narg('access_token_expires_at'), access_token_expires_at),
    "scope" = COALESCE(sqlc.narg('scope'), scope),
    "updated_at" = NOW()
WHERE "id" = sqlc.arg('account_id');




-- name: GetUserByID :one
SELECT * FROM "public"."user"
WHERE id = $1
LIMIT 1;

-- name: GetUserSessions :many
SELECT * FROM "public"."session"
WHERE "user_id" = $1;


-- name: UpdateUser :one
UPDATE "public"."user"
SET "name" = $2,
    "email" = $3,
    "image" = $4
WHERE id = $1
RETURNING id;


-- name: DeleteUser :one
DELETE FROM "public"."user"
WHERE id = $1
RETURNING id;



