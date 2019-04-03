CREATE TABLE IF NOT EXISTS kevlarweb.prToken (
    uid UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    prToken JSONB
);

CREATE INDEX prTokenIdx ON kevlarweb.prToken USING GIN (prToken);
