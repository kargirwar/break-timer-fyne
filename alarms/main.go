package alarms

import (
	"time"
	"github.com/kargirwar/break-timer/player"
	"github.com/kargirwar/break-timer/types"
	log "github.com/sirupsen/logrus"
)

//receive new rules on timerCh, communicate with player on playerCh, use initial rules
func Start(timerCh chan []types.Rule, rules []types.Rule) {
	log.Debug("Starting alarm thread")

	var alarms map[string]map[int][]int
	alarms = getAlarms(rules)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case rules = <-timerCh:
			log.Debug(rules)
			alarms = getAlarms(rules)
			log.Debug(alarms)

		case t := <-ticker.C:
			//log.Printf("Current time: %s %d %d", t.Weekday(), t.Hour(), t.Minute())
			log.WithFields(log.Fields{
				"day": t.Weekday(),
				"hour": t.Hour(),
				"minute": t.Minute(),
			}).Debug()
			for _, m := range alarms[t.Weekday().String()][t.Hour()] {
				if m == t.Minute() {
					player.Play(10, 500)
				}
			}
		}
	}
}

//for each day, for each hour find the minutes at which alarm should be sounded
func getAlarms(rules []types.Rule) map[string]map[int][]int {
	var alarms = make(map[string]map[int][]int)
	for _, r := range rules {
		hours := make(map[int][]int)
		for _, d := range r.Days {
			alarms[d] = hours
			s := r.From
			e := r.To
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
						log.Printf("%s h: %v mins: %v\n", d, h, mins)
						i++
						break
					}
				}
				if e == hrs[i] {
					//if the alarm falls exactly on the end hour we should play it
					if m%60 == 0 {
						alarms[d][e] = []int{0}
						log.Printf("%s h: %v mins: %v\n", d, e, alarms[d][e])
					}
					break
				}
			}
		}
	}
	return alarms
}

//get all hours from start to end , both inclusive
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

