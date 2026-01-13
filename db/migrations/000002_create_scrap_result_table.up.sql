CREATE TABLE scrap_result
(
    id              BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    scrap_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scrap_target_id BIGINT      NOT NULL REFERENCES scrap_target (id),
    title           TEXT        NOT NULL,
    article_url     TEXT        NOT NULL,
    image_url       TEXT,
    content         TEXT
);

CREATE INDEX idx_scrap_result_scrap_target_id ON scrap_result (scrap_target_id);