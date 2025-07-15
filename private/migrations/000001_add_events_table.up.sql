CREATE TABLE events
(
    id DEFAULT generateUUIDv4(),
    data JSON,
    created_at DateTime64(3, 'UTC') DEFAULT now(),
    updated_at DateTime64(3, 'UTC') DEFAULT now()
)
ENGINE = MergeTree

ORDER BY id;
