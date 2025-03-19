package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "sync"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	f, err := os.Open(filePath)
	array := make([]student, 0, 20)
	if err!=nil {
		panic(fmt.Sprintf("Failed Reading the file: %v", err.Error()))
	}
	defer f.Close()
	// Did it concurrently earlier but had to skip it to pass tests.
	// wg := &sync.WaitGroup{}
	reader := bufio.NewReader(f)
	idx := 0
	var parseLine func(line string)= func(line string) {
		// defer wg.Done()
		sArr := strings.Split(line, ",")
		if len(sArr) < 7 {
			panic("Malformed Record")
		}
		var nInt func(n uint8) int = func(n uint8) int {
			num, err := strconv.Atoi(strings.TrimSpace(sArr[n]))
			if err!=nil {
				panic(fmt.Sprintf("Invalid record %v", line))
			}
			return num
		}
		array = append(array, student{
			firstName: strings.TrimSpace(sArr[0]),
			lastName: strings.TrimSpace(sArr[1]),
			university: strings.TrimSpace(sArr[2]),
			test1Score: nInt(3),
			test2Score: nInt(4),
			test3Score: nInt(5),
			test4Score: nInt(6),
		})
	}

	for {
		line, err := reader.ReadString('\n')
		if err!=nil {
			parseLine(line)
			break
		}
		idx++
		if idx==1 {
			continue
		} 
		// wg.Add(1)
		// go parseLine(line)
		parseLine(line)

	}
	// wg.Wait()
	return array
}

func calculateGrade(students []student) []studentStat {
	studentStatArray := make([]studentStat, 0, 20)
	for _, value := range students {
		avg := func(s student) float32{
			return float32(s.test1Score + s.test2Score + s.test3Score + s.test4Score)/4
		}
		var grade Grade
		finalScore := avg(value)
		switch {
			case finalScore<35:
				grade=F
			case finalScore>=35 && finalScore < 50:
				grade=C
			case finalScore>=50 && finalScore<70:
				grade=B
			case finalScore>=70:
				grade=A 
		}
		stat := studentStat{
			student: value,
			finalScore: finalScore,
			grade: grade,
		}
		studentStatArray = append(studentStatArray, stat)
	}

	return studentStatArray
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var topper studentStat
	if len(gradedStudents)<1 {
		panic("Empty Graded List received")
	}
	for i, value := range gradedStudents {
		if i==0 {
			topper = value
			continue
		}

		if topper.finalScore <= value.finalScore {
			topper = value
		}
	}

	return topper
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	topperMap := make(map[string]studentStat)
	for _, value := range gs {
		topper, ok := topperMap[value.university]
		if !ok {
			topperMap[value.university]=value
			continue
		}
		if topper.finalScore <= value.finalScore {
			topperMap[value.university] = value
		}
	}
	return topperMap
}
