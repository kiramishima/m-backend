CREATE TABLE IF NOT EXISTS users_profile (
    user_id BIGINT NOT NULL UNIQUE,
    username CHAR(25) NOT NULL CHECK ( CHAR_LENGTH(username) >= 6 ),
    photo TINYTEXT NOT NULL DEFAULT 'default_profile.png',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT FK_User2Profile FOREIGN KEY (user_id) REFERENCES users(id)
         ON UPDATE CASCADE
) ENGINE=INNODB;