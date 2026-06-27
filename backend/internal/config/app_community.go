package config

import (
	"fmt"
	"strings"
)

const (
	DefaultCommunityAuthAccessTokenTTLSeconds  = 900
	DefaultCommunityAuthRefreshTokenTTLSeconds = 604800
	DefaultCommunityAuthCookieNamePrefix       = "community"
	DefaultCommunityAuthCookiePath             = "/"
	DefaultCommunityAuthCookieSameSite         = "lax"
	DefaultCommunityAuthCSRFCookieName         = "community_csrf"
	DefaultCommunityAuthCSRFHeaderName         = "X-Community-CSRF-Token"
	DefaultCommunityAuthClientType             = "community_web"
)

type CommunityConfig struct {
	Auth CommunityAuthConfig `mapstructure:"auth" json:"auth" yaml:"auth" toml:"auth"`
}

type CommunityAuthConfig struct {
	AccessTokenTTLSeconds  int                       `mapstructure:"access_token_ttl_seconds" envname:"COMMUNITY_AUTH_ACCESS_TOKEN_TTL_SECONDS" json:"access_token_ttl_seconds" yaml:"access_token_ttl_seconds" toml:"access_token_ttl_seconds"`
	RefreshTokenTTLSeconds int                       `mapstructure:"refresh_token_ttl_seconds" envname:"COMMUNITY_AUTH_REFRESH_TOKEN_TTL_SECONDS" json:"refresh_token_ttl_seconds" yaml:"refresh_token_ttl_seconds" toml:"refresh_token_ttl_seconds"`
	Cookie                 CommunityAuthCookieConfig `mapstructure:"cookie" json:"cookie" yaml:"cookie" toml:"cookie"`
	CSRF                   CommunityAuthCSRFConfig   `mapstructure:"csrf" json:"csrf" yaml:"csrf" toml:"csrf"`
	DefaultClientType      string                    `mapstructure:"default_client_type" envname:"COMMUNITY_AUTH_DEFAULT_CLIENT_TYPE" json:"default_client_type" yaml:"default_client_type" toml:"default_client_type"`
}

type CommunityAuthCookieConfig struct {
	NamePrefix string `mapstructure:"name_prefix" envname:"COMMUNITY_AUTH_COOKIE_NAME_PREFIX" json:"name_prefix" yaml:"name_prefix" toml:"name_prefix"`
	Domain     string `mapstructure:"domain" envname:"COMMUNITY_AUTH_COOKIE_DOMAIN" json:"domain" yaml:"domain" toml:"domain"`
	Path       string `mapstructure:"path" envname:"COMMUNITY_AUTH_COOKIE_PATH" json:"path" yaml:"path" toml:"path"`
	SameSite   string `mapstructure:"same_site" envname:"COMMUNITY_AUTH_COOKIE_SAME_SITE" json:"same_site" yaml:"same_site" toml:"same_site"`
	Secure     bool   `mapstructure:"secure" envname:"COMMUNITY_AUTH_COOKIE_SECURE" json:"secure" yaml:"secure" toml:"secure"`
}

type CommunityAuthCSRFConfig struct {
	Enabled    *bool  `mapstructure:"enabled" envname:"COMMUNITY_AUTH_CSRF_ENABLED" json:"enabled" yaml:"enabled" toml:"enabled"`
	CookieName string `mapstructure:"cookie_name" envname:"COMMUNITY_AUTH_CSRF_COOKIE_NAME" json:"cookie_name" yaml:"cookie_name" toml:"cookie_name"`
	HeaderName string `mapstructure:"header_name" envname:"COMMUNITY_AUTH_CSRF_HEADER_NAME" json:"header_name" yaml:"header_name" toml:"header_name"`
}

func (c CommunityAuthCSRFConfig) EnabledValue() bool {
	if c.Enabled == nil {
		return true
	}
	return *c.Enabled
}

func (c *CommunityConfig) ValidateName() string {
	return AppCommunityName
}

func (c *CommunityConfig) ValidateRequired() bool {
	return false
}

func (c *CommunityConfig) Validate() error {
	c.ApplyDefaults()
	if c.Auth.AccessTokenTTLSeconds <= 0 || c.Auth.RefreshTokenTTLSeconds <= 0 {
		return fmt.Errorf("auth token ttl values must be positive")
	}
	if strings.TrimSpace(c.Auth.Cookie.NamePrefix) == "" {
		return fmt.Errorf("auth.cookie.name_prefix is required")
	}
	if strings.TrimSpace(c.Auth.Cookie.Path) == "" {
		return fmt.Errorf("auth.cookie.path is required")
	}
	if !validCookieSameSite(c.Auth.Cookie.SameSite) {
		return fmt.Errorf("auth.cookie.same_site must be one of lax, strict, none")
	}
	if c.Auth.CSRF.EnabledValue() {
		if strings.TrimSpace(c.Auth.CSRF.CookieName) == "" || strings.TrimSpace(c.Auth.CSRF.HeaderName) == "" {
			return fmt.Errorf("auth.csrf cookie_name and header_name are required when csrf is enabled")
		}
	}
	if strings.TrimSpace(c.Auth.DefaultClientType) == "" {
		return fmt.Errorf("auth.default_client_type is required")
	}
	return nil
}

func (c *CommunityConfig) ApplyDefaults() {
	if c.Auth.AccessTokenTTLSeconds == 0 {
		c.Auth.AccessTokenTTLSeconds = DefaultCommunityAuthAccessTokenTTLSeconds
	}
	if c.Auth.RefreshTokenTTLSeconds == 0 {
		c.Auth.RefreshTokenTTLSeconds = DefaultCommunityAuthRefreshTokenTTLSeconds
	}
	if strings.TrimSpace(c.Auth.Cookie.NamePrefix) == "" {
		c.Auth.Cookie.NamePrefix = DefaultCommunityAuthCookieNamePrefix
	}
	if strings.TrimSpace(c.Auth.Cookie.Path) == "" {
		c.Auth.Cookie.Path = DefaultCommunityAuthCookiePath
	}
	if strings.TrimSpace(c.Auth.Cookie.SameSite) == "" {
		c.Auth.Cookie.SameSite = DefaultCommunityAuthCookieSameSite
	}
	c.Auth.Cookie.SameSite = strings.ToLower(strings.TrimSpace(c.Auth.Cookie.SameSite))
	if strings.TrimSpace(c.Auth.CSRF.CookieName) == "" {
		c.Auth.CSRF.CookieName = DefaultCommunityAuthCSRFCookieName
	}
	if strings.TrimSpace(c.Auth.CSRF.HeaderName) == "" {
		c.Auth.CSRF.HeaderName = DefaultCommunityAuthCSRFHeaderName
	}
	if strings.TrimSpace(c.Auth.DefaultClientType) == "" {
		c.Auth.DefaultClientType = DefaultCommunityAuthClientType
	}
}
