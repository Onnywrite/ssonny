package postgres

import (
	"context"
	"time"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/Onnywrite/ssonny/internal/lib/cuteql"
	"github.com/Onnywrite/ssonny/internal/storage/repo"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func (pg *PgStorage) TruncateTableTokens(ctx context.Context) error {
	tx, err := cuteql.Execute(ctx, pg.db, `TRUNCATE TABLE tokens CASCADE`)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) SaveToken(ctx context.Context, token models.Token) (uint64, repo.Transactor, error) {
	id, tx, err := cuteql.Get[uint64](ctx, pg.db, `
		INSERT INTO tokens (token_user_fk, token_app_fk, token_rotation,
			token_rotated_at, token_platform, token_agent)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING token_id
	`, token.UserId, token.AppId, token.Rotation,
		token.RotatedAt, token.Platform, token.Agent)
	if err != nil {
		return 0, nil, err
	}

	return *id, tx, err
}

func (pg *PgStorage) UpdateToken(ctx context.Context, id uint64, newValues map[string]any) error {
	if len(newValues) == 0 {
		return eris.Wrap(repo.ErrEmptyResult, "no fields to update")
	}

	if _, ok := newValues["token_id"]; ok {
		return eris.Wrap(repo.ErrInternal, "token_id must not be changed")
	}

	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db,
		squirrel.
			Update("tokens").
			SetMap(newValues).
			Where("token_id = ?", id),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) Token(ctx context.Context, id uint64) (*models.Token, error) {
	token, tx, err := cuteql.Get[models.Token](ctx, pg.db, `
		SELECT * FROM tokens
		WHERE token_id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	return token, cuteql.Commit(tx)
}

func (pg *PgStorage) DeleteTokens(ctx context.Context, userId uuid.UUID, appId *uint64) error {
	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db,
		squirrel.
			Delete("tokens").
			Where("token_user_fk = ?", userId).
			Where(squirrel.Eq{"token_app_fk": appId}),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) DeleteToken(ctx context.Context, tokenId uint64) error {
	tx, err := cuteql.Execute(ctx, pg.db, `
		DELETE FROM tokens
		WHERE token_id = $1
	`, tokenId)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) CountTokens(ctx context.Context, userId uuid.UUID, appId *uint64) (uint64, error) {
	count, tx, err := cuteql.GetSquirreled[uint64](ctx, pg.db,
		squirrel.
			Select("COUNT(*)").
			From("tokens").
			Where("token_user_fk = ?", userId).
			Where(squirrel.Eq{"token_app_fk": appId}),
	)
	if err != nil {
		return 0, err
	}

	return *count, cuteql.Commit(tx)
}

func (pg *PgStorage) DeleteTokensWithRotatedAtBefore(ctx context.Context, before time.Time,
) ([]uint64, error) {
	ids, tx, err := cuteql.Query[uint64](ctx, pg.db, `
		DELETE FROM tokens
		WHERE token_rotated_at < $1
		RETURNING token_id
	`, before)
	if err != nil {
		return nil, err
	}

	return ids, cuteql.Commit(tx)
}
