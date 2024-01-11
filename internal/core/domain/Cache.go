package domain

type Cache struct {
	Addr     string `envconfig:"CACHE_ADDR" required:"true" default:"localhost"`
	Password string `envconfig:"CACHE_PWD" default:""`
}
