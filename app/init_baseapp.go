package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/tendermint/go-wire"
	crypto "github.com/tendermint/go-crypto"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *SearchApp) initBaseApp() {
	bapp := baseapp.NewBaseApp(AppName)
	app.BaseApp = bapp
	app.router = bapp.Router()
	app.initBaseAppTxDecoder()
}

func (app *SearchApp) initBaseAppTxDecoder() {
	cdc := makeTxCodec()
	app.BaseApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = sdk.StdTx{}
		// StdTx.Msg is an interface whose concrete
		// types are registered in app/msgs.go.
		err := cdc.UnmarshalBinary(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxParse("").TraceCause(err, "")
		}
		return tx, nil
	})
}

// Wire requires registration of interfaces & concrete types.  All
// interfaces to be encoded/decoded in a Msg must be registered
// here, along with all the concrete types that implement them.
func makeTxCodec() (cdc *wire.Codec) {
	cdc = wire.NewCodec()

	// Register crypto.[PubKey,PrivKey,Signature] types.
	crypto.RegisterWire(cdc)

	// Register bank.[SendMsg,IssueMsg] types.
	bank.RegisterWire(cdc)

	return
}
