package gui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

// Struct to define tree node data
type TreeNode struct {
	ID       string // Unique identifier for this node
	ParentID string // Parent node ID
	Title    string // Title to display on the tree, for branch nodes
	Value    string // Value to display, for leaf nodes
	NodeType string // To adapt node's creation and update
}

/*
Node Types to handle :
	- main_title
		branch node with title of a main category
		displayed as a bigger text

	- sub_title
		branch node with title of a sub category
		displayed as a intermediate text

	- medium_title
		branch node with the medium's title
		buttons for update/delete medium or record and for expand/collapse

	- metadata_main
		branch node for metadata
		buttons for expand/collapse

	- single_line
		leaf value node, with single line of text
		displayed as a single line text

	- single_line_with_title
		leaf value node, with both a title and a value
		displayed as a single line, Title Value

	- multi_line
		leaf value node, with multiple lines of text
		displayed as a multi line text, scrollable, with buttons for expand/collapse

	- multi_line_with_title
		leaf value node, with multiple lines of text and a title
		displayed as a multi line text, scrollable, with buttons for expand/collapse

*/

func createAndPopulateTree(appCtxt *context.AppContext, mediaType string, mediaList []models.MediumWithRecord) *widget.Tree {
	// Prepare a map for Tree widget data
	treeData := make(map[string][]string) // Parent -> Children IDs
	nodes := make(map[string]TreeNode)    // NodeID -> TreeNode

	// Create top-level nodes (by status)
	treeData[""] = []string{"Finished", "In Progress", "Unstarted"}
	nodes["Finished"] = TreeNode{ID: "Finished", Title: "Finished", NodeType: "main_title"}
	nodes["In Progress"] = TreeNode{ID: "In Progress", Title: "In Progress", NodeType: "main_title"}
	nodes["Unstarted"] = TreeNode{ID: "Unstarted", Title: "Unstarted", NodeType: "main_title"}

	// Populate media into second and third level
	for _, medium := range mediaList {
		// Medium title branch node (2nd level)

		// Determine parent based on IsFinished and StartDate
		parent := "Unstarted"
		if medium.StartDate != "" {
			parent = "In Progress"
		}
		if medium.IsFinished {
			parent = "Finished"
		}
		// Create a unique ID for this media node
		mediaNodeID := fmt.Sprintf("media-%s", medium.ID)
		treeData[parent] = append(treeData[parent], mediaNodeID)

		// Add the media node with its title
		nodes[mediaNodeID] = TreeNode{
			ID:       mediaNodeID,
			ParentID: parent,
			Title:    medium.Title,
			Value:    medium.MediaID, // Value field will hold MediaID for edit/delete functions
			NodeType: "medium_title",
		}

		// Add third-level nodes for details about the media's PubDate, Creator Personal record and Metadata
		detailsParent := mediaNodeID

		// Creator Leaf node (3rd level)
		creatorNodeID := fmt.Sprintf("%s-creator", mediaNodeID)
		treeData[detailsParent] = append(treeData[detailsParent], creatorNodeID)
		nodes[creatorNodeID] = TreeNode{
			ID:       creatorNodeID,
			ParentID: detailsParent,
			Value:    fmt.Sprintf("Creators: %v", medium.Creator),
			NodeType: "single_line",
		}

		// Pubdate Leaf node (3rd level)
		pubDateNodeID := fmt.Sprintf("%s-pubdate", mediaNodeID)
		treeData[detailsParent] = append(treeData[detailsParent], pubDateNodeID)
		nodes[pubDateNodeID] = TreeNode{
			ID:       pubDateNodeID,
			ParentID: detailsParent,
			Value:    fmt.Sprintf("Publication date: %v", medium.PubDate),
			NodeType: "single_line",
		}

		// Personal record Branch node (3rd level)
		persRecordNodeID := fmt.Sprintf("%s-personal_record", mediaNodeID)
		treeData[detailsParent] = append(treeData[detailsParent], persRecordNodeID)
		nodes[persRecordNodeID] = TreeNode{
			ID:       persRecordNodeID,
			ParentID: detailsParent,
			Title:    "Personal record",
			NodeType: "sub_title",
		}
		// Values Leaf nodes for personal record (4th level)
		// Made one by one to ensure correct order when displayed
		startDateLeafID := fmt.Sprintf("%s-%s", persRecordNodeID, "start_date")
		treeData[persRecordNodeID] = append(treeData[persRecordNodeID], startDateLeafID)
		nodes[startDateLeafID] = TreeNode{
			ID:       startDateLeafID,
			ParentID: persRecordNodeID,
			Title:    "Began: ",
			Value:    fmt.Sprintf("%v", medium.StartDate),
			NodeType: "single_line_with_title",
		}
		endDateLeafID := fmt.Sprintf("%s-%s", persRecordNodeID, "end_date")
		treeData[persRecordNodeID] = append(treeData[persRecordNodeID], endDateLeafID)
		nodes[endDateLeafID] = TreeNode{
			ID:       endDateLeafID,
			ParentID: persRecordNodeID,
			Title:    "Finished on: ",
			Value:    fmt.Sprintf("%v", medium.EndDate),
			NodeType: "single_line_with_title",
		}
		durationLeafID := fmt.Sprintf("%s-%s", persRecordNodeID, "duration")
		treeData[persRecordNodeID] = append(treeData[persRecordNodeID], durationLeafID)
		nodes[durationLeafID] = TreeNode{
			ID:       durationLeafID,
			ParentID: persRecordNodeID,
			Title:    "Duration: ",
			Value:    fmt.Sprintf("%v", medium.Duration),
			NodeType: "single_line_with_title",
		}
		commentsLeafID := fmt.Sprintf("%s-%s", persRecordNodeID, "comments")
		treeData[persRecordNodeID] = append(treeData[persRecordNodeID], commentsLeafID)
		nodes[commentsLeafID] = TreeNode{
			ID:       commentsLeafID,
			ParentID: persRecordNodeID,
			Title:    "Comments: ",
			Value:    medium.Comments,
			NodeType: "multi_line_with_title",
		}

		// Metadata Branch node (3rd level)
		metadataNodeID := fmt.Sprintf("%s-metadata", mediaNodeID)
		treeData[detailsParent] = append(treeData[detailsParent], metadataNodeID)

		nodes[metadataNodeID] = TreeNode{
			ID:       metadataNodeID,
			ParentID: detailsParent,
			Title:    "Metadata",
			NodeType: "metadata_main",
		}

		// Metadata fields and values nodes (4th and 5th level)
		for metaKey, metaValue := range medium.Metadata {
			// Metadata Branch fields nodes (4th level)
			metaFieldID := fmt.Sprintf("%s-%s", metadataNodeID, metaKey)
			treeData[metadataNodeID] = append(treeData[metadataNodeID], metaFieldID)
			nodes[metaFieldID] = TreeNode{
				ID:       metaFieldID,
				ParentID: metadataNodeID,
				Title:    metaKey,
				NodeType: "sub_title",
			}
			// Metadata Leaf values nodes (5th level)
			metaValuesID := fmt.Sprintf("%s-value", metaFieldID)
			treeData[metaFieldID] = append(treeData[metaFieldID], metaValuesID)

			metadataNodeType := "single_line"
			if metaKey == "description" || metaKey == "overview" {
				// Change NodeType to multi line if needed
				metadataNodeType = "multi_line"
			}

			nodes[metaValuesID] = TreeNode{
				ID:       metaValuesID,
				ParentID: metaFieldID,
				Value:    formatMetadataValueForEntry(appCtxt, metaKey, metaValue),
				NodeType: metadataNodeType,
			}

		}
	}

	var tree *widget.Tree

	// Construct the tree using the populated data
	tree = widget.NewTree(
		// Callback function to list children of a node
		func(uid string) []string {
			return treeData[uid]
		},

		// Callback function to determine if a node is a branch (= has children)
		func(uid string) bool {
			_, hasChildren := treeData[uid]
			return hasChildren
		},

		// Callback function to create tree node widgets from IDs
		func(branch bool) fyne.CanvasObject {
			// Create a basic container for now, customization and objects addition will happen on update callback function
			if branch {
				branchContainer := container.NewHBox()
				return branchContainer
			} else {
				leafContainer := container.NewHBox()
				return leafContainer
			}

		},

		// Callback function to update a node's widget with its data
		func(uid string, branch bool, obj fyne.CanvasObject) {
			node := nodes[uid]
			// Update depends on NodeType

			branchContainer := obj.(*fyne.Container)
			// Clear the container first
			branchContainer.Objects = nil
			branchContainer.Refresh()

			// Branch nodes will use a canvas.Text
			branchTextObject := canvas.NewText(node.Title, color.White)
			branchTextObject.Resize(fyne.NewSize(branchContainer.MinSize().Width, 100)) // Seems useless

			switch node.NodeType {
			case "main_title":
				// Main titles will have bold bigger text
				branchTextObject.TextSize = 16
				branchTextObject.TextStyle.Bold = true
				branchContainer.Add(branchTextObject)
			case "sub_title":
				// Sub title will have standard text
				branchTextObject.TextSize = 14
				branchContainer.Add(branchTextObject)
			case "medium_title":
				// Medium title will have standard text, edit/delete button, expand/collapse button
				// color depends of category (finished, in progress, unstarted)
				branchTextObject.TextSize = 14
				switch node.ParentID {
				case "Finished":
					greenColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
					branchTextObject.Color = greenColor
				case "In Progress":
					orangeColor := color.RGBA{R: 255, G: 120, B: 0, A: 255}
					branchTextObject.Color = orangeColor
				case "Unstarted":
					redColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
					branchTextObject.Color = redColor
				}
				branchContainer.Add(branchTextObject)
				// Edit/delete buttons
				mediumEditButton := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
					buttonFuncMediumEdit(appCtxt, node, mediaType, mediaList)
				})
				mediumDeleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					buttonFuncMediumDelete(appCtxt, node)
				})
				expandButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropDown), func() {
					buttonFuncExpandBranches(tree, treeData, node.ID)
				})
				collapseButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropUp), func() {
					buttonFuncCollapseBranches(tree, treeData, node.ID)
				})

				branchContainer.Add(mediumEditButton)
				branchContainer.Add(mediumDeleteButton)
				branchContainer.Add(expandButton)
				branchContainer.Add(collapseButton)

			case "metadata_main":
				// Metadata title will have standard text
				// and a "expand all" button
				branchTextObject.TextSize = 14
				expandButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropDown), func() {
					buttonFuncExpandBranches(tree, treeData, node.ID)
				})
				collapseButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropUp), func() {
					buttonFuncCollapseBranches(tree, treeData, node.ID)
				})

				branchContainer.Add(branchTextObject)
				branchContainer.Add(collapseButton)
				branchContainer.Add(expandButton)

			case "single_line":
				// Use a single canvas.Text
				leafTextObject := canvas.NewText(node.Value, color.White)
				branchContainer.Add(leafTextObject)
			case "single_line_with_title":
				// Use two canvas.Text : one for title, one for value
				leafTitleObject := canvas.NewText(node.Title, color.White)
				leafValueObject := canvas.NewText(node.Value, color.White)
				branchContainer.Add(leafTitleObject)
				branchContainer.Add(leafValueObject)
			case "multi_line":
				// Use a scrollable label widget and a read button
				leafTextObject := widget.NewLabel(node.Value)
				scrollable := container.NewVScroll(leafTextObject)
				readButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
					dialog.ShowInformation(node.Title, node.Value, appCtxt.MainWindow)
				})
				branchContainer.Add(readButton)
				branchContainer.Add(scrollable)
			case "multi_line_with_title":
				// Use a scrollable label widget and a read button
				leafTitleObject := canvas.NewText(node.Title, color.White)
				leafTextObject := widget.NewLabel(node.Value)
				scrollable := container.NewVScroll(leafTextObject)
				readButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
					dialog.ShowInformation(node.Title, node.Value, appCtxt.MainWindow)
				})
				branchContainer.Add(leafTitleObject)
				branchContainer.Add(readButton)
				branchContainer.Add(scrollable)
			default:
				log.Printf("--GUI-- Tree unexpected nodeType: %v", node.NodeType)
			}
		},
	)

	// Expand first-level branches by default
	tree.OpenBranch("Finished")
	tree.OpenBranch("In Progress")
	tree.OpenBranch("Unstarted")

	return tree
}

