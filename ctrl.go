package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dchest/uniuri"
)

func intervals() []string {
	var intervals []string
	for i := 10; i <= 60; i += 5 {
		intervals = append(intervals, strconv.Itoa(i)+" mins")
	}
	return intervals
}

func hours() []string {
	var hours []string
	for i := 0; i <= 24; i += 1 {
		hours = append(hours, strconv.Itoa(i)+" hours")
	}
	return hours
}

func getDays(kind string) []string {
	var s []string
	switch kind {
	case "Weekdays":
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	case "Weekends":
		return []string{"Saturday", "Sunday"}
	}
	return s
}

type TimerCtrl struct {
	id       string
	vbox     *fyne.Container
	interval int
	days     []string
	start    int
	end      int
}

type Rule struct {
	Interval int
	Days     []string
	Start    int
	End      int
}

func NewTimerCtrl() *TimerCtrl {
	return &TimerCtrl{id: uniuri.New()}
}

func (c *TimerCtrl) uid() string {
	return c.id
}

func (r *Rule) String() string {
	return fmt.Sprintf("Interval: %d Days %v Start %d End %d", r.Interval, r.Days, r.Start, r.End)
}

func (c *TimerCtrl) rule() Rule {
	return Rule{Interval: c.interval, Days: c.days, Start: c.start, End: c.end}
}

func (c *TimerCtrl) ui(r Rule) *fyne.Container {
	if c.vbox != nil {
		return c.vbox
	}

	c.vbox = container.NewVBox()
	//top row
	txt := widget.NewLabel("Take a break every: ")
	sel := widget.NewSelect(intervals(), func(value string) {
		r := regexp.MustCompile(" mins")
		v := r.ReplaceAllString(value, "")
		//safe to ignore error here. The strings are all our own
		c.interval, _ = strconv.Atoi(v)
		log.Printf("%s: Set interval to %d", c.id, c.interval)
	})

	if r.Interval > 0 {
		sel.SetSelected(strconv.Itoa(r.Interval) + " mins")
	}

	btn := widget.NewButton("X", func() {
		ch <- &uiCmd{cmd: "remove", data: c.id}
		log.Println("tapped")
	})
	ctrl := container.NewHBox(txt, sel, btn)

	c.vbox.Add(ctrl)

	//middle row
	txt = widget.NewLabel("On: ")
	sel = widget.NewSelect([]string{"Weekdays", "Weekends"}, func(value string) {
		log.Printf("%s: Select days to %s", c.id, value)
		c.days = getDays(value)
	})
	ctrl = container.NewHBox(txt, sel)

	c.vbox.Add(ctrl)

	daysType := getDaysType(r)

	if daysType != "" {
		sel.SetSelected(daysType)
	}

	//bottom row
	txt = widget.NewLabel("From: ")
	sel1 := widget.NewSelect(hours(), func(value string) {
		r := regexp.MustCompile(" hours")
		v := r.ReplaceAllString(value, "")
		start, _ := strconv.Atoi(v)

		//if start >= end {
		//dialog := widget.
		//}

		c.start = start
		log.Printf("%s: Select start to %d", c.id, c.start)
	})

	sel2 := widget.NewSelect(hours(), func(value string) {
		r := regexp.MustCompile(" hours")
		v := r.ReplaceAllString(value, "")
		end, _ := strconv.Atoi(v)
		c.end = end
		log.Printf("%s: Select end to %d", c.id, c.end)
	})
	ctrl = container.NewHBox(txt, sel1, sel2)

	c.vbox.Add(ctrl)

	sel1.SetSelected(strconv.Itoa(r.Start) + " hours")
	sel2.SetSelected(strconv.Itoa(r.End) + " hours")

	return c.vbox
}

func getDaysType(r Rule) string {
	days := r.Days
	if len(days) == 0 {
		return ""
	}

	for _, d := range days {
		if d == "Monday" {
			return "Weekdays"
		}

		if d == "Saturday" {
			return "Weekends"
		}

		if d == "Sunday" {
			return "Weekends"
		}
	}

	return ""
}
