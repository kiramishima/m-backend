CREATE TABLE IF NOT EXISTS bonds (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(50) NOT NULL UNIQUE CHECK (uuid REGEXP '^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$'),
    name CHAR(40) NOT NULL UNIQUE CHECK (CHAR_LENGTH(name) >= 3),
    number INT NOT NULL CHECK(number >= 1 AND number <= 10000),
    price DECIMAL(13, 4) NOT NULL CHECK(price > 0),
    currency_id INT NOT NULL DEFAULT 1,
    created_by BIGINT NOT NULL,
    status ENUM('on_hold', 'on_sell') NOT NULL DEFAULT 'on_hold' CHECK ( status IN ('on_hold', 'on_sell')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT FK_UserBond FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_CurrencyBond FOREIGN KEY (currency_id) REFERENCES currencies(id)
) ENGINE=INNODB;

