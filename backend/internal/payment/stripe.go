package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"quant-trader/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/webhook"
	"go.uber.org/zap"
)

var (
	ErrInvalidPriceID      = errors.New("invalid price ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrWebhookVerification  = errors.New("webhook signature verification failed")
)

type Config struct {
	SuccessURL string
	CancelURL  string
}

type StripeService struct {
	db     *pgxpool.Pool
	logger *zap.Logger
	config *Config
}

func NewStripeService(db *pgxpool.Pool, logger *zap.Logger, apiKey string) *StripeService {
	if apiKey == "" {
		logger.Warn("Stripe API key not configured, payment features will be disabled")
	}

	stripe.Key = apiKey

	return &StripeService{
		db:  db,
		logger: logger,
		config: &Config{
			SuccessURL: "https://quant-trader.com/success?session_id={CHECKOUT_SESSION_ID}",
			CancelURL:  "https://quant-trader.com/canceled",
		},
	}
}

func (s *StripeService) CreateCheckoutSession(userID int64, priceID string) (string, error) {
	if userID <= 0 {
		return "", ErrInvalidUserID
	}
	if priceID == "" {
		return "", ErrInvalidPriceID
	}

	s.logger.Info("creating checkout session",
		zap.Int64("user_id", userID),
		zap.String("price_id", priceID))

	if stripe.Key == "" {
		return "", errors.New("Stripe is not configured")
	}

	validPriceIDs, err := s.getValidPriceIDs(context.Background())
	if err != nil {
		s.logger.Error("failed to fetch valid price IDs", zap.Error(err))
		return "", fmt.Errorf("failed to validate price: %w", err)
	}

	if !contains(validPriceIDs, priceID) {
		s.logger.Warn("invalid price ID provided", zap.String("price_id", priceID))
		return "", ErrInvalidPriceID
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(s.config.SuccessURL),
		CancelURL:  stripe.String(s.config.CancelURL),
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", userID))

	sess, err := session.New(params)
	if err != nil {
		s.logger.Error("failed to create checkout session",
			zap.Error(err),
			zap.Int64("user_id", userID))
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	s.logger.Info("checkout session created",
		zap.String("session_id", sess.ID),
		zap.Int64("user_id", userID))

	return sess.URL, nil
}

func (s *StripeService) HandleWebhook(payload []byte, sigHeader string, endpointSecret string) error {
	if len(payload) == 0 {
		return errors.New("empty webhook payload")
	}
	if endpointSecret == "" {
		s.logger.Warn("webhook endpoint secret not configured")
		return errors.New("webhook secret not configured")
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		s.logger.Error("webhook signature verification failed",
			zap.Error(err),
			zap.ByteString("payload", payload))
		return fmt.Errorf("%w: %v", ErrWebhookVerification, err)
	}

	s.logger.Info("received webhook event",
		zap.String("event_type", event.Type),
		zap.String("event_id", event.ID))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutSessionCompleted(ctx, event)

	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)

	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)

	case "invoice.payment_failed":
		return s.handlePaymentFailed(ctx, event)

	default:
		s.logger.Info("unhandled webhook event type", zap.String("event_type", event.Type))
	}

	return nil
}

