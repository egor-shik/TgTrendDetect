package config 

type Config struct {
    AppID int 
    AppHash string 
    Phone string 
    SessionFilePath string 
}

func NewConfig() *Config {
    return &Config{...}
}