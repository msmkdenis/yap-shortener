package db

import (
	"context"
	_ "embed"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

//go:embed queries/insert_url_and_return.sql
var insertURLAndReturn string

//go:embed queries/select_url_by_id.sql
var selectURLById string

//go:embed queries/select_all_urls.sql
var selectAllURLs string

//go:embed queries/delete_all_urls.sql
var deleteAllURLs string

type PostgresURLRepository struct {
	PostgresPool *PostgresPool
	logger       *zap.Logger
}

func NewPostgresURLRepository(postgresPool *PostgresPool, logger *zap.Logger) *PostgresURLRepository {
	return &PostgresURLRepository{
		PostgresPool: postgresPool,
		logger:       logger,
	}
}

func (r *PostgresURLRepository) Ping(ctx context.Context) error {
	return r.PostgresPool.db.Ping(ctx)
}

func (r *PostgresURLRepository) Insert(ctx context.Context, url model.URL) (*model.URL, error) {
	var savedURL model.URL
	err := r.PostgresPool.db.QueryRow(ctx, insertURLAndReturn,
		url.ID, url.Original, url.Shortened).Scan(&savedURL.ID, &savedURL.Original, &savedURL.Shortened)

	if err != nil {
		return nil, apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	return &savedURL, nil
}

func (r *PostgresURLRepository) SelectByID(ctx context.Context, key string) (*model.URL, error) {
	var url model.URL
	err := r.PostgresPool.db.QueryRow(ctx, selectURLById, key).Scan(&url.ID, &url.Original, &url.Shortened, &url.CorrelationID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = apperrors.NewValueError("url not found", utils.Caller(), apperrors.ErrURLNotFound)
		} else {
			err = apperrors.NewValueError("query failed", utils.Caller(), err)
		}
		return nil, err
	}

	return &url, nil
}

func (r *PostgresURLRepository) SelectAll(ctx context.Context) ([]model.URL, error) {
	queryRows, err := r.PostgresPool.db.Query(ctx, selectAllURLs)
	if err != nil {
		return nil, apperrors.NewValueError("query failed", utils.Caller(), err)
	}
	defer queryRows.Close()

	urls, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.URL])
	if err != nil {
		return nil, apperrors.NewValueError("unable to collect rows", utils.Caller(), err)
	}

	return urls, nil
}

func (r *PostgresURLRepository) DeleteAll(ctx context.Context) error {
	_, err := r.PostgresPool.db.Exec(ctx, deleteAllURLs)

	if err != nil {
		return apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	return nil
}

func (r *PostgresURLRepository) InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx)

	rows := make([][]interface{}, len(urls))
	for i, url := range urls {
		row := []interface{}{url.ID, url.Original, url.Shortened, url.CorrelationID}
		rows[i] = row
	}

	tempTable := uuid.New().String()
	re := regexp.MustCompile(`\d|-`)
	tempTable = re.ReplaceAllString(tempTable, "")

	_, err = tx.Exec(ctx,
		`
		create temporary table `+tempTable+` (like url_shortener.url) on commit drop
		`)
	if err != nil {
		return nil, apperrors.NewValueError("unable to create temp table", utils.Caller(), err)
	}

	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pg_temp", tempTable},
		[]string{"id", "original_url", "short_url", "correlation_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, apperrors.NewValueError("copy from failed", utils.Caller(), err)
	}
	if count != int64(len(urls)) {
		return nil, apperrors.NewValueError("not all rows were inserted", utils.Caller(), err)
	}

	queryRows, err := tx.Query(ctx,
		`
		insert into url_shortener.url (id, original_url, short_url, correlation_id) 
		select id, original_url, short_url, correlation_id from pg_temp.`+tempTable+` 
		on conflict (id) do update set original_url = excluded.original_url, short_url = excluded.short_url, correlation_id = excluded.correlation_id
		returning id, original_url, short_url, correlation_id 
		`)
	if err != nil {
		return nil, apperrors.NewValueError("unable to upsert batch", utils.Caller(), err)
	}
	defer queryRows.Close()

	savedURLs, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.URL])
	if err != nil {
		return nil, apperrors.NewValueError("unable to collect rows", utils.Caller(), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return savedURLs, nil
}
