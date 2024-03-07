-- name: GetTeamMembersPerPage :many
select * from team_member 
LIMIT $1 
OFFSET (($1 * @page::int ) - $1);

-- name: GetTotalTeamMembers :one
select count(*) from team_member;

-- name: CreateTeamMember :one
INSERT INTO team_member (
  member_id, team_id
) VALUES (
  $1, $2
)
RETURNING *;


-- name: DeleteTeamMember :exec
DELETE FROM team_member WHERE member_id=$1 AND team_id=$2;