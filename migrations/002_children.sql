CREATETABLEIFNOTEXISTSchildren(
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  document_number VARCHAR(32)NOT NULL,
  first_name VARCHAR(80)NOT NULL,
  last_name VARCHAR(80) NOT NULL,
  guardian_name VARCHAR(120) NOT NULL,
  guardian_phone VARCHAR(30) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id),
  UNIQUE KEY uq_children_document(
    document_number
  )
)ENGINE = InnoDBDEFAULTCHARSET = utf8mb4;
