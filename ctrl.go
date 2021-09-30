package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewCtrl() func() *fyne.Container {
	vbox := container.NewVBox()
	return func() *fyne.Container {
		//top row
		txt := widget.NewLabel("Take a break every: ")
		sel := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
			log.Println("Select interval to", value)
		})

		btn := widget.NewButton("X", func() {
			ch <- &uiCmd{cmd: "remove", data: vbox}
			log.Println("tapped")
		})
		ctrl := container.NewHBox(txt, sel, btn)

		vbox.Add(ctrl)

		//middle row
		txt = widget.NewLabel("On: ")
		sel = widget.NewSelect([]string{"Weekdays", "Weekends"}, func(value string) {
			log.Println("Select days to", value)
		})
		ctrl = container.NewHBox(txt, sel)

		vbox.Add(ctrl)

		//bottom row
		txt = widget.NewLabel("From: ")
		sel1 := widget.NewSelect([]string{"0", "1"}, func(value string) {
			log.Println("Set from: ", value)
		})

		sel2 := widget.NewSelect([]string{"0", "1"}, func(value string) {
			log.Println("Set to: ", value)
		})
		ctrl = container.NewHBox(txt, sel1, sel2)

		vbox.Add(ctrl)

		return vbox
	}
}
