package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kidinamoto01/sdk-test/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func TestSendMsg(t *testing.T) {
	tba := newTestSearchApp()
	tba.RunBeginBlock()

	// Construct a SendMsg.
	var msg = types.ProposeMsg{
		Sender:  crypto.GenPrivKeyEd25519().PubKey().Address(),
		Index: 1,
		Candidate:10,
		Name:"first vote",
	}

	// Run a Check on SendMsg.
	res := tba.RunCheckMsg(msg)
	assert.Equal(t, sdk.CodeOK, res.Code, res.Log)

	// Run a Deliver on SendMsg.
	res = tba.RunDeliverMsg(msg)
	assert.Equal(t, sdk.CodeUnrecognizedAddress, res.Code, res.Log)
}
