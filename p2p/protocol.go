package p2p

import (
	"context"
	"github.com/hashicorp/go-hclog"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"
	"reflect"
)

const bufferSize = 1024

const GossipProtocolPrefix = "/quantos/karod/gossip/1.0.0"
const GossipProtocolTopic = "quantos"

const ConsensusProtocolPrefiix = "/quantos/consensus/1.0.0"
const QuantosLiveProtocolPrefix = "/quantos/live/1.0.0"
const QuantosTestProtocolPrefix = "/quantostestnet/test/1.0.0"

type TopicSub struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

type Topic struct {
	logger  hclog.Logger
	topic   *pubsub.Topic
	typ     reflect.Type
	closeCh chan struct{}
}

func (t *Topic) createObj() proto.Message {
	message, ok := reflect.New(t.typ).Interface().(proto.Message)
	if !ok {
		return nil
	}

	return message
}

func (t *Topic) Publish(obj proto.Message) error {
	data, err := proto.Marshal(obj)
	if err != nil {
		return err
	}

	return t.topic.Publish(context.Background(), data)
}

func (t *Topic) Subscribe(handler func(obj interface{})) error {
	sub, err := t.topic.Subscribe(pubsub.WithBufferSize(bufferSize))
	if err != nil {
		return err
	}

	go t.readLoop(sub, handler)

	return nil
}

func (t *Topic) readLoop(sub *pubsub.Subscription, handler func(obj interface{})) {
	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		<-t.closeCh
		cancelFn()
	}()

	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			t.logger.Error("failed to get topic", "err", err)

			continue
		}

		go func() {
			obj := t.createObj()
			if err := proto.Unmarshal(msg.Data, obj); err != nil {
				t.logger.Error("failed to unmarshal topic", "err", err)

				return
			}

			handler(obj)
		}()
	}
}

func (s *Server) NewTopic(protoID string, obj proto.Message) (*Topic, error) {
	topic, err := s.ps.Join(protoID)
	if err != nil {
		return nil, err
	}

	tt := &Topic{
		logger: s.logger.Named(protoID),
		topic:  topic,
		typ:    reflect.TypeOf(obj).Elem(),
	}

	return tt, nil
}
