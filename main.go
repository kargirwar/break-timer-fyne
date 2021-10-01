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
var ctrls map[string]*TimerCtrl

func main() {
	ch = make(chan *uiCmd)
	ctrls = make(map[string]*TimerCtrl)

	go uiHandler()

	app := app.New()
	window := app.NewWindow("Break Timer")

	gui = container.NewVBox()
	btns := container.NewHBox(
		widget.NewButton("Add rule", func() {
			log.Println("Adding new rule")
			ctrl := NewTimerCtrl()
			ctrls[ctrl.uid()] = ctrl
			gui.Add(ctrl.ui())
		}),

		widget.NewButton("Save rules", func() {
			setRules()
		}),
	)

	gui.Add(btns)

	window.SetContent(gui)
	window.Resize(fyne.NewSize(300, 200))
	window.ShowAndRun()
}

//func runPopUp(w fyne.Window) (modal *widget.PopUp) {
//modal = widget.NewModalPopUp(
//container.NewVBox(
//widget.NewLabel("bar"),
//widget.NewButton("Close", func() { modal.Hide() }),
//),
//w.Canvas(),
//)
//modal.Show()
//return modal
//}

func setRules() {
	for _, v := range ctrls {
		log.Println(v.rule())
	}
}

func uiHandler() {
	for {
		select {
		case c := <-ch:
			log.Println("Cmd received:" + c.cmd)
			id := c.data.(string)
			ctrl, ok := ctrls[id]
			if !ok {
				log.Println("Unable to remove " + id)
				continue
			}
			gui.Remove(ctrl.ui())
			delete(ctrls, id)
		}
	}
}
