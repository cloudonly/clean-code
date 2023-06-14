package config

type Server struct {
	//JWT   JWT   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap   Zap   `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	// gorm
	MySQL MySQL `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	//PgSQL PgSQL `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	// cors
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}
