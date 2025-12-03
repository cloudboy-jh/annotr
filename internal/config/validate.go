package config

import "strings"

func ValidateAPIKey(provider, key string) bool {
	if key == "" {
		return false
	}

	switch provider {
	case "anthropic":
		return strings.HasPrefix(key, "sk-ant-")
	case "openai":
		return strings.HasPrefix(key, "sk-")
	case "groq":
		return strings.HasPrefix(key, "gsk_")
	case "ollama":
		return true
	default:
		return false
	}
}

func GetProviderFromKey(key string) string {
	switch {
	case strings.HasPrefix(key, "sk-ant-"):
		return "anthropic"
	case strings.HasPrefix(key, "gsk_"):
		return "groq"
	case strings.HasPrefix(key, "sk-"):
		return "openai"
	default:
		return ""
	}
}
