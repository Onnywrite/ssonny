package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/cuteql"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
)

func (pg *PgStorage) SaveUser(ctx context.Context, user models.User) (*models.User, repo.Transactor, error) {
	return cuteql.GetSquirreled[models.User](ctx, pg.db, nil,
		squirrel.
			Insert("users").
			Columns(
				"user_nickname", "user_email", "user_verified",
				"user_gender", "user_password_hash", "user_birthday").
			Values(user.Nickname, user.Email, user.IsVerified,
				user.Gender, user.PasswordHash, user.Birthday).
			Suffix("RETURNING *"),
	)
}
