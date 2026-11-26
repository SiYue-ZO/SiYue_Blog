-- Preserve existing totals. Only new browser page views create events.
CREATE TABLE post_view_events (
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    event_id   UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, event_id)
);
