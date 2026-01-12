CREATE TABLE scrap_target
(
    id                       BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    url                      TEXT NOT NULL,
    base_goquery_selector    TEXT NOT NULL,
    item_goquery_selector    TEXT NOT NULL,
    image_goquery_selector   TEXT NOT NULL,
    article_goquery_selector TEXT NOT NULL,
    UNIQUE (url, base_goquery_selector)
);