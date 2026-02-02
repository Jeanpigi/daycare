CREATETABLEIFNOTEXISTSsettings_pricing(
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  standard_price DECIMAL(12,2) NOT NULL,
  currency CHAR(3) NOT NULL DEFAULT 'COP',
  active TINYINT(1) NOT NULL DEFAULT 1,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(
    id
  ),
  KEY idx_pricing_active(
    active
  ),
  CONSTRAINT fk_pricing_created_by FOREIGN KEY(
    created_by
  ) REFERENCES users(id)
)ENGINE = InnoDBDEFAULTCHARSET = utf8mb4;
