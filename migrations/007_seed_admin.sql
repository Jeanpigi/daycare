-- bcrypt hash para "admin123"
INSERT INTO users(
  name,
  email,
  password_hash,
  role
)
VALUES(
  'Admin',
  'admin@daycare.local',
  '$2a$10$CwTycUXWue0Thq9StjUM0uJ8u2QW2X8nQ2eC3x4hJx2xR7q4w8W2S',
  'ADMIN'
)
  ON DUPLICATEKEYUPDATEemail = email;
-- Pricing base por defecto (20000 COP), activo
INSERT INTO settings_pricing(
  standard_price,
  currency,
  active,
  created_by
)SELECT
  20000,
  'COP',
  1,
  id
FROM
  users
WHERE
  email = 'admin@daycare.local'
    ON DUPLICATEKEYUPDATEactive = active;
-- Promo default: LOYALTY_MONTH => 15000 si 2 días o 300 mins, prioridad 10
INSERT INTO promotions(
  name,
  rule_type,
  promo_price,
  min_days,
  min_minutes,
  priority,
  active,
  created_by
)SELECT
  'Promo fidelidad',
  'LOYALTY_MONTH',
  15000,
  2,
  300,
  10,
  1,
  id
FROM
  users
WHERE
  email = 'admin@daycare.local'
    ON DUPLICATEKEYUPDATEactive = active;
