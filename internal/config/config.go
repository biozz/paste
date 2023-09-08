package config

type Server struct {
	Env     string `env:"ENV,default=prod"`
	DSN     string `env:"DB_DSN,default=file:sqlite.db"`
	BaseURL string `env:"BASE_URL,default=http://localhost:3000"`
	Bind    string `env:"BIND,default=:3000"`
}
