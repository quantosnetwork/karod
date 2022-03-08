package tests

import (
	"context"
	"github.com/hashicorp/go-hclog"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/test"
	mnet "github.com/libp2p/go-libp2p-testing/mocks/network"
	"github.com/quantosnetwork/karod/p2p"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSomething(t *testing.T) {

	assert.True(t, true, "True is true!")

}

type TestSuite struct {
	suite suite.TestingSuite
	KarodTestSuite
	ctx     context.Context
	Peer    mnet.MockPeerScope
	SK      ci.PrivKey
	PK      ci.PubKey
	Server  *p2p.Server
	SConfig *p2p.ServerConfig
}

type KarodTestSuite interface {
	SetupSuite(t *testing.T)
	TearDownSuite()
	InitializeContext()
	InitializeP2P(t *testing.T)
}

func (ts *TestSuite) SetupSuite(t *testing.T) {
	ts.suite.SetT(t)
	T := ts.suite.T()
	ts.InitializeP2P(T)

}

func (ts *TestSuite) InitializeContext() {
	ts.ctx = context.Background()
}

func (ts *TestSuite) InitializeP2P(t *testing.T) {
	sk, pk, err := test.RandTestKeyPair(ci.Secp256k1, 256)

	if err != nil {
		t.Error(err.Error())
	}
	ts.SK = sk
	ts.PK = pk

	sc := p2p.SetDefaultConfig()
	ts.SConfig = sc
	logger := hclog.NewNullLogger()
	s, err := p2p.NewServer(logger, sc)
	if err != nil {
		t.Error(err.Error())
	}
	ts.Server = s
	ts.suite.SetT(t)

}

func (ts *TestSuite) TearDownSuite() {
	ts.Server.Disconnect(ts.Peer.Peer(), "test ended")
}

func NewP2PTestSuite(t *testing.T) KarodTestSuite {
	ts := &TestSuite{}
	ts.SetupSuite(t)
	return ts
}
