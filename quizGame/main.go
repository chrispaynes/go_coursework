package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Quiz struct {
	correctAnswerCount int
	length             *int
	questions          [][]string
	selectedInts       *[]int
	source             *os.File
	timeLimit          *int
	timer              *time.Timer
}

func NewQuiz(defaultLength, defaultLimeLimit int, defaultSource string) *Quiz {
	q := Quiz{}
	q.correctAnswerCount = 0
	q.length = flag.Int("questionCount", defaultLength, "the amount of quiz questions")
	q.timeLimit = flag.Int("timeLimit", defaultLimeLimit, "the amount of time (in seconds) to answer a question")
	csvFilename := flag.String("csv", defaultSource, "csv file name")
	flag.Parse()

	q.timer = time.NewTimer(time.Duration(*q.timeLimit) * time.Second)
	q.selectedInts = &[]int{}

	q.readQuestions(*csvFilename)

	return q.randomizeQuestionSelection(time.Now().UnixNano())
}

func main() {
	q := NewQuiz(10, 30, "problems.csv")

	for i, sInt := range *q.selectedInts {
		reader := bufio.NewReader(os.Stdin)
		reader.Buffered()

		fmt.Printf("\nQuestion #%d:\n%s\n", (i + 1), strings.TrimSpace(q.questions[sInt][0]))

		responseCh := make(chan string, 1)

		go func() {
			response, err := reader.ReadString('\n')

			if err != nil {
				log.Fatal(err)
			}

			responseCh <- strings.TrimSpace(response)
		}()

		select {
		case <-q.timer.C:
			q.displayResults()
			return
		case response := <-responseCh:
			answer := strings.TrimSpace(q.questions[sInt][1])

			if strings.Compare(answer, response) == 0 {
				q.correctAnswerCount++
				color.Green("Correct!")
				continue
			}

			color.Red("Incorrect!")
		}
	}

	q.displayResults()
}

func isNotInCollection(i int, ints []int) bool {
	if len(ints) == 0 {
		return true
	}

	for _, n := range ints {
		if i == n {
			return false
		}
	}

	return true
}

func (q *Quiz) displayResults() {
	fmt.Printf("\nYou scored %d points out of %d\n", q.correctAnswerCount, *q.length)
}

func (q *Quiz) randomizeQuestionSelection(seed int64) *Quiz {
	rand.Seed(seed)

	for {
		randomInt := rand.Intn(*q.length)

		if isNotInCollection(randomInt, *q.selectedInts) {
			*q.selectedInts = append(*q.selectedInts, randomInt)
		}

		if len(*q.selectedInts) == *q.length {
			break
		}
	}

	return q
}

func (q *Quiz) readQuestions(f string) {
	var err error
	q.source, err = os.Open(f)

	if err != nil {
		log.Fatal(err)
	}

	defer q.source.Close()

	r := csv.NewReader(q.source)
	q.questions, err = r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}
}
