package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ScrapTarget struct {
	URL                    string
	baseGoquerySelector    string
	itemGoquerySelector    string
	imageGoquerySelector   string
	articleGoquerySelector string
}

func NewDatabaseConnection() (*sql.DB, error) {
	url, urlSet := os.LookupEnv("DATABASE_URL")

	if urlSet == false {
		url = "postgres://postgres:postgres@localhost:5432/postgres"
	}

	db, err := sql.Open("pgx", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Warning: failed to close DB: %v", err)
	}
}

func GetScrapTargets(db *sql.DB) ([]ScrapTarget, error) {
	rows, err := db.Query(`SELECT url,
										base_goquery_selector,
										item_goquery_selector,
										image_goquery_selector,
										article_goquery_selector
								   FROM scrap_target`)

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
		err := rows.Scan(&scrapTarget.URL, &scrapTarget.baseGoquerySelector, &scrapTarget.itemGoquerySelector, &scrapTarget.imageGoquerySelector, &scrapTarget.articleGoquerySelector)

		if err != nil {
			return nil, err
		}

		scrapTargets = append(scrapTargets, scrapTarget)
	}

	return scrapTargets, nil
}
