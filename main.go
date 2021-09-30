package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type uiCmd struct {
	cmd  string
	data interface{}
}

var ch chan *uiCmd
var gui *fyne.Container

func main() {
	ch = make(chan *uiCmd)

	go uiHandler()

	app := app.New()
	window := app.NewWindow("Break Timer")

	gui = container.NewVBox()
	btn := widget.NewButton("Add rule", func() {
		log.Println("tapped")
		ctrl := NewCtrl()
		gui.Add(ctrl())
	})
	gui.Add(btn)

	window.SetContent(gui)
	window.ShowAndRun()
}

func uiHandler() {
	for {
		select {
		case c := <-ch:
			log.Println("Cmd received:" + c.cmd)
			e := c.data.(*fyne.Container)
			gui.Remove(e)
		}
	}
}
