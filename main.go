package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
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
	gs := make([]student, 0, 20)
	if err != nil {
		return nil
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var parseLine func(line string) = func(line string) {
		sArr := strings.Split(line, ",")
		if len(sArr) < 7 {
			return
		}
		var getNumber func(arr []string, n uint8) int = func(arr []string, n uint8) int {
			num, err := strconv.Atoi(strings.TrimSpace(arr[n]))
			if err != nil {
				return 0
			}
			return num
		}
		gs = append(gs, student{
			firstName:  strings.TrimSpace(sArr[0]),
			lastName:   strings.TrimSpace(sArr[1]),
			university: strings.TrimSpace(sArr[2]),
			test1Score: getNumber(sArr, 3),
			test2Score: getNumber(sArr, 4),
			test3Score: getNumber(sArr, 5),
			test4Score: getNumber(sArr, 6),
		})
	}

	idx := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				parseLine(line)
				break
			}
			return nil
		}
		idx++
		if idx == 1 {
			continue
		}
		parseLine(line)

	}
	return gs
}

func calculateGrade(students []student) []studentStat {
	studentStatArray := make([]studentStat, 0, 20)
	for _, value := range students {
		avg := func(s student) float32 {
			return float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4
		}
		var grade Grade
		finalScore := avg(value)
		switch {
		case finalScore < 35:
			grade = F
		case finalScore >= 35 && finalScore < 50:
			grade = C
		case finalScore >= 50 && finalScore < 70:
			grade = B
		case finalScore >= 70:
			grade = A
		}
		stat := studentStat{
			student:    value,
			finalScore: finalScore,
			grade:      grade,
		}
		studentStatArray = append(studentStatArray, stat)
	}

	return studentStatArray
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var topper studentStat
	if len(gradedStudents) < 1 {
		return studentStat{}
	}
	for i, value := range gradedStudents {
		if i == 0 {
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
	findTopperByUniversity := func(gs []studentStat, university string) []studentStat {
		filteredList := make([]studentStat, 0, 20)
		for _, value := range gs {
			if value.university == university {
				filteredList = append(filteredList, value)
			}
		}
		return filteredList
	}
	for _, value := range gs {
		_, ok := topperMap[value.university]
		if !ok {
			topperMap[value.university] = findOverallTopper(findTopperByUniversity(gs, value.university))
		}
	}
	return topperMap
}
