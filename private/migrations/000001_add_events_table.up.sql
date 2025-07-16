-- migrations/001_create_logs_table.sql
CREATE TABLE logs
(
    _timestamp    DateTime,
    _namespace    LowCardinality(String),
    host          LowCardinality(String),
    service       LowCardinality(String),
    level         LowCardinality(String),
    user_id       LowCardinality(String),
    session_id    Nullable(String),
    trace_id      Nullable(String),
    _source       String,
    string_names  Array(String),
    string_values Array(String),
    int_names     Array(String),
    int_values    Array(Int64),
    float_names   Array(String),
    float_values  Array(Float64),
    bool_names    Array(String),
    bool_values   Array(String)
) ENGINE = MergeTree
      PARTITION BY toYYYYMM(_timestamp)
      ORDER BY (_namespace, service, _timestamp)
      TTL _timestamp + INTERVAL 30 DAY
      SETTINGS index_granularity = 8192;
