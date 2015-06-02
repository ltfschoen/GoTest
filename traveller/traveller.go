// Go Language program to practice using two goroutines
// concurrently to prepare a basic travel itinerary
// Instructions:
// > go build
// > ./traveller

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

const ticket = "Ticket Departing=%s Arriving=%s ID=%i Locations=%s \n"

var prompt = "Enter departure and arrival locations or %s to quit."

type locations struct {
    start string
    end   string
}

type tid struct {
	id	  int
	name  locations
}

func init() {
    if runtime.GOOS == "windows" {
        prompt = fmt.Sprintf(prompt, "Ctrl+Z, Enter")
    } else { // Unix-like
        prompt = fmt.Sprintf(prompt, "Ctrl+D")
    }
}

func main() {
	route := make(chan locations)
	defer close(route)
	itinerary := createItinerary(route)
	defer close(itinerary)
	interact(route, itinerary)
}

func createItinerary(route chan locations) chan tid {
	itinerary := make(chan tid)
	i := 0
	go func() {
		for {
			locationsCombinator := <-route
			id := i + 1
			name := locationsCombinator
			itinerary <- tid{id, name} 
		}
	}()
	return itinerary
}

func interact(route chan locations, itinerary chan tid) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		fmt.Printf("From and To: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var start, end string
		if _, err := fmt.Sscanf(line, "%s %s", &start, &end); err != nil {
			fmt.Fprintln(os.Stderr, "invalid input")
			continue
		}
		route <- locations{start, end}
		fid := <-itinerary
		fmt.Printf(ticket, start, end, fid.id, fid.name)
	}
	fmt.Println()
}