// Button function
func buttonFuncExpandBranches(tree *widget.Tree, treeData map[string][]string, uid string) {
	// Step 1: Open the current branch
	tree.OpenBranch(uid)

	// Step2: Retrieve child UIDs for this branch
	children, exists := treeData[uid]
	if !exists {
		return
	}

	// Step 3: For each child, call function recursively to open it and all it's children.
	for _, child := range children {
		buttonFuncExpandBranches(tree, treeData, child)
	}
}

// Button function
func buttonFuncCollapseBranches(tree *widget.Tree, treeData map[string][]string, uid string) {
	// Step 1: Retrieve child UIDs for this branch
	children, exists := treeData[uid]
	if !exists {
		return
	}

	// Step 2 : For each child, call function recursively to close all children and close it
	for _, child := range children {
		buttonFuncCollapseBranches(tree, treeData, child)
		tree.CloseBranch(uid)
	}
}

// Button function
func buttonFuncMediumEdit(appCtxt *context.AppContext, node TreeNode, mediaType string, mediaList []models.MediumWithRecord) {
	// Ask if user wants to edit their personal record or the medium itself
	var editDialog dialog.Dialog

	line1 := canvas.NewText("What do you want to edit ?", color.White)
	line1.Alignment = fyne.TextAlignCenter

	mediumEditButton := widget.NewButton("Medium's info", func() {
		appCtxt.PageManager.ShowUpdateMediaPage(mediaType, node.Value, mediaList)
		editDialog.Hide()
	})

	recordEditButton := widget.NewButton("My personal record", func() {
		appCtxt.PageManager.ShowUpdateRecordPage(mediaType, node.Value, mediaList)
		editDialog.Hide()
	})

	editDialog = dialog.NewCustom("Edit Medium", "Cancel", container.NewVBox(
		line1,
		container.NewHBox(layout.NewSpacer(), mediumEditButton, layout.NewSpacer(), recordEditButton, layout.NewSpacer()),
	), appCtxt.MainWindow)

	editDialog.Show()
}

