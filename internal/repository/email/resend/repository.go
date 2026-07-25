package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/rizkiar00/homework/pkg/config"
)

const sendEmailURL = "https://api.resend.com/emails"

type Repository struct {
	cfg    config.EmailConfig
	client *http.Client
}

func New(cfg config.Config) *Repository {
	return &Repository{
		cfg: cfg.Email,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

func (r *Repository) SendVerificationCode(ctx context.Context, to string, name string, code string) error {
	return r.sendCode(ctx, to, name, "Verify your Rizki Starter account", "verification", code)
}

func (r *Repository) SendPasswordResetCode(ctx context.Context, to string, name string, code string) error {
	return r.sendCode(ctx, to, name, "Reset your Rizki Starter password", "reset password", code)
}

func (r *Repository) sendCode(ctx context.Context, to string, name string, subject string, codeLabel string, code string) error {
	if !r.cfg.IsConfigured() {
		return errors.New("email config is not complete")
	}

	body := sendEmailRequest{
		From:    r.cfg.From,
		To:      []string{to},
		Subject: subject,
		Text:    fmt.Sprintf("Your %s code is %s. This code expires in 15 minutes.", codeLabel, code),
		HTML: fmt.Sprintf(
			"<p>Hello %s,</p><p>Your %s code is:</p><h2>%s</h2><p>This code expires in 15 minutes.</p>",
			html.EscapeString(name),
			html.EscapeString(codeLabel),
			html.EscapeString(code),
		),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendEmailURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend send email failed with status %d", resp.StatusCode)
	}

	return nil
}

func (r *Repository) Provider() string {
	return strings.TrimSpace(r.cfg.ProviderValue())
}