func (s *StripeService) handleCheckoutSessionCompleted(ctx context.Context, event stripe.Event) error {
	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return fmt.Errorf("failed to unmarshal session: %w", err)
	}

	s.logger.Info("processing checkout session completed",
		zap.String("session_id", sess.ID))

	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")
	fullSess, err := session.Get(sess.ID, params)
	if err != nil {
		return fmt.Errorf("failed to fetch full session: %w", err)
	}

	userIDStr, ok := fullSess.Metadata["user_id"]
	if !ok {
		s.logger.Error("user_id not found in session metadata",
			zap.String("session_id", sess.ID))
		return errors.New("user_id not found in session metadata")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		s.logger.Error("invalid user_id in metadata",
			zap.String("user_id_str", userIDStr),
			zap.Error(err))
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	if fullSess.LineItems == nil || len(fullSess.LineItems.Data) == 0 {
		s.logger.Error("no line items in session",
			zap.String("session_id", sess.ID))
		return errors.New("no line items found in session")
	}

	priceID := fullSess.LineItems.Data[0].Price.ID

	tierName, err := s.mapPriceToTier(ctx, priceID)
	if err != nil {
		s.logger.Error("failed to map price to tier",
			zap.String("price_id", priceID),
			zap.Error(err))
		tierName = model.TierFree
	}

	var customerID, subscriptionID string
	if fullSess.Customer != nil {
		customerID = fullSess.Customer.ID
	}
	if fullSess.Subscription != nil {
		subscriptionID = fullSess.Subscription.ID
	}

	if err := s.UpdateUserTier(ctx, userID, tierName, model.SubscriptionStatusActive, customerID, subscriptionID); err != nil {
		s.logger.Error("failed to update user tier",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	s.logger.Info("subscription created successfully",
		zap.Int64("user_id", userID),
		zap.String("tier", tierName))

	return nil
}

func (s *StripeService) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	s.logger.Info("processing subscription deletion",
		zap.String("subscription_id", sub.ID))

	userIDStr, ok := sub.Metadata["user_id"]
	if !ok {
		s.logger.Warn("user_id not in subscription metadata, looking up by customer ID",
			zap.String("customer_id", sub.Customer.ID))
		return s.handleSubscriptionChangeByStripeID(ctx, sub.Customer.ID, model.TierFree, model.SubscriptionStatusCanceled, sub.ID)
	}

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	if userID <= 0 {
		return ErrInvalidUserID
	}

	return s.UpdateUserTier(ctx, userID, model.TierFree, model.SubscriptionStatusCanceled, sub.Customer.ID, sub.ID)
}

func (s *StripeService) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	s.logger.Info("processing subscription update",
		zap.String("subscription_id", sub.ID),
		zap.String("status", string(sub.Status)))

	userIDStr, ok := sub.Metadata["user_id"]
	if !ok {
		s.logger.Warn("user_id not in subscription metadata, looking up by customer ID")
		tierName, _ := s.mapPriceToTierByPriceID(ctx, sub.Items.Data[0].Price.ID)
		return s.handleSubscriptionChangeByStripeID(ctx, sub.Customer.ID, tierName, string(sub.Status), sub.ID)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		s.logger.Error("invalid user_id in subscription metadata",
			zap.String("user_id_str", userIDStr),
			zap.Error(err))
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	tierName, err := s.mapPriceToTierByPriceID(ctx, sub.Items.Data[0].Price.ID)
	if err != nil {
		s.logger.Error("failed to map price to tier",
			zap.String("price_id", sub.Items.Data[0].Price.ID),
			zap.Error(err))
		tierName = model.TierFree
	}

	return s.UpdateUserTier(ctx, userID, tierName, string(sub.Status), sub.Customer.ID, sub.ID)
}

func (s *StripeService) handlePaymentFailed(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	s.logger.Warn("payment failed",
		zap.String("invoice_id", invoice.ID),
		zap.String("customer_id", invoice.Customer.ID))

	if invoice.Subscription == nil {
		return nil
	}

	return s.handleSubscriptionChangeByStripeID(
		ctx,
		invoice.Customer.ID,
		model.TierFree,
		model.SubscriptionStatusPastDue,
		invoice.Subscription.ID,
	)
}

func (s *StripeService) mapPriceToTier(ctx context.Context, priceID string) (string, error) {
	validPriceIDs, err := s.getValidPriceIDs(ctx)
	if err != nil {
		return "", err
	}

	for _, p := range validPriceIDs {
		if p.PriceID == priceID {
			return p.TierName, nil
		}
	}

	s.logger.Warn("unknown price ID, defaulting to free tier",
		zap.String("price_id", priceID))
	return model.TierFree, nil
}

func (s *StripeService) mapPriceToTierByPriceID(ctx context.Context, priceID string) (string, error) {
	return s.mapPriceToTier(ctx, priceID)
}

type validPrice struct {
	PriceID  string
	TierName string
}

func (s *StripeService) getValidPriceIDs(ctx context.Context) ([]validPrice, error) {
	rows, err := s.db.Query(ctx, `
		SELECT st.name, sp.stripe_price_id 
		FROM subscription_tiers st
		JOIN subscription_prices sp ON st.id = sp.tier_id
		WHERE sp.stripe_price_id IS NOT NULL AND sp.active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []validPrice
	for rows.Next() {
		var p validPrice
		if err := rows.Scan(&p.TierName, &p.PriceID); err != nil {
			s.logger.Warn("failed to scan price row", zap.Error(err))
			continue
		}
		prices = append(prices, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		s.logger.Warn("no valid price IDs found in database")
		return []validPrice{
			{PriceID: "price_pro_default", TierName: model.TierPro},
			{PriceID: "price_enterprise", TierName: model.TierEnterprise},
		}, nil
	}

	return prices, nil
}

func (s *StripeService) handleSubscriptionChangeByStripeID(ctx context.Context, customerID string, tierName string, status string, subID string) error {
	if customerID == "" {
		s.logger.Warn("empty customer ID in subscription change handler")
		return nil
	}

	var userID int64
	err := s.db.QueryRow(ctx, `
		SELECT user_id 
		FROM user_subscriptions 
		WHERE stripe_customer_id = $1`, customerID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("could not find user for stripe customer ID",
				zap.String("customer_id", customerID))
			return nil
		}
		s.logger.Error("database error looking up customer",
			zap.String("customer_id", customerID),
			zap.Error(err))
		return fmt.Errorf("failed to lookup customer: %w", err)
	}

	s.logger.Info("updating subscription by customer ID",
		zap.Int64("user_id", userID),
		zap.String("customer_id", customerID),
		zap.String("tier", tierName),
		zap.String("status", status))

	return s.UpdateUserTier(ctx, userID, tierName, status, customerID, subID)
}

func (s *StripeService) UpdateUserTier(ctx context.Context, userID int64, tierName string, status string, customerID, subID string) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}

	s.logger.Info("updating user tier",
		zap.Int64("user_id", userID),
		zap.String("tier", tierName),
		zap.String("status", status),
		zap.String("customer_id", customerID),
		zap.String("subscription_id", subID))

	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var tierID int64
	err = tx.QueryRow(ctx, "SELECT id FROM subscription_tiers WHERE name = $1", tierName).Scan(&tierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error("tier not found", zap.String("tier_name", tierName))
			return fmt.Errorf("tier not found: %w", err)
		}
		s.logger.Error("failed to get tier ID", zap.Error(err))
		return fmt.Errorf("failed to get tier ID: %w", err)
	}

	expiresAt := time.Now().Add(32 * 24 * time.Hour)
	if status == model.SubscriptionStatusCanceled || status == model.SubscriptionStatusPastDue {
		var existingExpiresAt time.Time
		err = tx.QueryRow(ctx,
			"SELECT expires_at FROM user_subscriptions WHERE user_id = $1",
			userID).Scan(&existingExpiresAt)
		if err == nil && existingExpiresAt.After(expiresAt) {
			expiresAt = existingExpiresAt
		}
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_subscriptions (user_id, tier_id, status, expires_at, stripe_customer_id, stripe_subscription_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			tier_id = $2,
			status = $3,
			expires_at = CASE 
				WHEN $3 IN ('active', 'trialing') THEN GREATEST(user_subscriptions.expires_at, $4)
				ELSE user_subscriptions.expires_at 
			END,
			stripe_customer_id = COALESCE($5, user_subscriptions.stripe_customer_id),
			stripe_subscription_id = COALESCE($6, user_subscriptions.stripe_subscription_id),
			updated_at = NOW()`,
		userID, tierID, status, expiresAt, customerID, subID)

	if err != nil {
		s.logger.Error("failed to update subscription",
			zap.Error(err),
			zap.Int64("user_id", userID))
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("user subscription updated successfully",
		zap.Int64("user_id", userID),
		zap.String("tier", tierName),
		zap.String("status", status))

	return nil
}

func contains(slice []validPrice, priceID string) bool {
	for _, p := range slice {
		if p.PriceID == priceID {
			return true
		}
	}
	return false
}
