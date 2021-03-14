package configs

type Config struct {
	ServAddr   string `toml:"server_addr"`
	LogLevel   string `toml:"log_level"`
	PostgresBD string `toml:"postgres_bd"`
}

func NewConfig() *Config {
	return &Config{
		ServAddr:   ":5000",
		LogLevel:   "DEBUG",
		PostgresBD: "",
	}
}
