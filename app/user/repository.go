package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"training/app"
	"training/persistence"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Insert(ctx context.Context, model persistence.User) (*persistence.User, error)
	Update(ctx context.Context, model persistence.User) error
	Delete(ctx context.Context, userId uuid.UUID) error
	SelectById(ctx context.Context, userId uuid.UUID) (*persistence.User, error)
}

type repository struct {
	db    *pgxpool.Pool
	redis redis.UniversalClient
}

func NewRepository(db *pgxpool.Pool, redis redis.UniversalClient) Repository {
	return &repository{
		db:    db,
		redis: redis,
	}
}

func (r *repository) Insert(ctx context.Context, user persistence.User) (*persistence.User, error) {
	user.UserId = uuid.Must(uuid.NewV7())

	_, err := r.db.Exec(ctx,
		"insert into users(user_id, user_email, user_name, created_by, created_at) values($1, $2, $3, $4, $5)",
		user.UserId,
		user.UserEmail,
		user.UserName,
		user.CreatedBy,
		user.CreatedAt,
	)
	return &user, err
}

func (r *repository) Update(ctx context.Context, user persistence.User) error {
	cmd, err := r.db.Exec(ctx,
		"update users set user_email = $1, user_name = $2, updated_by = $3, updated_at = $4 where user_id = $5",
		user.UserEmail,
		user.UserName,
		user.UpdatedBy,
		user.UpdatedAt,
		user.UserId,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return app.ErrNotFound
	}

	err = r.removeCacheFromRedis(ctx, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, userId uuid.UUID) error {
	cmd, err := r.db.Exec(ctx,
		"delete from users where user_id = $1",
		userId,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return app.ErrNotFound
	}

	err = r.removeCacheFromRedis(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) SelectById(ctx context.Context, userId uuid.UUID) (*persistence.User, error) {
	user, err := r.getByIdFromRedis(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = r.selectByIdFromDB(ctx, userId)
		if err != nil {
			return nil, err
		}

		if user != nil {
			err = r.cacheToRedis(ctx, *user)
			if err != nil {
				return nil, err
			}
		}
	}
	return user, nil
}

func (r *repository) selectByIdFromDB(ctx context.Context, userId uuid.UUID) (*persistence.User, error) {
	var user persistence.User
	if err := r.db.QueryRow(ctx,
		"SELECT user_id, user_email, user_name, created_by, created_at, updated_by, updated_at FROM users WHERE user_id = $1",
		userId,
	).Scan(&user.UserId,
		&user.UserEmail,
		&user.UserName,
		&user.CreatedBy,
		&user.CreatedAt,
		&user.UpdatedBy,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) getByIdFromRedis(ctx context.Context, id uuid.UUID) (*persistence.User, error) {
	byteResult, err := r.redis.Get(ctx, fmt.Sprintf("%v:%v", app.RedisUserKey, id.String())).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var user persistence.User
	if err := json.NewDecoder(bytes.NewReader(byteResult)).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) cacheToRedis(ctx context.Context, user persistence.User) error {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(user); err != nil {
		return err
	}

	if err := r.redis.Set(ctx, fmt.Sprintf("%v:%v", app.RedisUserKey, user.UserId.String()), buffer.String(), 10*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) removeCacheFromRedis(ctx context.Context, id uuid.UUID) error {
	if err := r.redis.Del(ctx, fmt.Sprintf("%v:%v", app.RedisUserKey, id.String())).Err(); err != nil {
		return err
	}
	return nil
}
