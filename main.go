package main

import (
	"bufio"
	"fmt"
	"github.com/mgutz/ansi"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type student struct {
	ID          string
	Name        string
	Retning     string
	Ambition    int
	Gruppe      string
	Prioriteter []string
}

func main() {

	fmt.Println(ansi.Green,`
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$_$$$$$$$$$$$$$$$$_$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$__$$$$$$$$$$$$$$_$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$_______________$$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$___________________$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$____$$$_________$$$____$$$$$$$$$$$$$
$$$$$$$$$$$$$_____$$$_________$$$_____$$$$$$$$$$$$
$$$$$$$$$$$$___________________________$$$$$$$$$$$
$$$$$$$$$$$$___________________________$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
$$$$$____$$$____________________________$$$____$$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$______$$____________________________$$______$$
$$$$$____$$$____________________________$$$____$$$
$$$$$$$$$$$$____________________________$$$$$$$$$$
$$$$$$$$$$$$____________________________$$$$$$$$$$
$$$$$$$$$$$$___________________________$$$$$$$$$$$
$$$$$$$$$$$$$$$$$______$$$$$$_____$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$______$$$$$$_____$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$______$$$$$$_____$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$______$$$$$$_____$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$______$$$$$$_____$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
`,ansi.Reset)
	fmt.Println("Programmet vil vise dig medkursister samme prioriteter som dig,\n" +
		"og lignende ambitionsniveau.")
	data := get_data()
	you := main_flow(data)
	get_possible_groupmates(you, data, 0)

}

func main_flow(data []student) student {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(ansi.Color("Student id: ","blue+b"))
		studienummer, _ := reader.ReadString('\n')
		studienummer = strings.Trim(studienummer, "\n\t\r")
		match, _ := regexp.MatchString(`s\d{6}`, studienummer)
		if match {
			svar, you := findStudent(data, studienummer)
			if you != nil {
				fmt.Println(ansi.Cyan+" ->"+ansi.Reset, svar)
				return *you
			} else {
				fmt.Println(ansi.Cyan+" ->"+ansi.Reset, svar)
			}
		} else {
			fmt.Println(ansi.Cyan + " -> " + ansi.Red + "student id skal skrives: sXXXXXX" + ansi.Reset)
		}
	}
}

func get_amdition_difference(a int, b int, e student) string {
	if a-b >= 10 {
		return ansi.Red + "\t* " + e.ID + " - " + e.Name + ansi.Reset
	} else {
		return ansi.Green + "\t* " + e.ID + " - " + e.Name + ansi.Reset
	}
}

func get_helper(you student, e student) {
	if e.Ambition > you.Ambition {
		fmt.Println(get_amdition_difference(e.Ambition, you.Ambition, e))
	} else {
		fmt.Println(get_amdition_difference(you.Ambition, e.Ambition, e))
	}
}

func get_possible_groupmates(you student, students []student, i int) {
	fmt.Println("\"" + you.Prioriteter[i] + "\"")
	fmt.Println("\nFørste prioritet!")
	for _, e := range students {
		if e.Prioriteter[0] == you.Prioriteter[i] {
			get_helper(you, e)
		}
	}
	fmt.Println("\nAnden prioritet!")
	for _, e := range students {
		if e.Prioriteter[1] == you.Prioriteter[i] {
			get_helper(you, e)
		}
	}
	fmt.Println("\nTredje prioritet!")
	for _, e := range students {
		if e.Prioriteter[1] == you.Prioriteter[i] {
			get_helper(you, e)
		}
	}

}

func findStudent(students []student, student string) (string, *student) {
	for _, e := range students {
		if e.ID == student {
			return ansi.Green + "Registrering påbegyndes, vent venligst" + ansi.Reset, &e
		}
	}
	return ansi.Red + "student med nummer " + student + " fides ikke på listen" + ansi.Reset, nil
}

func get_data() (data []student) {
	url := "" +
		"https://docs.google.com/spreadsheets/" +
		"d/1zPrQqne3SrSG78wCKp52pzEBJ73C9SoNemZJLgOTrGU/" +
		"export?format=csv&" +
		"id=1zPrQqne3SrSG78wCKp52pzEBJ73C9SoNemZJLgOTrGU&" +
		"gid=1666970525"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	for _, e := range strings.Split(string(body), "\n") {
		line := clean_data(e)
		if len(line) > 10 {
			re, _ := regexp.MatchString(`s\d{6}`, line[1])
			if re {
				s := student{
					ID:          line[1],
					Name:        line[2] + " " + line[3],
					Retning:     line[4],
					Ambition:    get_ambition(line[6]),
					Gruppe:      strings.Trim(line[15], "gr"),
					Prioriteter: get_priorities(line[23], line[24], line[25]),
				}
				data = append(data, s)
			}
		}
	}
	return data
}

func get_priorities(p1 string, p2 string, p3 string) (priorities []string) {
	priorities = []string{p1, p2, p3}
	for i, e := range priorities {
		if e == fmt.Sprintf("din %d. prioritet", i+1) {
			priorities[i] = "Ingen"
		} else {
			priorities[i] = strings.Trim(e, "\n\t\r")
		}
	}
	return priorities
}

func clean_data(data string) []string {
	return strings.Split(strings.Trim(data, "\""), ",")
}

func get_ambition(ambition string) int {
	if ambition != "" {
		tal, err := strconv.Atoi(strings.Trim(ambition, "\" "))
		if err != nil {
			return -1
		}
		return tal
	} else {
		return 0
	}
}