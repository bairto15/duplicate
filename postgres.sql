CREATE TABLE IF NOT EXISTS conn_log (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    ip_addr VARCHAR(15),
    ts TIMESTAMP
);

TRUNCATE TABLE conn_log;

INSERT INTO conn_log (user_id, ip_addr, ts)
SELECT
  floor(random() * 10000) AS user_id,
  (
    '127.0.' ||
    (floor(random() * 255) + 1)::text || '.' ||
    (floor(random() * 255) + 1)::text
  ) AS ip_addr,
  current_timestamp - (random() * interval '365 days')
FROM generate_series(1, 1000000);
