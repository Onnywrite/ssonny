package postgres

import (
	"context"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/cuteql"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func (pg *PgStorage) TruncateTableUsers(ctx context.Context) error {
	tx, err := cuteql.Execute(ctx, pg.db, nil, `TRUNCATE TABLE users CASCADE`)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) SaveUser(ctx context.Context, user models.User) (*models.User, repo.Transactor, error) {
	return cuteql.GetSquirreled[models.User](ctx, pg.db, nil,
		squirrel.
			Insert("users").
			Columns(
				"user_nickname", "user_email", "user_verified",
				"user_gender", "user_password_hash", "user_birthday").
			Values(user.Nickname, user.Email, user.Verified,
				user.Gender, user.PasswordHash, user.Birthday).
			Suffix("RETURNING *"),
	)
}

func (pg *PgStorage) UpdateUser(ctx context.Context, userId uuid.UUID, newValues map[string]any) error {
	if len(newValues) == 0 {
		return eris.Wrap(repo.ErrEmptyResult, "no fields to update")
	}

	if _, ok := newValues["user_id"]; ok {
		return eris.Wrap(repo.ErrInternal, "user_id must not be changed")
	}

	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db, nil,
		squirrel.
			Update("users").
			SetMap(newValues).
			Where("user_id = ?", userId),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_email": email})
}

func (pg *PgStorage) UserByNickname(ctx context.Context, nickname string) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_nickname": nickname})
}

func (pg *PgStorage) UserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_id": id})
}

func (pg *PgStorage) getUserWhere(ctx context.Context, where squirrel.Sqlizer) (*models.User, error) {
	user, tx, err := cuteql.GetSquirreled[models.User](ctx, pg.db, nil,
		squirrel.
			Select("*").
			From("users").
			Where(where),
	)
	if err != nil {
		return nil, err
	}

	return user, cuteql.Commit(tx)
}
