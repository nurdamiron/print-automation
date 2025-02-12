// internal/config/config.go
package config

type Config struct {
    DatabaseURL string
    ServerAddr  string
    JWTSecret   string
    S3Bucket    string
    // Payment gateway configs
    PaymentGatewayURL  string
    PaymentGatewayKey  string
}

func Load() (*Config, error) {
    // Load configuration from environment variables
    // or configuration file
    return &Config{ 
        DatabaseURL: "mysql://root:print0101@print.czwiyugwum02.eu-north-1.rds.amazonaws.com/root", 
        ServerAddr:  ":3306",
    }, nil
}