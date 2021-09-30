package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type uiCmd struct {
	index int
}

var ch chan *uiCmd
var gui *fyne.Container
var i int

var m map[int]*fyne.Container

func main() {
	ch = make(chan *uiCmd)
	m = make(map[int]*fyne.Container)

	go uiHandler()

	app := app.New()
	window := app.NewWindow("Break Timer")

	gui = container.NewVBox()
	btn := widget.NewButton("Add rule", func() {
		log.Println("tapped")
		ctrl := NewCtrl()
		m[i] = ctrl(i)
		gui.Add(m[i])
		i++
	})
	gui.Add(btn)

	window.SetContent(gui)
	window.ShowAndRun()
}

func uiHandler() {
	for {
		select {
		case c := <-ch:
			log.Println("Cmd received:" + strconv.Itoa(c.index))
			e := m[c.index]
			if gui == nil {
				log.Println("Returning due to nil")
				return
			}
			gui.Remove(e)
		}
	}
}
