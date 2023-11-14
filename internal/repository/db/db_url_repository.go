package db

import (
	"errors"
	"regexp"

	//"strings"

	"github.com/google/uuid"
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
	return r.PostgresPool.db.Ping(ctx.Request().Context())
}

func (r *PostgresURLRepository) Insert(ctx echo.Context, url model.URL) (*model.URL, error) {
	var savedURL model.URL
	err := r.PostgresPool.db.QueryRow(ctx.Request().Context(),
		`
		insert into url_shortener.url (id, original_url, short_url) 
		values ($1, $2, $3) 
		returning id, original_url, short_url
		`,
		url.ID, url.Original, url.Shortened).Scan(&savedURL.ID, &savedURL.Original, &savedURL.Shortened)

	if err != nil {
		return nil, apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	return &savedURL, nil
}

func (r *PostgresURLRepository) SelectByID(ctx echo.Context, key string) (*model.URL, error) {
	var url model.URL
	err := r.PostgresPool.db.QueryRow(ctx.Request().Context(),
		`
		select id, original_url, short_url, coalesce(correlation_id, '')
		from url_shortener.url
		where id = $1
		`,
		key).Scan(&url.ID, &url.Original, &url.Shortened, &url.CorrelationID)

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

func (r *PostgresURLRepository) SelectAll(ctx echo.Context) ([]model.URL, error) {
	queryRows, err := r.PostgresPool.db.Query(ctx.Request().Context(),
		`
		select id, original_url, short_url, coalesce(correlation_id, '')
		from url_shortener.url
		`)
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

func (r *PostgresURLRepository) DeleteAll(ctx echo.Context) error {
	_, err := r.PostgresPool.db.Exec(ctx.Request().Context(),
		`
		delete from url_shortener.url
		`)

	if err != nil {
		return apperrors.NewValueError("query failed", utils.Caller(), err)
	}

	return nil
}

func (r *PostgresURLRepository) InsertAllOrUpdate(ctx echo.Context, urls []model.URL) ([]model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("unable to start transaction", utils.Caller(), err)
	}
	defer tx.Rollback(ctx.Request().Context())

	rows := make([][]interface{}, len(urls))
	for i, url := range urls {
		row := []interface{}{url.ID, url.Original, url.Shortened, url.CorrelationID}
		rows[i] = row
	}

	tempTable := uuid.New().String()
	re, err := regexp.Compile(`\d|-`) // mustcompile throws panic instead of error
	if err != nil {
		return nil, apperrors.NewValueError("unable to compile regexp", utils.Caller(), err)
	}
	tempTable = re.ReplaceAllString(tempTable, "") // must not contain digits and '-' chars or sql can throw exception

	_, err = tx.Exec(ctx.Request().Context(),
		`
		create temporary table `+tempTable+` (like url_shortener.url) on commit drop
		`)
	if err != nil {
		return nil, apperrors.NewValueError("unable to create temp table", utils.Caller(), err)
	}

	count, err := tx.CopyFrom(
		ctx.Request().Context(),
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

	queryRows, err := tx.Query(ctx.Request().Context(),
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

	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return nil, apperrors.NewValueError("commit failed", utils.Caller(), err)
	}

	return savedURLs, nil
}
