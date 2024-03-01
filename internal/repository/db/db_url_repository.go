package db

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/apperr"
)

//go:embed queries/insert_url_and_return.sql
var insertURLAndReturn string

//go:embed queries/select_url_by_id.sql
var selectURLByID string

//go:embed queries/select_all_urls.sql
var selectAllURLs string

//go:embed queries/select_all_urls_by_userid.sql
var selectAllURLsByUserID string

//go:embed queries/delete_all_urls.sql
var deleteAllURLs string

//go:embed queries/block_urls_by_userid_and_urlids.sql
var blockURLsByUserIDAndURLsIDs string

//go:embed queries/set_true_deleted_to_urls_by_userid_and_urlsids.sql
var setDeletedByUserIDandURLsIDs string

//go:embed queries/create_tmp_table_like_url.sql
var createTmpTableLikeURL string

//go:embed queries/upsert_and_return_urls_from_tmp_table.sql
var upsertAndReturnURLsFromTmpTable string

//go:embed queries/select_stats.sql
var selectStats string

// PostgresURLRepository represents a PostgreSQL implementation of the URLRepository interface.
type PostgresURLRepository struct {
	PostgresPool *PostgresPool
	logger       *zap.Logger
}

// NewPostgresURLRepository returns a new instance of PostgresURLRepository.
//
// It takes a PostgresPool pointer and a zap.Logger pointer as parameters and returns a PostgresURLRepository pointer.
func NewPostgresURLRepository(postgresPool *PostgresPool, logger *zap.Logger) *PostgresURLRepository {
	return &PostgresURLRepository{
		PostgresPool: postgresPool,
		logger:       logger,
	}
}

// Ping pings the PostgresURLRepository.
func (r *PostgresURLRepository) Ping(ctx context.Context) error {
	return r.PostgresPool.db.Ping(ctx)
}

// DeleteURLByUserID deletes from PostgreSQL DB URLs by user ID.
//
// Performed via select for update block within batched transaction
func (r *PostgresURLRepository) DeleteURLByUserID(ctx context.Context, userID string, shortURLs string) error {
	r.logger.Info("DeleteURLByUserID", zap.String("userID", userID), zap.String("shortURLs", shortURLs))
	tx, err := r.PostgresPool.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return apperr.NewValueError("unable to start transaction", apperr.Caller(), err)
	}
	defer tx.Rollback(ctx)

	block, err := tx.Prepare(ctx, "block", blockURLsByUserIDAndURLsIDs)
	if err != nil {
		return apperr.NewValueError("unable to prepare query", apperr.Caller(), err)
	}

	update, err := tx.Prepare(ctx, "update", setDeletedByUserIDandURLsIDs)
	if err != nil {
		return apperr.NewValueError("unable to prepare query", apperr.Caller(), err)
	}

	batch := &pgx.Batch{}
	batch.Queue(block.Name, userID, shortURLs)
	batch.Queue(update.Name, userID, shortURLs)
	result := tx.SendBatch(ctx, batch)

	err = result.Close()
	if err != nil {
		return apperr.NewValueError("close failed", apperr.Caller(), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperr.NewValueError("commit failed", apperr.Caller(), err)
	}

	return nil
}

// SelectAllByUserID retrieves from PostgreSQL DB URLs by user ID.
func (r *PostgresURLRepository) SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error) {
	queryRows, err := r.PostgresPool.db.Query(ctx, selectAllURLsByUserID, userID)
	if err != nil {
		return nil, apperr.NewValueError("query failed", apperr.Caller(), err)
	}
	defer queryRows.Close()

	urls, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.URL])
	if err != nil {
		return nil, apperr.NewValueError("unable to collect rows", apperr.Caller(), err)
	}

	if len(urls) == 0 {
		return nil, apperr.NewValueError(fmt.Sprintf("urls not found by user %s", userID), apperr.Caller(), urlErr.ErrURLNotFound)
	}

	return urls, nil
}

// Insert inserts to PostgreSQL DB URL.
func (r *PostgresURLRepository) Insert(ctx context.Context, url model.URL) (*model.URL, error) {
	var savedURL model.URL
	err := r.PostgresPool.db.QueryRow(ctx, insertURLAndReturn,
		url.ID, url.Original, url.Shortened, url.UserID, url.DeletedFlag).
		Scan(&savedURL.ID, &savedURL.Original, &savedURL.Shortened, &savedURL.UserID, &savedURL.DeletedFlag)
	if err != nil {
		return nil, apperr.NewValueError("query failed", apperr.Caller(), err)
	}

	return &savedURL, nil
}

