package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

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

func (r *PostgresURLRepository) Ping(ctx echo.Context) error {
	err := r.PostgresPool.db.Ping(ctx.Request().Context())
	return err
}

func (r *PostgresURLRepository) Insert(ctx echo.Context, url model.URL) (*model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	var savedURL model.URL
	err = tx.QueryRow(ctx.Request().Context(),
		`
		insert into url_shortener.url (id, original_url, short_url) 
		values ($1, $2, $3) 
		on conflict (id) do update 
		set original_url = $2, short_url = $3 
		returning id, original_url, short_url
		`,
		url.ID, url.Original, url.Shortened).Scan(&savedURL.ID, &savedURL.Original, &savedURL.Shortened)

	if err != nil {
		return nil, apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return &savedURL, nil
}

func (r *PostgresURLRepository) SelectByID(ctx echo.Context, key string) (*model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	var url model.URL
	err = tx.QueryRow(ctx.Request().Context(),
		`
		select id, original_url, short_url
		from url_shortener.url
		where id = $1
		`,
		key).Scan(&url.ID, &url.Original, &url.Shortened)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = apperrors.NewValueError("url not found", utils.Caller(), apperrors.ErrorURLNotFound)
		} else {
			err = apperrors.NewValueError("query failed", utils.Caller(), err)
		}
		return nil, err
	}

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return &url, nil
}

func (r *PostgresURLRepository) SelectAll(ctx echo.Context) ([]model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	var urls []model.URL
	rows, err := tx.Query(ctx.Request().Context(),
		`
		select id, original_url, short_url
		from url_shortener.url
		`)
	if err != nil {
		return nil, apperrors.NewValueError("query failed", utils.Caller(), err)
	}
	defer rows.Close()

	for rows.Next() {
		var url model.URL
		err = rows.Scan(&url.ID, &url.Original, &url.Shortened)
		if err != nil {
			return nil, apperrors.NewValueError("unable to scan values", utils.Caller(), err)
		}
		urls = append(urls, url)
	}

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return urls, nil
}

func (r *PostgresURLRepository) DeleteAll(ctx echo.Context) error {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	_, err = tx.Exec(ctx.Request().Context(),
		`
	delete from url_shortener.url
	`)

	if err != nil {
		return apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return nil
}

func (r *PostgresURLRepository) InsertBatch(ctx echo.Context, urls []model.URL) ([]model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	rows := make([][]interface{}, len(urls))
	for i, url := range urls {
		row := []interface{}{url.ID, url.Original, url.Shortened}
		rows[i] = row
	}

	_, err = tx.Exec(ctx.Request().Context(),
		`
		create temporary table _temp_upsert_urls (like url_shortener.url) on commit drop
		`)
	if err != nil {
		return nil, apperrors.NewValueError("unable to create temp table", utils.Caller(), err)
	}

	count, err := tx.CopyFrom(
		ctx.Request().Context(),
		pgx.Identifier{"pg_temp", "_temp_upsert_urls"},
		[]string{"id", "original_url", "short_url"}, 
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, apperrors.NewValueError("copy from failed", utils.Caller(), err)
	}
	if count != int64(len(urls)) {
		return nil, apperrors.NewValueError("not all rows were inserted", utils.Caller(), err)
	}

	quaryRows, err := tx.Query(ctx.Request().Context(),
		`
		insert into url_shortener.url (id, original_url, short_url) 
		select id, original_url, short_url from pg_temp._temp_upsert_urls 
		on conflict (id) do update set original_url = excluded.original_url, short_url = excluded.short_url
		returning id, original_url, short_url 
		`)
	if err != nil {
		return nil, apperrors.NewValueError("unable to upsert batch", utils.Caller(), err)
	}
	defer quaryRows.Close()

	var savedUrls []model.URL
	for quaryRows.Next() {
		var url model.URL
		err = quaryRows.Scan(&url.ID, &url.Original, &url.Shortened)
		if err != nil {
			return nil, apperrors.NewValueError("unable to scan values", utils.Caller(), err)
		}
		savedUrls = append(savedUrls, url)
	}

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return savedUrls, nil
}
