package configs

type Configs struct {
	Port string `json:"port"`
	DB   *DBCfg `json:"db"`
}

type DBCfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"db_name"`
	SslMode  string `json:"ssl_mode"`
}

func NewConfig() *Configs {
	return &Configs{
		DB: &DBCfg{},
	}
}
