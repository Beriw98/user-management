CREATE TABLE IF NOT EXISTS users {
    id              VARCHAR(36) PRIMARY KEY,
--     name            VARCHAR(50),
--     surname         VARCHAR(50),
    password        VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
}