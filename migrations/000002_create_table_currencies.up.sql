CREATE TABLE IF NOT EXISTS currencies (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    currency VARCHAR(70) NOT NULL UNIQUE CHECK(currency != ""),
    currency_short_name CHAR(4) NOT NULL UNIQUE CHECK(currency_short_name != ""),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
) ENGINE=INNODB;