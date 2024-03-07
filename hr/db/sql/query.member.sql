-- name: GetMemberByID :one
select * from member where id=$1;

-- name: GetMemberByUsernameAndPassword :one
select * from member where username=$1 AND password=$2;

-- name: GetMemberByUsername :one
select * from member where username=$1;

-- name: GetMemberByEmail :one
select * from member where email=$1;

-- name: GetMembersPerPage :many
select * from member
LIMIT $1 
OFFSET (($1 * @page::int ) - $1);

-- name: GetTotalMembers :one
select count(*) from member;

-- name: CreateMember :one
INSERT INTO member (
  email, role, registration_secret
) VALUES (
  $1,$2,$3
)
RETURNING *;

-- name: UpdateMemberForRegistration :exec
UPDATE member
  set username = $2,
  fullname = $3,
  password = $4,
  email_confirmed = $5,
  registration_secret = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateMemberRole :exec
UPDATE member
  set role = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteMember :exec
DELETE FROM member where id = $1;


