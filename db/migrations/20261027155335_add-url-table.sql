-- migrate:up
CREATE TABLE url_code (
    id      bigserial PRIMARY KEY,
    code    TEXT NOT NULL UNIQUE,
    url     TEXT NOT NULL,
    created_at    TIMESTAMPTZ not null default now()
);


-- migrate:down
DROP TABLE url_code;
