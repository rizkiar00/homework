package config

import "strings"

type EmailConfig struct {
	Provider     string `env:"EMAIL_PROVIDER" envDefault:"resend"`
	From         string `env:"EMAIL_FROM"`
	ResendAPIKey string `env:"RESEND_API_KEY"`
}

func (c EmailConfig) IsConfigured() bool {
	return strings.TrimSpace(c.From) != "" && strings.TrimSpace(c.ResendAPIKey) != ""
}

func (c EmailConfig) ProviderValue() string {
	if strings.TrimSpace(c.Provider) == "" {
		return "resend"
	}

	return c.Provider
}
