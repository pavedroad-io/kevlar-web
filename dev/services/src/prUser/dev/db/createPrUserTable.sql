CREATE TABLE IF NOT EXISTS pavedroad.prUser (
    uuid UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    prUser JSONB
);

CREATE INDEX IF NOT EXISTS prUserIdx ON pavedroad.prUser USING GIN (prUser);
