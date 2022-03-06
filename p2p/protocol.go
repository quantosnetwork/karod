package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const GossipProtocolPrefix = "/quantos/karod/gossip/1.0.0"
const GossipProtocolTopic = "quantos"

const ConsensusProtocolPrefiix = "/quantos/consensus/1.0.0"
const QuantosLiveProtocolPrefix = "/quantos/live/1.0.0"
const QuantosTestProtocolPrefix = "/quantostestnet/test/1.0.0"

type TopicSub struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}
