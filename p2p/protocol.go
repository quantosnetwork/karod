package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"karod/protocol/libp2p_grpc_protocol"
)

const GossipProtocolPrefix = "/quantos/karod/gossip/1.0.0"
const GossipProtocolTopic = "quantos"
const GRPCProtocolPrefix = libp2p_grpc_protocol.Protocol
const ConsensusProtocolPrefiix = "/quantos/consensus/1.0.0"
const QuantosLiveProtocolPrefix = "/quantos/live/1.0.0"
const QuantosTestProtocolPrefix = "/quantostestnet/test/1.0.0"

type TopicSub struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}
