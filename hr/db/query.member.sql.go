// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.member.sql

package db

import (
	"context"
)

const createMember = `-- name: CreateMember :one
INSERT INTO member (
  email, role, registration_secret
) VALUES (
  $1,$2,$3
)
RETURNING id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at
`

type CreateMemberParams struct {
	Email              string
	Role               Roles
	RegistrationSecret string
}

func (q *Queries) CreateMember(ctx context.Context, arg CreateMemberParams) (Member, error) {
	row := q.db.QueryRowContext(ctx, createMember, arg.Email, arg.Role, arg.RegistrationSecret)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Password,
		&i.Email,
		&i.RegistrationSecret,
		&i.EmailConfirmed,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMember = `-- name: DeleteMember :exec
DELETE FROM member where id = $1
`

func (q *Queries) DeleteMember(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteMember, id)
	return err
}

const getMemberByEmail = `-- name: GetMemberByEmail :one
select id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at from member where email=$1
`

func (q *Queries) GetMemberByEmail(ctx context.Context, email string) (Member, error) {
	row := q.db.QueryRowContext(ctx, getMemberByEmail, email)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Password,
		&i.Email,
		&i.RegistrationSecret,
		&i.EmailConfirmed,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMemberByID = `-- name: GetMemberByID :one
select id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at from member where id=$1
`

func (q *Queries) GetMemberByID(ctx context.Context, id int32) (Member, error) {
	row := q.db.QueryRowContext(ctx, getMemberByID, id)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Password,
		&i.Email,
		&i.RegistrationSecret,
		&i.EmailConfirmed,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMemberByUsername = `-- name: GetMemberByUsername :one
select id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at from member where username=$1
`

func (q *Queries) GetMemberByUsername(ctx context.Context, username string) (Member, error) {
	row := q.db.QueryRowContext(ctx, getMemberByUsername, username)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Password,
		&i.Email,
		&i.RegistrationSecret,
		&i.EmailConfirmed,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMemberByUsernameAndPassword = `-- name: GetMemberByUsernameAndPassword :one
select id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at from member where username=$1 AND password=$2
`

type GetMemberByUsernameAndPasswordParams struct {
	Username string
	Password string
}

func (q *Queries) GetMemberByUsernameAndPassword(ctx context.Context, arg GetMemberByUsernameAndPasswordParams) (Member, error) {
	row := q.db.QueryRowContext(ctx, getMemberByUsernameAndPassword, arg.Username, arg.Password)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Fullname,
		&i.Password,
		&i.Email,
		&i.RegistrationSecret,
		&i.EmailConfirmed,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMembersPerPage = `-- name: GetMembersPerPage :many
select id, username, fullname, password, email, registration_secret, email_confirmed, role, created_at, updated_at from member
LIMIT $1 
OFFSET (($1 * $2::int ) - $1)
`

type GetMembersPerPageParams struct {
	Limit int32
	Page  int32
}

func (q *Queries) GetMembersPerPage(ctx context.Context, arg GetMembersPerPageParams) ([]Member, error) {
	rows, err := q.db.QueryContext(ctx, getMembersPerPage, arg.Limit, arg.Page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Member
	for rows.Next() {
		var i Member
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Fullname,
			&i.Password,
			&i.Email,
			&i.RegistrationSecret,
			&i.EmailConfirmed,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTotalMembers = `-- name: GetTotalMembers :one
select count(*) from member
`

func (q *Queries) GetTotalMembers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalMembers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateMemberForRegistration = `-- name: UpdateMemberForRegistration :exec
UPDATE member
  set username = $2,
  fullname = $3,
  password = $4,
  email_confirmed = $5,
  registration_secret = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateMemberForRegistrationParams struct {
	ID                 int32
	Username           string
	Fullname           string
	Password           string
	EmailConfirmed     bool
	RegistrationSecret string
}

func (q *Queries) UpdateMemberForRegistration(ctx context.Context, arg UpdateMemberForRegistrationParams) error {
	_, err := q.db.ExecContext(ctx, updateMemberForRegistration,
		arg.ID,
		arg.Username,
		arg.Fullname,
		arg.Password,
		arg.EmailConfirmed,
		arg.RegistrationSecret,
	)
	return err
}

const updateMemberRole = `-- name: UpdateMemberRole :exec
UPDATE member
  set role = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateMemberRoleParams struct {
	ID   int32
	Role Roles
}

func (q *Queries) UpdateMemberRole(ctx context.Context, arg UpdateMemberRoleParams) error {
	_, err := q.db.ExecContext(ctx, updateMemberRole, arg.ID, arg.Role)
	return err
}