package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"go.uber.org/zap"
)

// Notifier 告警通知接口
type Notifier interface {
	Send(ctx context.Context, alert Alert) error
}

// Alert 告警信息
type Alert struct {
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Level       string            `json:"level"` // info, warning, error, critical
	Symbol      string            `json:"symbol,omitempty"`
	Price       float64           `json:"price,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// DingTalkNotifier 钉钉通知
type DingTalkNotifier struct {
	webhookURL string
	logger     *zap.Logger
}

// NewDingTalkNotifier 创建钉钉通知器
func NewDingTalkNotifier(webhookURL string, logger *zap.Logger) *DingTalkNotifier {
	return &DingTalkNotifier{
		webhookURL: webhookURL,
		logger:     logger,
	}
}

// Send 发送钉钉通知
func (d *DingTalkNotifier) Send(ctx context.Context, alert Alert) error {
	content := fmt.Sprintf("**%s**\n\n%s\n\n时间: %s",
		alert.Title,
		alert.Message,
		alert.Timestamp.Format("2006-01-02 15:04:05"),
	)

	if alert.Symbol != "" {
		content += fmt.Sprintf("\n交易对: %s", alert.Symbol)
	}
	if alert.Price > 0 {
		content += fmt.Sprintf("\n价格: %.2f", alert.Price)
	}

	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": alert.Title,
			"text":  content,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dingtalk webhook returned status: %d", resp.StatusCode)
	}

	d.logger.Info("dingtalk alert sent", zap.String("title", alert.Title))
	return nil
}

// WeChatNotifier 企业微信通知
type WeChatNotifier struct {
	webhookURL string
	logger     *zap.Logger
}

// NewWeChatNotifier 创建企业微信通知器
func NewWeChatNotifier(webhookURL string, logger *zap.Logger) *WeChatNotifier {
	return &WeChatNotifier{
		webhookURL: webhookURL,
		logger:     logger,
	}
}

// Send 发送企业微信通知
func (w *WeChatNotifier) Send(ctx context.Context, alert Alert) error {
	content := fmt.Sprintf("%s\n\n%s\n\n时间: %s",
		alert.Title,
		alert.Message,
		alert.Timestamp.Format("2006-01-02 15:04:05"),
	)

	if alert.Symbol != "" {
		content += fmt.Sprintf("\n交易对: %s", alert.Symbol)
	}
	if alert.Price > 0 {
		content += fmt.Sprintf("\n价格: %.2f", alert.Price)
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wechat webhook returned status: %d", resp.StatusCode)
	}

	w.logger.Info("wechat alert sent", zap.String("title", alert.Title))
	return nil
}

// EmailNotifier 邮件通知
type EmailNotifier struct {
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
	logger   *zap.Logger
}

// NewEmailNotifier 创建邮件通知器
func NewEmailNotifier(smtpHost, smtpPort, username, password, from string, logger *zap.Logger) *EmailNotifier {
	return &EmailNotifier{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
		logger:   logger,
	}
}

// Send 发送邮件通知
func (e *EmailNotifier) Send(ctx context.Context, alert Alert) error {
	to := alert.Metadata["email"]
	if to == "" {
		return fmt.Errorf("recipient email not provided")
	}

	subject := fmt.Sprintf("[%s] %s", alert.Level, alert.Title)
	body := fmt.Sprintf("Alert: %s\n\n%s\n\nTime: %s",
		alert.Title,
		alert.Message,
		alert.Timestamp.Format("2006-01-02 15:04:05"),
	)

	if alert.Symbol != "" {
		body += fmt.Sprintf("\nSymbol: %s", alert.Symbol)
	}
	if alert.Price > 0 {
		body += fmt.Sprintf("\nPrice: %.2f", alert.Price)
	}

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	err := smtp.SendMail(addr, auth, e.from, []string{to}, msg)
	if err != nil {
		return err
	}

	e.logger.Info("email alert sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}

// MultiNotifier 多通道通知器
type MultiNotifier struct {
	notifiers []Notifier
	logger    *zap.Logger
}

// NewMultiNotifier 创建多通道通知器
func NewMultiNotifier(logger *zap.Logger) *MultiNotifier {
	return &MultiNotifier{
		notifiers: make([]Notifier, 0),
		logger:    logger,
	}
}

// AddNotifier 添加通知器
func (m *MultiNotifier) AddNotifier(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Send 发送通知到所有通道
func (m *MultiNotifier) Send(ctx context.Context, alert Alert) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(ctx, alert); err != nil {
			m.logger.Error("failed to send alert", zap.Error(err))
			lastErr = err
		}
	}
	return lastErr
}

// AlertManager 告警管理器
type AlertManager struct {
	notifier Notifier
	logger   *zap.Logger
}

// NewAlertManager 创建告警管理器
func NewAlertManager(notifier Notifier, logger *zap.Logger) *AlertManager {
	return &AlertManager{
		notifier: notifier,
		logger:   logger,
	}
}

// SendPriceAlert 发送价格告警
func (am *AlertManager) SendPriceAlert(ctx context.Context, symbol string, price float64, condition string) error {
	alert := Alert{
		Title:     fmt.Sprintf("价格告警: %s", symbol),
		Message:   fmt.Sprintf("%s 价格 %s 触发告警，当前价格: %.2f", symbol, condition, price),
		Level:     "warning",
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now(),
	}
	return am.notifier.Send(ctx, alert)
}

// SendTradeAlert 发送交易告警
func (am *AlertManager) SendTradeAlert(ctx context.Context, symbol, side string, qty, price float64) error {
	alert := Alert{
		Title:     fmt.Sprintf("交易执行: %s", symbol),
		Message:   fmt.Sprintf("%s %s %.4f @ %.2f", side, symbol, qty, price),
		Level:     "info",
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"side": side,
			"qty":  fmt.Sprintf("%.4f", qty),
		},
	}
	return am.notifier.Send(ctx, alert)
}

// SendSystemAlert 发送系统告警
func (am *AlertManager) SendSystemAlert(ctx context.Context, title, message string, level string) error {
	alert := Alert{
		Title:     title,
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
	}
	return am.notifier.Send(ctx, alert)
}
