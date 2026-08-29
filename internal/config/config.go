package config 

type Config struct {
    AppID int 
    AppHash string 
    Phone string 
    SessionFilePath string 
}

func NewConfig(appID int, appHash, phone, sessionFilePath string) *Config {
    return &Config{
    AppID:  appID,
    AppHash:  appHash,
    Phone: phone,
    SessionFilePath: sessionFilePath,
    }
}

