package gui

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/VincNT21/kallaxy/client/context"
)

func StartGui(appCtxt *context.AppContext) {
	// Create a New App
	appGui := app.New()
	log.Print("--INFO-- Client GUI started")

	// Initialize Page Manager
	pageManager := GuiPageManager{
		appCtxt: appCtxt,
		appGui:  appGui,
	}

	// Assign the PageManager to the appContext's pageManager field
	appCtxt.PageManager = &pageManager

	appCtxt.LoadsAppstate()
	if appCtxt.APIClient.CurrentUser.Username != "" {
		pageManager.GetBackWindow()
	} else {
		pageManager.GetLoginWindow()
	}

	appGui.Run()
	exitGui(appCtxt)
}

func exitGui(appCtxt *context.AppContext) {
	log.Print("--INFO-- Client GUI exited")
	appCtxt.DumpAppstate()
	appCtxt.APIClient.Cache.DumpCacheFile()
}
