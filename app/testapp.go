package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
)

type testBasecoinApp struct {
	*SearchApp
	*bam.TestApp
}

func newTestSearchApp() *testBasecoinApp {
	app := NewSearchApp()
	tba := &testBasecoinApp{
		SearchApp: app,
	}
	tba.TestApp = bam.NewTestApp(app.BaseApp)
	return tba
}

