package config

type Config struct {
	Server ServerConfig
	DB     DBConfig
	CORS   CORSConfig
	App    AppConfig
	Bot    BotConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

type AppConfig struct {
	Dir string
}

type BotConfig struct {
	URL   string
	Token string
}

type CORSConfig struct {
	AllowedOrigins []string
}
