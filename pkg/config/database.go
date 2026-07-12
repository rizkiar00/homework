package config

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	URL      string `env:"DATABASE_URL"`
	Name     string `env:"DB_NAME"`
	Driver   string `env:"DB_DRIVER" envDefault:"postgres"`
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Schema   string `env:"DB_SCHEMA" envDefault:"public"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

func (c DatabaseConfig) IsConfigured() bool {
	if strings.TrimSpace(c.URL) != "" {
		return true
	}

	return c.Host != "" && c.Name != "" && c.Username != ""
}

func (c DatabaseConfig) GenerateConnectionString() string {
	if strings.TrimSpace(c.URL) != "" {
		return c.URL
	}

	if !c.IsConfigured() {
		return ""
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=%s&connect_timeout=5",
		c.Username,
		c.Password,
		c.ResolveHost(),
		c.Port,
		c.Name,
		c.Schema,
		c.SSLModeValue(),
	)
}

func (c DatabaseConfig) SSLModeValue() string {
	if strings.TrimSpace(c.SSLMode) == "" {
		return "disable"
	}

	return c.SSLMode
}

func (c DatabaseConfig) ResolveHost() string {
	if c.Host != "auto" && c.Host != "windows-host" {
		return c.Host
	}

	host, err := resolveWSLWindowsHost()
	if err != nil {
		return "127.0.0.1"
	}

	return host
}

func resolveWSLWindowsHost() (string, error) {
	host, err := resolveWSLDefaultGateway()
	if err == nil {
		return host, nil
	}

	return resolveWSLNameserver()
}

func resolveWSLDefaultGateway() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 || fields[1] != "00000000" {
			continue
		}

		gatewayHex := fields[2]
		if len(gatewayHex) != 8 {
			return "", fmt.Errorf("invalid gateway format")
		}

		bytes, err := hex.DecodeString(gatewayHex)
		if err != nil {
			return "", err
		}

		return strings.Join([]string{
			strconv.Itoa(int(bytes[3])),
			strconv.Itoa(int(bytes[2])),
			strconv.Itoa(int(bytes[1])),
			strconv.Itoa(int(bytes[0])),
		}, "."), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("default gateway not found")
}

func resolveWSLNameserver() (string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "nameserver" {
			return fields[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("nameserver not found")
}
