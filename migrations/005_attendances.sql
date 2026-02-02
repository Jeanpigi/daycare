CREATETABLEIFNOTEXISTSattendances(
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  child_id BIGINT UNSIGNED NOT NULL,
  checked_in_at DATETIME NOT NULL,
  checked_out_at DATETIME NULL,
  minutes INT NULL,
  pricing_config_id BIGINT UNSIGNED NULL,
  promo_id BIGINT UNSIGNED NULL,
  gross_amount DECIMAL(12,2)NULL,
  discount_amount DECIMAL(12,2)NULL,
  net_amount DECIMAL(12,2)NULL,
  currency CHAR(3)NOT NULL DEFAULT 'COP',
  created_by BIGINT UNSIGNED NULL,
  closed_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(
    id
  ),
  KEY idx_att_child_in(
    child_id,
    checked_in_at
  ),
  KEY idx_att_child_open(
    child_id,
    checked_out_at
  ),
  CONSTRAINT fk_att_child FOREIGN KEY(
    child_id
  ) REFERENCES children(id)
    ON DELETE RESTRICT
      ON UPDATE CASCADE CONSTRAINT fk_att_pricing FOREIGN KEY(
        pricing_config_id
      ) REFERENCES settings_pricing(id)
        ON DELETE
      SET NULL
        ON UPDATE CASCADE CONSTRAINT fk_att_promo FOREIGN KEY(
          promo_id
        ) REFERENCES promotions(id)
          ON DELETE
        SET NULL
          ON UPDATE CASCADE
)ENGINE = InnoDBDEFAULTCHARSET = utf8mb4;
