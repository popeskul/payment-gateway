CREATE INDEX IF NOT EXISTS idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX IF NOT EXISTS idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

ALTER TABLE payments ADD CONSTRAINT chk_payment_status CHECK (status IN ('pending', 'completed', 'failed'));
ALTER TABLE refunds ADD CONSTRAINT chk_refund_status CHECK (status IN ('pending', 'completed', 'failed'));

ALTER TABLE payments ADD CONSTRAINT chk_currency CHECK (currency ~ '^[A-Z]{3}$');

ALTER TABLE merchants ADD CONSTRAINT chk_merchant_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$');
ALTER TABLE users ADD CONSTRAINT chk_user_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$');
