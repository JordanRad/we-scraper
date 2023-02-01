package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Exercise struct {
	Name        string
	MuscleGroup string
}

var URLs = []string{"https://www.muscleandstrength.com/exercises/abductors.html",
	"https://www.muscleandstrength.com/exercises/abs",
	"https://www.muscleandstrength.com/exercises/adductors.html",
	"https://www.muscleandstrength.com/exercises/biceps",
	"https://www.muscleandstrength.com/exercises/calves",
	"https://www.muscleandstrength.com/exercises/chest",
	"https://www.muscleandstrength.com/exercises/glutes",
	"https://www.muscleandstrength.com/exercises/forearms",
	"https://www.muscleandstrength.com/exercises/hamstrings",
	"https://www.muscleandstrength.com/exercises/lats",
	"https://www.muscleandstrength.com/exercises/lower-back",
	"https://www.muscleandstrength.com/exercises/triceps",
	"https://www.muscleandstrength.com/exercises/neck.html",
	"https://www.muscleandstrength.com/exercises/quads",
	"https://www.muscleandstrength.com/exercises/shoulders",
	"https://www.muscleandstrength.com/exercises/traps",
}

func contains(arr []Exercise, str string) bool {
	for _, e := range arr {
		if strings.Contains(e.Name, str) {
			return true
		}
	}
	return false
}

func toSQLScript(arr []Exercise) {
	f, err := os.Create("insert_exercises.sql")
	if err != nil {
		log.Fatalf("Error creating an output file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString("INSERT INTO exercise_entity (name, muscle_group) VALUES \n")
	if err != nil {
		fmt.Errorf("error writing new string: %w", err)
	}
	for i, ex := range arr {
		if i+1 == len(arr) {
			_, err = f.WriteString(fmt.Sprintf("('%v', '%v')\n", ex.Name, ex.MuscleGroup))
			if err != nil {
				fmt.Errorf("error writing new string: %w", err)
			}
		} else {
			_, err = f.WriteString(fmt.Sprintf("('%v', '%v'),\n", ex.Name, ex.MuscleGroup))
			if err != nil {
				fmt.Errorf("error writing new string: %w", err)
			}
		}
	}
	_, err = f.WriteString("ON CONFLICT DO NOTHING;")
	if err != nil {
		fmt.Errorf("error writing new string: %w", err)
	}
}
func main() {
	c := colly.NewCollector()

	var exercises []Exercise
	fmt.Println("Scrapping loading...")
	for _, url := range URLs {
		c.OnHTML("div.cell", func(e *colly.HTMLElement) {
			mg := strings.Split(url, "exercises/")[1]
			if strings.Contains(mg, ".html") {
				mg = strings.Split(mg, ".html")[0]
			}

			name := e.ChildText("div.node-title>a")
			if !contains(exercises, e.ChildText("div.node-title>a")) {
				ex := Exercise{
					Name:        name,
					MuscleGroup: mg,
				}
				exercises = append(exercises, ex)
			}

		})
		c.Visit(url)
	}
	fmt.Println("URLs have been traversed successfully!")
	var unique []Exercise
	for _, exercise := range exercises {
		if exercise.Name != "" {
			unique = append(unique, exercise)
		}
	}
	fmt.Println("Processing results...")
	toSQLScript(unique)
	fmt.Println(len(unique), " exercises have been processed successfully")

}
