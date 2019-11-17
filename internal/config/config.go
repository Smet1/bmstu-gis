package config

type DB struct {
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
	Host     string `yaml:"host"`
}

type Config struct {
	ServeAddr string `yaml:"serve_addr"`
	DB        DB     `yaml:"db"`
}
