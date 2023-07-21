package main

import (
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/driver/desktop"
)

// PointerButton is a button with pointer mouse coursor when enabled.
type PointerButton struct {
	widget.Button
}

// NewPointerButton creates a new button widget with the set label and tap
// handler.
func NewPointerButton(text string, onTapped func()) *PointerButton {
	btn := &PointerButton{}
	btn.ExtendBaseWidget(btn)
	btn.Text = text
	btn.OnTapped = onTapped
	return btn
}

// Cursor returns the cursor type of this widget.
func (b *PointerButton) Cursor() desktop.Cursor {
	if !b.Disabled() {
		return desktop.PointerCursor
	}
	return desktop.CrosshairCursor
}
