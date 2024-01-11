package psnats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"kiramishima/m-backend/internal/core/domain"
)

type NATSPubSub struct {
	client *nats.Conn
}

// NewNATSPubSub creates a new instance of NATSPubSub
func NewNATSPubSub(nats_addr string) (*NATSPubSub, error) {
	nc, err := nats.Connect(nats_addr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS Server: %v", err)
	}

	return &NATSPubSub{client: nc}, nil
}

func (nc *NATSPubSub) PublishEvent(topic string, event any) error {
	data, _ := json.Marshal(event)

	err := nc.client.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

// Module
var Module = fx.Module("pubsub",
	fx.Provide(func(cfg *domain.Configuration) *NATSPubSub {
		cache, _ := NewNATSPubSub(cfg.NATS_Addr)
		return cache
	}),
)
