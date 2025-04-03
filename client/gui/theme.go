package gui

import "fyne.io/fyne/v2"

type MyTheme struct {
	fyne.Theme
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		if style.Italic {
			return resourceBolditalicTtf
		}
		return resourceBoldTtf
	}
	if style.Italic {
		return resourceItalicTtf
	}
	return resourceRegularTtf
}
