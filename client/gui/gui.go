package gui

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/VincNT21/kallaxy/client/context"
)

func StartGui(appCtxt *context.AppContext) {
	a := app.New()
	appCtxt.LoadsAppstate()
	if appCtxt.APIClient.LastUser.Username == "" {
		getLoginWindow(a, appCtxt)
	}

	a.Run()
	log.Print("--DEBUG-- Client GUI exited")
	// appCtxt.DumpAppstate()
}
