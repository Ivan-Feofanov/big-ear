package stream

import (
	"context"
	"fmt"

	"braces.dev/errtrace"
	"github.com/Ivan-Feofanov/big-ear/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Stream struct {
	js   jetstream.JetStream
	name string
}

func (s *Stream) Name() string {
	return s.name
}

func NewStream(nc *nats.Conn, natsConfig config.StreamConfig) (*Stream, error) {
	js, _ := jetstream.New(nc)
	stream := &Stream{
		js:   js,
		name: natsConfig.StreamName,
	}

	if _, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:              natsConfig.StreamName,
		Subjects:          []string{fmt.Sprintf("%s.>", natsConfig.StreamName)},
		MaxMsgsPerSubject: natsConfig.MaxMsgs,
		MaxBytes:          natsConfig.MaxBytes,
	}); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return stream, nil
}

func (s *Stream) Publish(ctx context.Context, subjectName string, data []byte) error {
	_, err := s.js.Publish(ctx, fmt.Sprintf("%s.%s", s.name, subjectName), data)
	return errtrace.Wrap(err)
}
