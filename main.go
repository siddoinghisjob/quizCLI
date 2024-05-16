package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type item struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func fileOpener(filename *string) *os.File {
	csvFile, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Error opening the CSV file : %s", *filename))
	}
	return csvFile
}

func fileReader(csvfile *os.File, filename *string) [][]string {
	file := csv.NewReader(csvfile)
	data, err := file.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Error reading the CSV file : %s", *filename))
	}
	return data
}

func csvParser(data [][]string) (problems []item) {
	problems = make([]item, len(data))
	for i, d := range data {
		problems[i] = item{
			d[0],
			d[1],
		}
	}
	return
}

func questionReader(problems []item, timer *time.Timer) (score int, timesup bool) {
	answer := make(chan string)
	for i, d := range problems {
		fmt.Printf("Problem #%d : %s =", i+1, d.question)

		go func() {
			var inp string
			fmt.Scanf("%s\n", &inp)
			answer <- inp
		}()

		select {
		case <-timer.C:
			fmt.Println()
			timesup = false
			return
		case ans := <-answer:
			if ans == d.answer {
				score++
			}
		}
	}
	timesup = true
	return
}

func main() {
	filename := flag.String("csv", "problems.csv", "Name of CSV file in {question, answer} format.")
	timeLim := flag.Int("time", 20, "Set time limit for the program.")

	flag.Parse()

	csvfile := fileOpener(filename)
	data := fileReader(csvfile, filename)
	problems := csvParser(data)

	timer := time.NewTimer(time.Duration(*timeLim) * time.Second)
	score, timesup := questionReader(problems, timer)
	if !timesup {
		fmt.Println("Times Up!")
	}
	fmt.Printf("Your score is : %d out %d \n", score, len(problems))
}
