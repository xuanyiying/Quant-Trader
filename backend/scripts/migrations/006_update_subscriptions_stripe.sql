-- Migration: Add Stripe identifiers to user_subscriptions
ALTER TABLE user_subscriptions ADD COLUMN IF NOT EXISTS stripe_customer_id VARCHAR(255);
ALTER TABLE user_subscriptions ADD COLUMN IF NOT EXISTS stripe_subscription_id VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_user_subs_stripe_customer ON user_subscriptions(stripe_customer_id);
CREATE INDEX IF NOT EXISTS idx_user_subs_stripe_sub ON user_subscriptions(stripe_subscription_id);
