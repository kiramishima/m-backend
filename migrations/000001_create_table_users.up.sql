CREATE TABLE IF NOT EXISTS users (
     id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
     email VARCHAR(120) NOT NULL UNIQUE CHECK (email REGEXP '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
     password VARCHAR(90) NOT NULL CHECK(password != "" AND CHAR_LENGTH(password) >= 6),
     email_verified_at TIMESTAMP NULL,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP,
     deleted_at TIMESTAMP
) ENGINE=INNODB;