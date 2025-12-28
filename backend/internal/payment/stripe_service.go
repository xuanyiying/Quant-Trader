package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"quant-trader/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/webhook"
	"go.uber.org/zap"
)

type StripeService struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewStripeService(db *pgxpool.Pool, logger *zap.Logger, apiKey string) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		db:     db,
		logger: logger,
	}
}

func (s *StripeService) CreateCheckoutSession(userID int64, priceID string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("https://quant-trader.com/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://quant-trader.com/canceled"),
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", userID))

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}

func (s *StripeService) HandleWebhook(payload []byte, sigHeader string, endpointSecret string) error {
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		s.logger.Error("webhook signature verification failed", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch event.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return fmt.Errorf("failed to unmarshal session: %w", err)
		}

		// Expand line items to get the price ID
		params := &stripe.CheckoutSessionParams{}
		params.AddExpand("line_items")
		fullSess, err := session.Get(sess.ID, params)
		if err != nil {
			return fmt.Errorf("failed to fetch full session: %w", err)
		}

		userIDStr, ok := fullSess.Metadata["user_id"]
		if !ok {
			return fmt.Errorf("user_id not found in session metadata")
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user_id in metadata: %w", err)
		}

		if fullSess.LineItems == nil || len(fullSess.LineItems.Data) == 0 {
			return fmt.Errorf("no line items found in session")
		}

		priceID := fullSess.LineItems.Data[0].Price.ID
		tierName := s.mapPriceToTier(priceID)

		return s.UpdateUserTier(ctx, userID, tierName, model.SubscriptionStatusActive, fullSess.Customer.ID, fullSess.Subscription.ID)

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return fmt.Errorf("failed to unmarshal subscription: %w", err)
		}

		userIDStr, ok := sub.Metadata["user_id"]
		if !ok {
			// If not in metadata, we might need to look up by stripe customer ID
			return s.handleSubscriptionChangeByStripeID(ctx, sub.Customer.ID, model.TierFree, model.SubscriptionStatusCanceled, sub.ID)
		}
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		return s.UpdateUserTier(ctx, userID, model.TierFree, model.SubscriptionStatusCanceled, sub.Customer.ID, sub.ID)

	case "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return fmt.Errorf("failed to unmarshal subscription: %w", err)
		}

		userIDStr, ok := sub.Metadata["user_id"]
		if !ok {
			return s.handleSubscriptionChangeByStripeID(ctx, sub.Customer.ID, s.mapPriceToTier(sub.Items.Data[0].Price.ID), string(sub.Status), sub.ID)
		}
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		tierName := s.mapPriceToTier(sub.Items.Data[0].Price.ID)
		return s.UpdateUserTier(ctx, userID, tierName, string(sub.Status), sub.Customer.ID, sub.ID)
	}

	return nil
}

func (s *StripeService) mapPriceToTier(priceID string) string {
	// In a real app, these should come from config or DB
	priceToTier := map[string]string{
		"price_pro_default":    model.TierPro,
		"price_enterprise":     model.TierEnterprise,
		"price_123_pro":        model.TierPro,
		"price_456_enterprise": model.TierEnterprise,
	}

	if tier, ok := priceToTier[priceID]; ok {
		return tier
	}
	return model.TierFree
}

func (s *StripeService) handleSubscriptionChangeByStripeID(ctx context.Context, customerID string, tierName string, status string, subID string) error {
	var userID int64
	err := s.db.QueryRow(ctx, "SELECT user_id FROM user_subscriptions WHERE stripe_customer_id = $1", customerID).Scan(&userID)
	if err != nil {
		s.logger.Warn("could not find user for stripe customer id", zap.String("customer_id", customerID))
		return nil // Not necessarily an error we want to retry
	}
	return s.UpdateUserTier(ctx, userID, tierName, status, customerID, subID)
}

func (s *StripeService) UpdateUserTier(ctx context.Context, userID int64, tierName string, status string, customerID, subID string) error {
	// 1. Get tier ID
	var tierID int64
	err := s.db.QueryRow(ctx, "SELECT id FROM subscription_tiers WHERE name = $1", tierName).Scan(&tierID)
	if err != nil {
		return fmt.Errorf("tier not found: %w", err)
	}

	// 2. Upsert subscription
	_, err = s.db.Exec(ctx, fmt.Sprintf(`
		INSERT INTO user_subscriptions (user_id, tier_id, status, expires_at, stripe_customer_id, stripe_subscription_id)
		VALUES ($1, $2, $3, NOW() + INTERVAL '32 days', $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			tier_id = $2,
			status = $3,
			expires_at = CASE 
				WHEN $3 = '%s' THEN NOW() + INTERVAL '32 days'
				ELSE user_subscriptions.expires_at 
			END,
			stripe_customer_id = $4,
			stripe_subscription_id = $5,
			updated_at = NOW()`, model.SubscriptionStatusActive),
		userID, tierID, status, customerID, subID)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	s.logger.Info("user subscription updated",
		zap.Int64("user_id", userID),
		zap.String("tier", tierName),
		zap.String("status", status))
	return nil
}
