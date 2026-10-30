-- migrate:up

CREATE TABLE IF NOT EXISTS click_events (
    id          bigint GENERATED ALWAYS AS IDENTITY,
    link_id     bigint      NOT NULL,
    clicked_at  timestamptz NOT NULL DEFAULT now(),
    user_agent  text        NOT NULL DEFAULT '',
    browser     text        NOT NULL DEFAULT '',
    os          text        NOT NULL DEFAULT '',
    device_type text        NOT NULL DEFAULT '',
    country     text        NOT NULL DEFAULT '',
    ip_hash     bytea,
    is_bot      boolean     NOT NULL DEFAULT false,
    status_code smallint    NOT NULL,

    CONSTRAINT pk_click_events PRIMARY KEY (id),
    CONSTRAINT chck_click_events_status_code
        CHECK (status_code IN (301, 303, 302, 410)),
    CONSTRAINT chck_click_events_user_agent_length
        CHECK (char_length(user_agent) <= 512)
);

CREATE INDEX IF NOT EXISTS idx_click_events_link_id_clicked_at
    ON click_events (link_id, clicked_at);

-- migrate:down
DROP INDEX IF EXISTS idx_click_events_link_id_clicked_at;
DROP TABLE IF EXISTS click_events;
