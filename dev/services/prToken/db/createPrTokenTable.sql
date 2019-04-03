CREATE TABLE IF NOT EXISTS prToken (
    uid UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    prToken JSONB
);

CREATE INDEX prTokenIdx ON prToken USING GIN (prToken);
