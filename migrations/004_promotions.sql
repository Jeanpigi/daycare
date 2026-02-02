CREATETABLEIFNOTEXISTSpromotions(
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(120) NOT NULL,
  rule_type VARCHAR(30) NOT NULL,
  promo_price DECIMAL(12,2) NOT NULL,
  min_days INT NULL,
  min_minutes INT NULL,
  starts_at DATETIME NULL,
  ends_at DATETIME NULL,
  priority INTNOT NULL DEFAULT 100,
  active TINYINT(1) NOT NULL DEFAULT 1,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(
    id
  ),
  KEY idx_promos_active(
    active,
    priority
  ),
  CONSTRAINT fk_promos_created_by FOREIGN KEY(
    created_by
  ) REFERENCES users(id)
)ENGINE = InnoDBDEFAULTCHARSET = utf8mb4;
