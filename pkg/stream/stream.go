package stream

import (
	"context"
	"errors"
	"fmt"

	"braces.dev/errtrace"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Stream struct {
	js   jetstream.JetStream
	name string
}

func NewStream(nc *nats.Conn, name string) (*Stream, error) {
	js, _ := jetstream.New(nc)
	stream := &Stream{
		js:   js,
		name: name,
	}

	if _, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     name,
		Subjects: []string{fmt.Sprintf("%s.>", name)},
	}); err != nil {
		if errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
			return stream, nil
		}
		return nil, errtrace.Wrap(err)
	}

	return stream, nil
}

func (s *Stream) Publish(ctx context.Context, subjectName string, data []byte) error {
	_, err := s.js.Publish(ctx, fmt.Sprintf("%s.%s", s.name, subjectName), data)
	return errtrace.Wrap(err)
}
