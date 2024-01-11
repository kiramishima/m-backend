package domain

type Configuration struct {
	HTTPServer
	Database
	Cache
	ContextTimeout int    `envconfig:"CONTEXT_TIMEOUT" default:"2"`
	NATS_Addr      string `envconfig:"NATS_ADDR" default:"nats://localhost:4222"`
}
