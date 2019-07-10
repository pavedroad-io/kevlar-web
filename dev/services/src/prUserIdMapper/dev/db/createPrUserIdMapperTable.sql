CREATE TABLE IF NOT EXISTS pavedroad.prUserIdMapper (
    apiVersion  VARCHAR(254) NOT NULL,
    kind VARCHAR(25) NOT NULL,
    objVersion CHAR(10) NOT NULL,
    credential VARCHAR(254) NOT NULL PRIMARY KEY,
    userUUID VARCHAR(254),
    loginCount INT,
    created TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ,
    active CHAR(5) NOT NULL
);
