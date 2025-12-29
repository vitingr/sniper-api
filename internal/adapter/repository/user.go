package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/internal/domain/entity"
	"main.go/pkg/logger"
	util "main.go/pkg/utils/password"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	hashed_password, err := util.HashPassword(user.Password)

	if err != nil {
		logger.Error("Error hashing password")
	}

	query := `
		INSERT INTO users (name, email, password, profile_picture, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return repository.db.QueryRow(
		ctx, query,
		user.Name, user.Email, hashed_password, user.ProfilePicture, user.CreatedAt,
	).Scan(&user.ID, &user.CreatedAt)
}

func (repository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, profile_picture, created_at, updated_at FROM users WHERE email = $1`

	user := &entity.User{}
	err := repository.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) ListAllUsers(ctx context.Context, email string) ([]*entity.User, error) {
	query := `SELECT * FROM users WHERE email != $1`

	rows, err := repository.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*entity.User{}

	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password,
			&user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, profile_picture = $4, updated_at = NOW()
		WHERE id = $5
	`

	_, err := repository.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.ProfilePicture, user.ID)

	return err
}

func (repository *UserRepository) DeleteUser(ctx context.Context, email string) (*entity.User, error) {
	query := `
		DELETE FROM users WHERE email = $1
		RETURNING id, name, email, password, profile_picture, created_at, updated_at
	`

	user := &entity.User{}
	err := repository.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	return user, nil
}
