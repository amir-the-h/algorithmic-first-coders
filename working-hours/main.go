package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var Priorities = map[int]string{
	3: "Station",
	2: "Store",
	1: "Tenant",
}

type Schedule struct {
	Day         string
	OpenHour    int
	OpenMinute  int
	CloseHour   int
	CloseMinute int
}

type Schedules []*Schedule

type Exception struct {
	StartTime *time.Time
	EndTime   *time.Time
	Status    string
}

type Exceptions map[string][]*Exception

type Query struct {
	Time *time.Time
}

type Queries []*Query

func newSchedule(d string, oh, om, ch, cm int) *Schedule {
	return &Schedule{
		Day:         d,
		OpenHour:    oh,
		OpenMinute:  om,
		CloseHour:   ch,
		CloseMinute: cm,
	}
}

func newException(s string, st, et *time.Time) *Exception {
	return &Exception{
		StartTime: st,
		EndTime:   et,
		Status:    s,
	}
}

func newQuery(t *time.Time) *Query {
	return &Query{
		Time: t,
	}
}

func (q *Query) Check(s Schedules, e Exceptions) bool {
	for _, s := range Priorities {
		if exs, ok := e[s]; ok {
			for _, ex := range exs {
				if (q.Time.After(*ex.StartTime) || q.Time.Equal(*ex.StartTime)) &&
					(q.Time.Before(*ex.EndTime) || q.Time.Equal(*ex.EndTime)) {
					return ex.Status == "open"
				}
			}
		}
	}

	for _, sc := range s {
		if q.Time.Weekday().String() == sc.Day {
			ot := time.Date(q.Time.Year(), q.Time.Month(), q.Time.Day(), sc.OpenHour, sc.OpenMinute, q.Time.Second(), q.Time.Nanosecond(), q.Time.Location())
			ct := time.Date(q.Time.Year(), q.Time.Month(), q.Time.Day(), sc.CloseHour, sc.CloseMinute, q.Time.Second(), q.Time.Nanosecond(), q.Time.Location())
			return (q.Time.After(ot) || q.Time.Equal(ot)) && (q.Time.Before(ct) || q.Time.Equal(ct))
		}
	}

	return false
}

func readInput(m string, a ...interface{}) bool {
	fmt.Print(m)
	_, err := fmt.Scan(a...)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func main() {
	// main calculation
	now := time.Now()

	var (
		n, m, q int
	)
	// getting initializing configurations
	if !readInput("Enter configuration inputs: ", &n, &m, &q) {
		return
	}

	// getting schedules one by one
	fmt.Println("Input schedules: ")
	var schedules Schedules
	for len(schedules) < n {
		var (
			d, o, c string
		)
		// get the schedule input
		if !readInput("\t", &d, &o, &c) {
			return
		}

		// split and reformat times => 8:7 -> 08:07
		openTimeSplit := strings.Split(o, ":")
		closeTimeSplit := strings.Split(c, ":")
		openHour, _ := strconv.Atoi(openTimeSplit[0])
		openMinute, _ := strconv.Atoi(openTimeSplit[1])
		closeHour, _ := strconv.Atoi(closeTimeSplit[0])
		closeMinute, _ := strconv.Atoi(closeTimeSplit[1])
		schedule := newSchedule(d, openHour, openMinute, closeHour, closeMinute)
		schedules = append(schedules, schedule)
	}
	fmt.Println("Schedules set")

	// getting exceptions one by one
	fmt.Println("Input exceptions: ")
	exceptions := make(Exceptions, m)
	for i := 0; i < m; i++ {
		var (
			l, st, et, s string
		)
		// get the exception input
		if !readInput("\t", &l, &st, &et, &s) {
			return
		}

		// convert to time
		sTime, err := time.Parse("2006-01-02T15:04", st)
		if err != nil {
			fmt.Println(err)
			return
		}
		eTime, err := time.Parse("2006-01-02T15:04", et)
		if err != nil {
			fmt.Println(err)
			return
		}

		exception := newException(s, &sTime, &eTime)
		currentException := []*Exception{}
		if _, ok := exceptions[l]; !ok {
			exceptions[l] = currentException
		}
		exceptions[l] = append(exceptions[l], exception)
	}
	fmt.Println("Exceptions set", exceptions)

	// getting querys one by one
	fmt.Println("Input queries: ")
	results := make([]bool, q)
	for i := 0; i < q; i++ {
		var (
			t string
		)
		// get the query input
		if !readInput("\t", &t) {
			return
		}

		// convert to time
		tt, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			fmt.Println(err)
			return
		}

		query := newQuery(&tt)
		results[i] = query.Check(schedules, exceptions)
	}

	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Println("Calculation finished in ", time.Now().Sub(now))
}
