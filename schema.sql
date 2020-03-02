--
-- DSN: postgres://iktest@localhost/iktest?sslmode=disable
--

CREATE TABLE IF NOT EXISTS ik_table (
    uuid_v1      uuid,
    uuid_v4      uuid,
    snowflake_id bigint,
    uuid_v1_as_string   VARCHAR(16)
    ctime TIMESTAMP WITH TIME ZONE default now()
);
CREATE UNIQUE INDEX IF NOT EXISTS uuid_v1 ON ik_table (uuid_v1);         -- uuid_ops
CREATE UNIQUE INDEX IF NOT EXISTS uuid_v4 ON ik_table (uuid_v4);         -- uuid_ops
CREATE UNIQUE INDEX IF NOT EXISTS snowflake ON ik_table (snowflake_id);  -- int8_ops
