// internal/config/config.go
package config

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
}

type ServerConfig struct {
    Host string
    Port string
}

type DatabaseConfig struct {
    DSN string
}

func LoadConfig() (*Config, error) {
    // В будущем здесь можно добавить загрузку из env файла
    return &Config{
        Server: ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
        Database: DatabaseConfig{
            DSN: "root:print0101@tcp(print.czwiyugwum02.eu-north-1.rds.amazonaws.com:3306)/root",
        },
    }, nil
}