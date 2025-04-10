package gui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
		can be clicked for update/delete medium or record

	- metadata_main
		branch node for metadata
		with a "expand all" button

	- single_line
		leaf value node, with single line of text
		displayed as a single line text

	- single_line_with_title
		leaf value node, with both a title and a value
		displayed as a single line, Title Value

	- multi_line
		leaf value node, with multiple lines of text
		displayed as a multi line text


*/

func createAndPopulateTree(appCtxt *context.AppContext, mediaList []models.MediumWithRecord) *widget.Tree {
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
				leafContainer := container.NewVBox()
				return leafContainer
			}

		},

		// Callback function to update a node's widget with its data
		func(uid string, branch bool, obj fyne.CanvasObject) {
			node := nodes[uid]
			// Update depends on NodeType and on branch/leaf
			if branch {
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
					// Sub title will have medium text
					branchTextObject.TextSize = 14
					branchContainer.Add(branchTextObject)
				case "medium_title":
					// Medium title will have medium text
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
						fmt.Println("UNSTARTED")
						redColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
						branchTextObject.Color = redColor
					}
					branchContainer.Add(branchTextObject)

				case "metadata_main":
					// Metadata title will have medium text
					// and a "expand all" button
					branchTextObject.TextSize = 14
					expandButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropDown), func() {
						ExpandBranches(tree, treeData, node.ID)
					})
					collapseButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameArrowDropUp), func() {
						CollapseBranches(tree, treeData, node.ID)
					})

					branchContainer.Add(branchTextObject)
					branchContainer.Add(collapseButton)
					branchContainer.Add(expandButton)

				default:
					log.Printf("--GUI-- Tree branch's unexpected nodeType: %v", node.NodeType)
				}

			} else {
				// Leaf nodes will use a widget.Label
				leafContainer := obj.(*fyne.Container)
				// Clear the container first
				leafContainer.Objects = nil
				leafContainer.Refresh()

				switch node.NodeType {
				case "single_line":
					// Use a single canvas.Text
					leafTextObject := canvas.NewText(node.Value, color.White)
					leafContainer.Add(leafTextObject)
				case "single_line_with_title":
					// Use two canvas.Text : one for title, one for value
					leafTitleObject := canvas.NewText(node.Title, color.White)
					leafValueObject := canvas.NewText(node.Value, color.White)
					leafContainer.Add(container.NewHBox(leafTitleObject, leafValueObject))
				case "multi_line":
					// Use a scrollable label widget
					leafTextObject := widget.NewLabel(node.Value)
					scrollable := container.NewVScroll(leafTextObject)
					leafContainer.Add(scrollable)
				default:
					log.Printf("--GUI-- Tree leaf's unexpected nodeType: %v", node.NodeType)
				}

			}

		},
	)

	// Expand first-level branches by default
	tree.OpenBranch("Finished")
	tree.OpenBranch("In Progress")
	tree.OpenBranch("Unstarted")

	return tree
}

func ExpandBranches(tree *widget.Tree, treeData map[string][]string, uid string) {
	// Step 1: Open the current branch
	tree.OpenBranch(uid)

	// Step2: Retrieve child UIDs for this branch
	children, exists := treeData[uid]
	if !exists {
		return
	}

	// Step 3: For each child, call function recursively to open it and all it's children.
	for _, child := range children {
		ExpandBranches(tree, treeData, child)
	}
}

func CollapseBranches(tree *widget.Tree, treeData map[string][]string, uid string) {
	// Step 1: Retrieve child UIDs for this branch
	children, exists := treeData[uid]
	if !exists {
		return
	}

	// Step 2 : For each child, call function recursively to close all children and close it
	for _, child := range children {
		CollapseBranches(tree, treeData, child)
		tree.CloseBranch(uid)
	}
}
