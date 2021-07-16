package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Employee struct {
	Name        string
	Loc         *time.Location
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
	Meetings    Meetings
}

type Employees map[string]*Employee
type Attenders []*Employee

type Meeting struct {
	Capacity  int
	Duration  int
	Attenders Attenders
	Time      *time.Time
}

type Meetings []*Meeting

func (m *Meeting) Schedule() bool {
	now := time.Now().UTC()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	validFlag := true
	// check every 30 minutes for availability of all attenders for the meeting
	for startTime.Before(endOfDay) {
		validFlag = true
		endTime := startTime.Add(time.Minute * time.Duration(m.Duration))
		for _, attender := range m.Attenders {
			// need to find when they start and end theire work
			startOfWork := time.Date(now.Year(), now.Month(), now.Day(), attender.StartHour, attender.StartMinute, 0, 0, attender.Loc).UTC()
			endOfWork := time.Date(now.Year(), now.Month(), now.Day(), attender.EndHour, attender.EndMinute, 0, 0, attender.Loc).UTC()

			// check if meetings fits in the working hours
			if (!(startTime.After(startOfWork) || startTime.Equal(startOfWork))) ||
				(!(endTime.Before(endOfWork) || endTime.Equal(endOfWork))) {
				validFlag = false
				continue
			}

			// now check if meeting overlaps other meetings of attender
			for _, meeting := range attender.Meetings {
				endMeeting := meeting.Time.Add(time.Minute * time.Duration(meeting.Duration))
				if (startTime.After(*meeting.Time) || startTime.Equal(*meeting.Time)) &&
					(startTime.Before(endMeeting) || startTime.Equal(endMeeting)) {
					validFlag = false
					continue
				}
			}
		}
		if validFlag {
			m.Time = &startTime
			return true
		}
		startTime = startTime.Add(time.Minute * 30)
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

func readInputLine(m string) string {
	fmt.Print(m)
	buffer := bufio.NewReader(os.Stdin)
	line, err := buffer.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return line
}

func newEmployee(n, l, s, e string) *Employee {
	var (
		sh, sm, eh, em, lh, lm int
	)
	_, err := fmt.Sscanf(s, "%02d:%02d", &sh, &sm)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, err = fmt.Sscanf(e, "%02d:%02d", &eh, &em)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sign := 1
	_, err = fmt.Sscanf(l, "+%02d:%02d", &lh, &lm)
	if err != nil {
		_, err = fmt.Sscanf(l, "-%02d:%02d", &lh, &lm)
		sign = -1
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	offset := ((lh * 60) + lm) * 60 * sign
	loc := time.FixedZone(l, offset)

	return &Employee{
		Name:        n,
		Loc:         loc,
		StartHour:   sh,
		StartMinute: sm,
		EndHour:     eh,
		EndMinute:   em,
	}
}

func newMeeting(employees Employees, count, duration int, attendersNames ...string) *Meeting {
	attenders := make(Attenders, len(attendersNames))
	for i, attenderName := range attendersNames {
		attenders[i] = employees[attenderName]
	}

	return &Meeting{
		Capacity:  count,
		Duration:  duration,
		Attenders: attenders,
	}
}

func main() {
	var p, m int
	if !readInput("Input initializing config: ", &p, &m) {
		return
	}

	now := time.Now()

	fmt.Println("Input employees:")
	employees := make(Employees, p)
	for i := 0; i < p; i++ {
		var n, l, s, e string
		if !readInput("\t", &n, &l, &s, &e) {
			return
		}
		employee := newEmployee(n, l, s, e)
		if employee == nil {
			return
		}
		employees[n] = employee
	}
	fmt.Println("Employees set")

	fmt.Println("Input meetings:")
	meetings := make(Meetings, m)
	for i := 0; i < m; i++ {
		var c, d int
		if !readInput("\tMeeting: ", &c, &d) {
			return
		}
		input := strings.ReplaceAll(readInputLine("\tAttenders: "), "\n", "")
		att := strings.Split(input, " ")
		if len(att) < c {
			fmt.Println("Not enough attenders.")
			return
		}
		meeting := newMeeting(employees, c, d, att...)
		meetings[i] = meeting
	}

	for _, meeting := range meetings {
		if !meeting.Schedule() {
			fmt.Println("N/A")
			break
		}
		for _, attender := range meeting.Attenders {
			attender.Meetings = append(attender.Meetings, meeting)
		}
		fmt.Println(meeting.Time.Format("15:04"))
	}

	fmt.Println("Calculation finished in ", time.Now().Sub(now))
}