// SelectByID retrieves URL from PostgreSQL DB by ID.
func (r *PostgresURLRepository) SelectByID(ctx context.Context, key string) (*model.URL, error) {
	var url model.URL
	err := r.PostgresPool.db.QueryRow(ctx, selectURLByID, key).
		Scan(&url.ID, &url.Original, &url.Shortened, &url.CorrelationID, &url.UserID, &url.DeletedFlag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = apperr.NewValueError("url not found", apperr.Caller(), urlErr.ErrURLNotFound)
		} else {
			err = apperr.NewValueError("query failed", apperr.Caller(), err)
		}
		return nil, err
	}

	return &url, nil
}

// SelectAll retrieves all URL from PostgreSQL DB.
func (r *PostgresURLRepository) SelectAll(ctx context.Context) ([]model.URL, error) {
	queryRows, err := r.PostgresPool.db.Query(ctx, selectAllURLs)
	if err != nil {
		return nil, apperr.NewValueError("query failed", apperr.Caller(), err)
	}
	defer queryRows.Close()

	urls, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.URL])
	if err != nil {
		return nil, apperr.NewValueError("unable to collect rows", apperr.Caller(), err)
	}

	return urls, nil
}

// SelectStats retrieves stats from PostgreSQL DB.
func (r *PostgresURLRepository) SelectStats(ctx context.Context) (*model.URLStats, error) {
	var urlStats model.URLStats
	err := r.PostgresPool.db.QueryRow(ctx, selectStats).
		Scan(&urlStats.Urls, &urlStats.Users)
	if err != nil {
			err = apperr.NewValueError("query failed", apperr.Caller(), err)
		return nil, err
	}

	return &urlStats, nil
}

// DeleteAll deletes all URL from PostgreSQL DB.
func (r *PostgresURLRepository) DeleteAll(ctx context.Context) error {
	_, err := r.PostgresPool.db.Exec(ctx, deleteAllURLs)
	if err != nil {
		return apperr.NewValueError("query failed", apperr.Caller(), err)
	}

	return nil
}

// InsertAllOrUpdate upserts URLs to PostgreSQL DB.
//
// performed in a single transaction with copy protocol and temp table
func (r *PostgresURLRepository) InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	tx, err := r.PostgresPool.db.Begin(ctx)
	if err != nil {
		return nil, apperr.NewValueError("unable to start transaction", apperr.Caller(), err)
	}
	defer tx.Rollback(ctx)

	rows := make([][]interface{}, len(urls))
	for i, url := range urls {
		row := []interface{}{url.ID, url.Original, url.Shortened, url.CorrelationID, url.UserID, url.DeletedFlag}
		rows[i] = row
	}

	tempTable := uuid.New().String()
	re := regexp.MustCompile(`\d|-`)
	tempTable = re.ReplaceAllString(tempTable, "")

	createTmpTableQuery := fmt.Sprintf(createTmpTableLikeURL, tempTable)
	_, err = tx.Exec(ctx, createTmpTableQuery)
	if err != nil {
		return nil, apperr.NewValueError("unable to create temp table", apperr.Caller(), err)
	}

	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pg_temp", tempTable},
		[]string{"id", "original_url", "short_url", "correlation_id", "user_id", "deleted_flag"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, apperr.NewValueError("copy from failed", apperr.Caller(), err)
	}
	if count != int64(len(urls)) {
		return nil, apperr.NewValueError("not all rows were inserted", apperr.Caller(), err)
	}

	upsertFromTmpTableQuery := fmt.Sprintf(upsertAndReturnURLsFromTmpTable, tempTable)
	queryRows, err := tx.Query(ctx, upsertFromTmpTableQuery)
	if err != nil {
		return nil, apperr.NewValueError("unable to upsert batch", apperr.Caller(), err)
	}
	defer queryRows.Close()

	savedURLs, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.URL])
	if err != nil {
		return nil, apperr.NewValueError("unable to collect rows", apperr.Caller(), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, apperr.NewValueError("commit failed", apperr.Caller(), err)
	}

	return savedURLs, nil
}
