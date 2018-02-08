package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *SearchApp) initCapKeys() {

	// All top-level capabilities keys
	// should be constructed here.
	// For more information, see http://www.erights.org/elib/capability/ode/ode.pdf.
	app.capKeyMainStore = sdk.NewKVStoreKey("main")
}
