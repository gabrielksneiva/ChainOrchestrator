package eventbus

import "context"

// Publisher interface para publicação de eventos
type Publisher interface {
	Publish(ctx context.Context, message string, chainType string) error
}
