package app

import (
	"github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/kidinamoto01/sdk-test/types"
)

const AppName = "Ballot"

// ClearchainApp is basic application
type SearchApp struct {
	*baseapp.BaseApp
	accts sdk.AccountMapper
}

func NewSearchApp() *SearchApp {
	// var app = &SearchApp{}

	// make multistore with various keys
	mainKey := sdk.NewKVStoreKey("ballot")
	// ibcKey = sdk.NewKVStoreKey("ibc")

	bApp := baseapp.NewBaseApp(AppName)
	mountMultiStore(bApp, mainKey)
	err := bApp.LoadLatestVersion(mainKey)
	if err != nil {
		panic(err)
	}

	// register routes on new application
	accts := types.AccountMapper(mainKey)
	types.RegisterRoutes(bApp.Router(), accts)

	// set up ante and tx parsing
	setAnteHandler(bApp, accts)
	initBaseAppTxDecoder(bApp)

	return &SearchApp{
		BaseApp: bApp,
		accts:   accts,
	}
}

// RunForever starts the abci server
func (app *SearchApp) RunForever() {
	srv, err := server.NewServer("127.0.0.1:46658", "socket", app)
	if err != nil {
		panic(err)
	}
	srv.Start()
	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}

func (app *SearchApp) StoreAccount(acct sdk.Account) {
	// delivertx with fake tx bytes (we don't care for SetAccount)
	var ctx = app.NewContext(false, []byte{1, 2, 3, 4})
	app.accts.SetAccount(ctx, acct)
}

func mountMultiStore(bApp *baseapp.BaseApp,
	keys ...*sdk.KVStoreKey) {

	// create substore for every key
	for _, key := range keys {
		bApp.MountStore(key, sdk.StoreTypeIAVL)
	}
}

func setAnteHandler(bApp *baseapp.BaseApp, accts sdk.AccountMapper) {
	// this checks auth, but may take fee is future, check for compatibility
	bApp.SetDefaultAnteHandler(
		auth.NewAnteHandler(accts))
}

func initBaseAppTxDecoder(bApp *baseapp.BaseApp) {
	cdc := types.MakeTxCodec()
	bApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
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
