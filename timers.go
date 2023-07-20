package main

import (
	"context"
	"fmt"
	"github.com/kargirwar/golang/utils"
	"time"
)

const STOP_AFTER = 10 //seconds
func start(rules []Rule) {
	utils.Dbg(context.Background(), fmt.Sprintln("Starting timer thread"))

	var alarms map[string]map[int][]int
	alarms = getAlarms(rules)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case rules = <-timerCh:
			utils.Dbg(context.Background(), fmt.Sprintln(rules))
			alarms = getAlarms(rules)
			utils.Dbg(context.Background(), fmt.Sprintln(alarms))

		case t := <-ticker.C:
			utils.Dbg(context.Background(), fmt.Sprintf("Current time: %s %d %d", t.Weekday(), t.Hour(), t.Minute()))
			for _, m := range alarms[t.Weekday().String()][t.Hour()] {
				if m == t.Minute() {
					playerCh <- PLAY
					utils.Dbg(context.Background(), fmt.Sprintf("Playing alarm"))
					time.Sleep(STOP_AFTER * time.Second)
					playerCh <- STOP //stop alarm after STOP_AFTER unconditionally
				}
			}
		}
	}
}

// for each day, for each hour find the minutes at which alarm should be sounded
func getAlarms(rules []Rule) map[string]map[int][]int {
	var alarms = make(map[string]map[int][]int)
	for _, r := range rules {
		hours := make(map[int][]int)
		for _, d := range r.Days {
			alarms[d] = hours
			s := r.Start
			e := r.End
			f := r.Interval
			hrs := getHours(s, e)

			i := 0
			m := f
			h := 0

			for {
				mins := make([]int, 0)
				for {
					mins = append(mins, m%60)
					m += f

					if m-(60*i) >= 60 {

						h = hrs[i]
						if f == 60 && h == s {
							//if the frequency is 60 minutes, do not play on the starting hour
							i++
							break
						}

						alarms[d][h] = mins
						utils.Dbg(context.Background(), fmt.Sprintf("%s h: %v mins: %v\n", d, h, mins))
						i++
						break
					}
				}
				if e == hrs[i] {
					//if the alarm falls exactly on the end hour we should play it
					if m%60 == 0 {
						alarms[d][e] = []int{0}
						utils.Dbg(context.Background(), fmt.Sprintf("%s h: %v mins: %v\n", d, e, alarms[d][e]))
					}
					break
				}
			}
		}
	}
	return alarms
}

// get all hours from start to end , both inclusive
func getHours(s, e int) []int {
	var hrs []int

	for {
		hrs = append(hrs, s)
		s += 1
		if s > e {
			break
		}
	}
	return hrs
}
