package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Database driver import is required
)

type ScrapTarget struct {
	ID                     int64
	URL                    string
	baseGoquerySelector    string
	itemGoquerySelector    string
	imageGoquerySelector   string
	articleGoquerySelector string
}

type ScrapResult struct {
	Title      string
	ArticleUrl string
	ImageUrl   sql.NullString
	Content    sql.NullString
}

func NewDatabaseConnectionPool(config Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseUrl)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnectionPool(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Warning: failed to close DB: %v", err)
	}
}

func GetScrapTargets(db *sql.DB) ([]ScrapTarget, error) {
	query := `SELECT id,
				url, 
				base_goquery_selector,
				item_goquery_selector,
				image_goquery_selector,
				article_goquery_selector
		   FROM scrap_target`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()

		if err != nil {
			panic(err)
		}
	}(rows)

	var scrapTargets []ScrapTarget

	for rows.Next() {
		var scrapTarget ScrapTarget
		err := rows.Scan(
			&scrapTarget.ID,
			&scrapTarget.URL,
			&scrapTarget.baseGoquerySelector,
			&scrapTarget.itemGoquerySelector,
			&scrapTarget.imageGoquerySelector,
			&scrapTarget.articleGoquerySelector,
		)

		if err != nil {
			return nil, err
		}

		scrapTargets = append(scrapTargets, scrapTarget)
	}

	return scrapTargets, nil
}

func InsertScrapProcess(db *sql.DB) (int64, error) {
	query := `INSERT INTO scrap_process DEFAULT VALUES RETURNING id`
	var insertedRowId int64
	err := db.QueryRow(query).Scan(&insertedRowId)

	if err != nil {
		return 0, err
	}

	return insertedRowId, nil
}

func FinishScrapProcess(db *sql.DB, processId int64) error {
	query := `UPDATE scrap_process SET finished_at = now(), status = 'FINISHED' WHERE id = $1`
	_, err := db.Exec(query, processId)
	return err
}

func InsertScrapResult(db *sql.DB, processId int64, scrapTarget ScrapTarget, title string, articleUrl string, imageUrl string) (int64, error) {
	query := `INSERT INTO scrap_result (scrap_process_id, scrap_target_id, title, article_url, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var insertedRowId int64
	err := db.QueryRow(query, processId, scrapTarget.ID, title, articleUrl, imageUrl).Scan(&insertedRowId)

	if err != nil {
		return 0, err
	}

	return insertedRowId, nil
}

func UpdateScrapContent(db *sql.DB, id int64, content string) error {
	query := `UPDATE scrap_result SET content = $1 WHERE id = $2`
	_, err := db.Exec(query, content, id)

	if err != nil {
		return err
	}

	return nil
}

func GetLatestScrapResult(db *sql.DB, scrapTargetId int64) ([]ScrapResult, error) {
	query := `SELECT title, article_url, image_url, content FROM latest_scrap_results WHERE scrap_target_id = $1`
	rows, err := db.Query(query, scrapTargetId)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()

		if err != nil {
			panic(err)
		}
	}(rows)

	var scrapResults []ScrapResult

	for rows.Next() {
		var scrapResult ScrapResult
		err := rows.Scan(
			&scrapResult.Title,
			&scrapResult.ArticleUrl,
			&scrapResult.ImageUrl,
			&scrapResult.Content,
		)

		if err != nil {
			return nil, err
		}

		scrapResults = append(scrapResults, scrapResult)
	}

	return scrapResults, nil
}
