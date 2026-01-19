CREATE TYPE scrap_process_status AS ENUM ('RUNNING', 'FINISHED');

CREATE TABLE scrap_process
(
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at  TIMESTAMPTZ          NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ                   DEFAULT NULL,
    status      scrap_process_status NOT NULL DEFAULT 'RUNNING',
    CONSTRAINT chk_finished_after_created CHECK (finished_at IS NULL OR finished_at >= created_at),
    CONSTRAINT chk_finished_status CHECK (finished_at IS NULL OR status = 'FINISHED')
);

CREATE INDEX idx_scrap_process_created_at ON scrap_process (created_at DESC);