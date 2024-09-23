DROP INDEX IF EXISTS idx_payments_merchant_id;
DROP INDEX IF EXISTS idx_refunds_payment_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

ALTER TABLE payments DROP CONSTRAINT IF EXISTS chk_payment_status;
ALTER TABLE refunds DROP CONSTRAINT IF EXISTS chk_refund_status;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS chk_currency;
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS chk_merchant_email;
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_user_email;