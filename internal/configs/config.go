package configs

type Configs struct {
	Port string     `json:"port"`
	DB   *DBCfg     `json:"db"`
	Aws  *AwsConfig `json:"aws"`
}

type DBCfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"db_name"`
	SslMode  string `json:"ssl_mode"`
}

type AwsConfig struct {
	Host      string `json:"host"`
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewConfig() *Configs {
	return &Configs{
		DB:  &DBCfg{},
		Aws: &AwsConfig{},
	}
}
