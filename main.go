package main

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/renishb10/nhl_api/nhlapi"
)

func main() {
	now := time.Now()

	rosterFile, err := os.OpenFile("rosters.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error opening file roasters.txt %v", err)
	}

	defer rosterFile.Close()

	wrt := io.MultiWriter(os.Stdout, rosterFile)

	log.SetOutput(wrt)

	teams, err := nhlapi.GetAllTeams()
	if err != nil {
		log.Fatalf("Error fetching details from NHL api %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(teams))

	results := make(chan []nhlapi.Roster)

	for _, team := range teams {
		go func(team nhlapi.Team) {
			roster, err := nhlapi.GetRosters(team.ID)
			if err != nil {
				log.Fatalf("Error getting team %v", err)
			}

			results <- roster
			wg.Done()
		}(team)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	display(results)

	log.Printf("took %v", time.Now().Sub(now).String())
}

func display(results chan []nhlapi.Roster) {
	for r := range results {
		for _, ros := range r {
			log.Println("------------------------------")
			log.Printf("Name: %d\n", ros.Person.ID)
			log.Printf("Name: %s\n", ros.Person.FullName)
			log.Printf("Name: %s\n", ros.Position.Abbreviation)
			log.Printf("Name: %s\n", ros.JerseyNumber)
			log.Println("------------------------------")
		}
	}
}
