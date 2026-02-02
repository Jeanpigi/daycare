CREATETABLEIFNOTEXISTSaudit_log(
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  actor_user_id BIGINT UNSIGNED NOT NULL,
  action VARCHAR(50)NOT NULL,
  entity_type VARCHAR(50)NOT NULL,
  entity_id BIGINT UNSIGNED NULL,
  before_json JSON NULL,
  after_json JSON NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(
    id
  ),
  KEY idx_audit_actor(
    actor_user_id,
    created_at
  ),
  CONSTRAINT fk_audit_actor FOREIGN KEY(
    actor_user_id
  ) REFERENCES users(id)
)ENGINE = InnoDBDEFAULTCHARSET = utf8mb4;
