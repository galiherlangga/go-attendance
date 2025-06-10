package config

type AppConfig struct {
	Port string
	Host string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port: GetEnv("APP_PORT", "8080"),
		Host: GetEnv("APP_HOST", "localhost"),
	}
}