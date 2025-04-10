package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

type GuiPageManager struct {
	appCtxt    *context.AppContext
	mainWindow fyne.Window
}

// Methods to show different pages

func (pm *GuiPageManager) ShowLoginPage() {
	content := createLoginContent(pm.appCtxt)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - Login")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowWelcomeBackPage() {
	content := createWelcomeBackContent(pm.appCtxt)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - Welcome back")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowHomePage() {
	content := createHomepageContent(pm.appCtxt)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - Homepage")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowCreateUserPage() {
	content := createUserContent(pm.appCtxt)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - New user")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowCreateMediaPage(mediaType string) {
	content := createMediaCreationContent(pm.appCtxt, mediaType)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - New media")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowShelfPage() {
	mediaRecords, err := pm.appCtxt.APIClient.Media.GetMediaWithRecords()
	if err != nil || len(mediaRecords.MediaRecords) == 0 {
		dialog.ShowInformation("Information", "There is no media in your shelf\nGo create some !", pm.mainWindow)
		return
	}
	content := createShelfContent(pm.appCtxt, mediaRecords)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - My Shelf")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowCompartmentMediaPage(mediaType string, mediaList []models.MediumWithRecord) {
	content := createMediaListContent(pm.appCtxt, mediaType, mediaList)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle(fmt.Sprintf("Kallaxy - My Shelf %s Compartment", mediaType))
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}

func (pm *GuiPageManager) ShowParametersPage() {
	content := createParametersContent(pm.appCtxt)
	pm.mainWindow.SetContent(content)
	pm.mainWindow.SetTitle("Kallaxy - User Parameters")
	// Resize if needed
	pm.mainWindow.Resize(fyne.NewSize(1024, 768))
}
