-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(256) NOT NULL,
    password VARCHAR(256) NOT NULL
);

CREATE TABLE profile (
    user_id BIGINT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    city VARCHAR(35) NOT NULL,
    age TINYINT NOT NULL,
    gender TINYINT NOT NULL
);

CREATE TABLE profile_interest (
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY(user_id, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE profile_interest;
DROP TABLE profile;
DROP TABLE auth;
-- +goose StatementEnd
