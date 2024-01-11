CREATE TABLE IF NOT EXISTS transactions (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    seller_id BIGINT NOT NULL,
    buyer_id BIGINT NOT NULL,
    bond_id BIGINT NOT NULL,
    total_acquired BIGINT NOT NULL,
    status TINYINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT FK_UserSellerTransaction FOREIGN KEY(seller_id) REFERENCES users(id),
    CONSTRAINT FK_UserBuyerTransaction FOREIGN KEY(buyer_id) REFERENCES users(id),
    CONSTRAINT FK_BondTransaction FOREIGN KEY(bond_id) REFERENCES bonds(id)
) ENGINE=INNODB;