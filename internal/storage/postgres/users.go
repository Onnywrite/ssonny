package postgres

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/cuteql"
	"github.com/Onnywrite/ssonny/internal/storage/repo"
)

func (pg *PgStorage) SaveUser(ctx context.Context, user models.User) (*models.User, repo.Transactor, error) {
	return cuteql.GetNamed[models.User, models.User](ctx, pg.db, nil, `
		INSERT INTO users 
			(user_nickname, user_email, user_verified,
			user_gender, user_password_hash, user_birthday)
		VALUES
			(:user_nickname, :user_email, :user_verified,
			:user_gender, :user_password_hash, :user_birthday)
		RETURNING *
	`, user)
}
