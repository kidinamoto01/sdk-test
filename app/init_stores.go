package app

import (

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/kidinamoto01/sdk-test/types"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *SearchApp) initStores() {
	app.mountStores()
	app.initAccountMapper()
}

// Initialize root stores.
func (app *SearchApp) mountStores() {

	// Create MultiStore mounts.
	app.BaseApp.MountStore(app.capKeyMainStore, sdk.StoreTypeMulti)
}

// Initialize the AccountMapper.
func (app *SearchApp) initAccountMapper() {

	var accountMapper = auth.NewAccountMapper(
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)

	// Register all interfaces and concrete types that
	// implement those interfaces, here.
	cdc := accountMapper.WireCodec()
	auth.RegisterWireBaseAccount(cdc)

	// Make accountMapper's WireCodec() inaccessible.
	app.accts = accountMapper.Seal()
}