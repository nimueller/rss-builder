CREATE TABLE scrap_result
(
    id               BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    scrap_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scrap_process_id BIGINT      NOT NULL REFERENCES scrap_process (id),
    scrap_target_id  BIGINT      NOT NULL REFERENCES scrap_target (id),
    title            TEXT        NOT NULL,
    article_url      TEXT        NOT NULL,
    image_url        TEXT,
    content          TEXT
);

CREATE INDEX idx_scrap_result_scrap_process_id ON scrap_result (scrap_process_id);
CREATE INDEX idx_scrap_result_process_target ON scrap_result (scrap_process_id, scrap_target_id);

CREATE VIEW latest_scrap_results AS
WITH latest_process AS (SELECT *
                        FROM scrap_process
                        ORDER BY created_at DESC
                        LIMIT 1)
SELECT sr.scrap_target_id aS "scrap_target_id",
       sr.title           AS "title",
       sr.article_url     AS "article_url",
       sr.image_url       AS "image_url",
       sr.content         AS "content"
FROM latest_process lp
         CROSS JOIN LATERAL (SELECT * FROM scrap_result sr WHERE lp.id = sr.scrap_process_id) sr;