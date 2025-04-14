package gui

import (
	"fmt"
	"image/color"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
	err := appCtxt.LoadMetadataFieldsSpecs()
	if err != nil {
		appCtxt.MetadataFieldsSpecs = kallaxyapi.InitMetadataFieldsSpecs()
	}

	err = appCtxt.LoadMetadataFieldsMap()
	if err != nil {
		appCtxt.MetadataFieldsMap = kallaxyapi.InitMetadataFieldsMap()
	}

	// Show first page
	if appCtxt.APIClient.CurrentUser.Username != "" {
		pageManager.ShowWelcomeBackPage()
	} else {
		pageManager.ShowLoginPage()
	}

	// Set window to be centered and visible
	mainWindow.CenterOnScreen()
	mainWindow.Show()

	// Check if a new version of client exists
	isNewVersion, versionTag, versionDescription, dlUrl, err := appCtxt.APIClient.Admin.CheckForClientUpdate(appCtxt.APIClient.ClientVersion)
	if err != nil {
		dialog.ShowError(err, appCtxt.MainWindow)
	} else {
		if isNewVersion {
			line1 := canvas.NewText("A new client version has been found !", color.White)
			line1.Alignment = fyne.TextAlignCenter
			line2 := canvas.NewText(fmt.Sprintf("New version is: %s", versionTag), color.White)
			line2.Alignment = fyne.TextAlignCenter
			description := widget.NewLabel(versionDescription)
			description.Wrapping = fyne.TextWrapWord
			content := container.NewVBox(
				line1,
				line2,
				description,
			)
			dialog.ShowCustomConfirm("New version", "Download it", "Dismiss", content, func(b bool) {
				if b {
					urlObj, err := url.Parse(dlUrl)
					if err != nil {
						log.Printf("--ERROR-- invalid url for downloading new client: %s", err)
					}
					if err = appGui.OpenURL(urlObj); err != nil {
						log.Printf("--ERROR-- with appGui.OpenURL() : %s", err)
					}
				}
			}, appCtxt.MainWindow)
		}
	}

	// Check the server version
	serverVersion, err := appCtxt.APIClient.Admin.CheckForServerVersion()
	if err == nil {
		appCtxt.APIClient.ServerVersion = serverVersion
	}

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
