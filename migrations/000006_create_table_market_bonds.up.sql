CREATE TABLE IF NOT EXISTS market_bonds (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    bond_id BIGINT NOT NULL,
    available INT NOT NULL DEFAULT 0,
    status ENUM('available', 'bought') NOT NULL DEFAULT 'available' CHECK ( status IN ('available', 'bought')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT FK_MarketBondsBond FOREIGN KEY (bond_id) REFERENCES bonds(id)
) ENGINE=INNODB;

