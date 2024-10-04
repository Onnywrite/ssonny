package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	cute "github.com/Onnywrite/ssonny/pkg/cuteql"
	"github.com/jmoiron/sqlx"
)

func (pg *PgStorage) SaveApp(ctx context.Context, app models.App, domainsIds []uint64,
) (*models.App, error) {
	tx, err := pg.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin tx: %w", err)
	}

	savedApp, err := cute.Get[models.App](ctx, cute.WithTx(tx), cute.WithQuery(`
INSERT INTO applications (app_owner_fk, app_name, app_description, app_secret_hash)
VALUES ($1, $2, $3, $4)
RETURNING *`),
		cute.WithArgs(app.OwnerId, app.Name, app.Description, app.SecretHash))
	if err != nil {
		return nil, err
	}

	err = pg.tieDomainsToApp(ctx, tx, savedApp.Id, domainsIds)
	if err != nil {
		return nil, err
	}

	return savedApp, nil
}

func (pg *PgStorage) SaveDomains(ctx context.Context, domains []models.Domain,
) ([]models.Domain, error) {
	args := make([][]any, len(domains))
	for i, d := range domains {
		args[i] = []any{d.OwnerId, d.Name, d.IsVerified, d.VerifiedAt}
	}

	savedDomains, err := cute.Query[models.Domain](ctx, cute.WithDb(pg.db), cute.WithQuery(`
INSERT INTO domains (domain_owner_fk, domain_name, domain_verified, domain_verified_at)
VALUES ($1, $2, $3, $4)
RETURNING *`),
		cute.WithBatch(args...))
	if err != nil {
		return nil, err
	}

	return savedDomains, nil
}

func (pg *PgStorage) TieDomainsToApp(ctx context.Context, appId uint64, domainsIds []uint64,
) error {
	return pg.tieDomainsToApp(ctx, nil, appId, domainsIds)
}

func (pg *PgStorage) tieDomainsToApp(ctx context.Context, tx *sqlx.Tx,
	appId uint64, domainsIds []uint64,
) error {
	args := make([][]any, len(domainsIds))
	for i, id := range domainsIds {
		args[i] = []any{appId, id}
	}

	txOpt := cute.WithTx(tx)
	if txOpt == nil {
		txOpt = cute.WithDb(pg.db)
	}

	err := cute.Execute(ctx, txOpt, cute.WithCommit(), cute.WithQuery(`
	INSERT INTO applications_domains (app_fk, domain_fk)
	VALUES ($1, $2)`),
		cute.WithBatch(args...),
	)
	if err != nil {
		return err
	}

	return nil
}
