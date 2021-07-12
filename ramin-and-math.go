package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func isSquare(number int) bool {
	// negetive numbers are out of rule here
	if number <= 0 {
		return false
	}

	// the simple way to find a square number is
	// checking the second root is absolout or not
	sqrt := math.Sqrt(float64(number))
	return sqrt == float64(int64(sqrt))
}

func readInput(buffer *bufio.Reader) string {
	text, err := buffer.ReadString('\n')
	// not much of error handling here
	if err != nil {
		panic(err)
	}
	// replacing the endline character
	if runtime.GOOS == "windows" {
		text = strings.Replace(text, "\n\r", "", -1)
	} else {
		text = strings.Replace(text, "\n", "", -1)
	}

	return text
}

func main() {
	// defining a buffer to read the inputs
	buffer := bufio.NewReader(os.Stdin)
	var (
		q     int
		input string
		err   error
	)

	// getting the count
	for q < 1 {
		fmt.Print("How many questions to ask: ")
		input = readInput(buffer)

		// convert the string to an integer
		q, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter an integer.")
		}
	}

	// get the input series
	var i int
	series := make([][2]int, q)
	validFlag := true
	for i < q {
		// display proper message for first time, error time or normal time
		if !validFlag || i < 1 {
			fmt.Print("Please provide first and last numbers of the range with a space between them: ")
			validFlag = true
		} else {
			fmt.Print("Next range: ")
		}

		// get the input series
		input = readInput(buffer)
		inputSeries := strings.Split(input, " ")
		if len(inputSeries) != 2 {
			validFlag = false
			continue
		}

		// validate and replace them one by one
		r, err := strconv.Atoi(inputSeries[0])
		if err != nil {
			validFlag = false
			continue
		}
		series[i][0] = r

		r, err = strconv.Atoi(inputSeries[1])
		if err != nil {
			validFlag = false
			continue
		}
		series[i][1] = r

		// if they tried to manipulate us ;)
		if series[i][0] > series[i][1] {
			series[i][1] = series[i][0]
			series[i][0] = r
		}

		i++
	}

	now := time.Now()

	// now calculate the result and print it out
	for _, inputSeries := range series {
		var a int
		for number := inputSeries[0]; number <= inputSeries[1]; number++ {
			if isSquare(number) {
				a++
			}
		}
		fmt.Println(a)
	}

	fmt.Println("Calculation finished in ", time.Now().Sub(now))
}