// Button function
func buttonFuncMediumDelete(appCtxt *context.AppContext, node TreeNode) {
	// First dialog : sure to delete ?
	dialog.ShowConfirm("Confirm", fmt.Sprintf("Are you sure you want to delete this medium: %s ?", node.Title), func(b bool) {
		if b {
			// If yes, second dialog : what to delete ?
			line0 := canvas.NewText("What do you want to delete ?", color.White)
			line0.Alignment = fyne.TextAlignCenter
			line1 := canvas.NewText("This medium AND your personal record about it", color.White)
			line1.Alignment = fyne.TextAlignCenter
			line2 := canvas.NewText("OR", color.White)
			line2.Alignment = fyne.TextAlignCenter
			line3 := canvas.NewText("Just your personal record", color.White)
			line3.Alignment = fyne.TextAlignCenter
			line4 := canvas.NewText("(This allows you or other user to retrieve this medium later from server's database)", color.White)
			line4.Alignment = fyne.TextAlignCenter
			dialog.ShowCustomConfirm(
				"Confirm",
				"Medium and Record",
				"Only Record",
				container.NewVBox(line0, line1, line2, line3, line4),
				func(b bool) {
					// Call for delete commands, according to anwser
					// With a third dialog : last warning
					if b { // Medium and record delete
						dialog.ShowConfirm("Last Warning", "Are you sure you want to delete both the medium and your record ?", func(b bool) {
							if b {
								appCtxt.APIClient.Media.DeleteMedium(node.Value) // Deleting a medium will automatically delete linked record on cascade
								dialog.ShowInformation("Info", "Medium and Record deleted !", appCtxt.MainWindow)
								appCtxt.PageManager.ShowHomePage()
							}
						}, appCtxt.MainWindow)
					} else { // Record delete only
						dialog.ShowConfirm("Last Warning", "Are you sure you want to delete your personal record about this medium ?", func(b bool) {
							if b {
								appCtxt.APIClient.Records.DeleteRecord(node.Value)
								dialog.ShowInformation("Info", "Personal Record deleted !", appCtxt.MainWindow)
								appCtxt.PageManager.ShowHomePage()
							}
						}, appCtxt.MainWindow)
					}
				},
				appCtxt.MainWindow)
		}
	}, appCtxt.MainWindow)
}
