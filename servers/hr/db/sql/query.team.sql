-- name: GetTeamsPerPage :many
select * from team 
LIMIT $1 
OFFSET (($1 * @page::int ) - $1);

-- name: GetTeamByID :one
select * from team where id=$1;

-- name: GetTotalTeams :one
select count(*) from team;

-- name: CreateTeam :one
INSERT INTO team (
  name, description
) VALUES (
  $1,$2
)
RETURNING *;

-- name: UpdateTeam :exec
UPDATE team
  set name = $2,
  description = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
