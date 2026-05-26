package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultDevPaymentCode = "LOCAL-PAID"
	defaultTokenSymbol    = "MERGE"
	defaultGitHubOwner    = "mergeos-bounties"
)

type Config struct {
	TokenSymbol       string
	StatePath         string
	PlatformFeeBps    int64
	DevPaymentEnabled bool
	DevPaymentCode    string

	PayPalEnvironment  string
	PayPalClientID     string
	PayPalClientSecret string

	CryptoRPCURL           string
	CryptoReceiver         string
	CryptoAsset            string
	CryptoTokenContract    string
	CryptoTokenDecimals    int
	CryptoWeiPerUSDCent    string
	CryptoMinConfirmations int64

	GitHubToken     string
	GitHubOwner     string
	GitHubOwnerType string

	BountyRoot string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func LoadConfig() Config {
	statePath := getenv("MERGEOS_STATE_PATH", filepath.Join("data", "mergeos-state.json"))
	bountyRoot := getenv("BOUNTY_ROOT", filepath.Join("..", "bounties"))

	return Config{
		TokenSymbol:       getenv("TOKEN_SYMBOL", defaultTokenSymbol),
		StatePath:         statePath,
		PlatformFeeBps:    getenvInt64("PLATFORM_FEE_BPS", 1000),
		DevPaymentEnabled: getenvBool("DEV_PAYMENT_ENABLED", true),
		DevPaymentCode:    getenv("DEV_PAYMENT_CODE", defaultDevPaymentCode),

		PayPalEnvironment:  strings.ToLower(getenv("PAYPAL_ENV", "sandbox")),
		PayPalClientID:     os.Getenv("PAYPAL_CLIENT_ID"),
		PayPalClientSecret: os.Getenv("PAYPAL_CLIENT_SECRET"),

		CryptoRPCURL:           os.Getenv("CRYPTO_RPC_URL"),
		CryptoReceiver:         strings.ToLower(os.Getenv("CRYPTO_RECEIVER")),
		CryptoAsset:            strings.ToLower(getenv("CRYPTO_ASSET", "native")),
		CryptoTokenContract:    strings.ToLower(os.Getenv("CRYPTO_TOKEN_CONTRACT")),
		CryptoTokenDecimals:    int(getenvInt64("CRYPTO_TOKEN_DECIMALS", 6)),
		CryptoWeiPerUSDCent:    os.Getenv("CRYPTO_WEI_PER_USD_CENT"),
		CryptoMinConfirmations: getenvInt64("CRYPTO_MIN_CONFIRMATIONS", 1),

		GitHubToken:     os.Getenv("GITHUB_TOKEN"),
		GitHubOwner:     getenv("GITHUB_OWNER", defaultGitHubOwner),
		GitHubOwnerType: strings.ToLower(getenv("GITHUB_OWNER_TYPE", "org")),

		BountyRoot: bountyRoot,

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     getenv("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     getenv("SMTP_FROM", "noreply@mergeos.local"),
	}
}

func (c Config) PayPalReady() bool {
	return c.PayPalClientID != "" && c.PayPalClientSecret != ""
}

func (c Config) CryptoReady() bool {
	if c.CryptoRPCURL == "" || c.CryptoReceiver == "" {
		return false
	}
	if c.CryptoAsset == "erc20" {
		return c.CryptoTokenContract != ""
	}
	return c.CryptoWeiPerUSDCent != ""
}

func (c Config) GitHubReady() bool {
	return c.GitHubToken != "" && c.GitHubOwner != ""
}

func (c Config) SMTPReady() bool {
	return c.SMTPHost != "" && c.SMTPUsername != "" && c.SMTPPassword != "" && c.SMTPFrom != ""
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getenvInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
