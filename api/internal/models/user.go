package models

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
)

var UserNodeID int64 = 1

type User struct {
	ID int64 `json:"id"`
}

type UserRepo struct {
	context.Context

	db   *pgxpool.Pool
	sgen *snowflake.Node
}

func NewUserRepo(ctx context.Context, db *pgxpool.Pool) *UserRepo {
	sgen, _ := snowflake.NewNode(UserNodeID)

	return &UserRepo{
		Context: ctx,
		db:      db,
		sgen:    sgen,
	}
}

func (u *UserRepo) CreateUser() (*User, error) {
	id := u.sgen.Generate().Int64()

	query := `
		INSERT INTO users (id)
		VALUES ($1)
		RETURNING id
	`

	userRow := u.db.QueryRow(u, query, id)

	user := &User{}
	err := userRow.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
