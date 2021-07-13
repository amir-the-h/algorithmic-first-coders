package main

import (
	"fmt"
	"math"
	"time"
)

type Series struct {
	Start int
	End   int
}

func newSeries(s, e int) *Series {
	return &Series{
		Start: s,
		End:   e,
	}
}

func (s *Series) Check() (r int) {
	for i := s.Start; i <= s.End; i++ {
		if IsSquare(i) {
			r++
		}
	}

	return
}

func IsSquare(n int) bool {
	// negetive numbers are out of rule here
	if n <= 0 {
		return false
	}

	// the simple way to find a square number is
	// checking the second root is absolout or not
	sqrt := math.Sqrt(float64(n))
	return sqrt == float64(int64(sqrt))
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
	var (
		q int
	)

	// getting the count
	for q < 1 {
		if !readInput("How many questions to ask: ", &q) {
			return
		}
	}
	now := time.Now()

	// get the input series
	fmt.Println("Input questions: ")
	results := make([]int, q)
	for i := 0; i < q; i++ {
		var s, e int
		if !readInput("\t", &s, &e) {
			return
		}
		series := newSeries(s, e)
		results[i] = series.Check()
	}

	// now calculate the result and print it out
	for _, r := range results {
		fmt.Println(r)
	}

	fmt.Println("Calculation finished in ", time.Now().Sub(now))
}
