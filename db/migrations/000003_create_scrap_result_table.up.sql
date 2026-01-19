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