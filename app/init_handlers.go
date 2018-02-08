package app

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/sketchy"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *SearchApp) initHandlers() {
	app.initDefaultAnteHandler()
	app.initRouterHandlers()
}

func (app *SearchApp) initDefaultAnteHandler() {

	// Deducts fee from payer.
	// Verifies signatures and nonces.
	// Sets Signers to ctx.
	app.BaseApp.SetDefaultAnteHandler(
		auth.NewAnteHandler(app.accts))
}

func (app *SearchApp) initRouterHandlers() {

	// All handlers must be added here.
	// The order matters.
	app.router.AddRoute("bank", bank.NewHandler(app.accts))
	app.router.AddRoute("sketchy", sketchy.NewHandler())
}
