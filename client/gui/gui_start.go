package gui

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/VincNT21/kallaxy/client/context"
	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
)

func StartGui(appCtxt *context.AppContext) {
	// Create a New App
	appGui := app.New()
	log.Print("--INFO-- Client GUI started")

	// Create the main window tha will persist throughout the app lifecycle
	mainWindow := appGui.NewWindow("Kallaxy")

	// Store it in the AppContext
	appCtxt.MainWindow = mainWindow

	// Initialize the Page Manager
	pageManager := &GuiPageManager{
		appCtxt:    appCtxt,
		mainWindow: mainWindow,
	}

	// Assign the PageManager to the appContext
	appCtxt.PageManager = pageManager

	// Load application state and metadata fields
	appCtxt.LoadAppstate()
	// appCtxt.LoadMetadataFieldsSpecs()
	appCtxt.MetadataFieldsSpecs = kallaxyapi.InitMetadataFieldsSpecs()
	// appCtxt.LoadMetadataFieldsMap()
	appCtxt.MetadataFieldsMap = kallaxyapi.InitMetadataFieldsMap()

	if appCtxt.APIClient.CurrentUser.Username != "" {
		pageManager.ShowWelcomeBackPage()
	} else {
		pageManager.ShowLoginPage()
	}

	// Set window to be centered and visible
	mainWindow.CenterOnScreen()
	mainWindow.Show()

	// Start the application main loop
	appGui.Run()

	// Handle cleanup/saving when app exits
	exitGui(appCtxt)
}

func exitGui(appCtxt *context.AppContext) {
	log.Print("--INFO-- Client GUI exited")
	appCtxt.SaveAppstate()
	appCtxt.SaveMetadataFieldsSpecs()
	appCtxt.SaveMedataFieldsMap()
	appCtxt.APIClient.Cache.DumpCacheFile()
}
