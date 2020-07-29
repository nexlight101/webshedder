// Package webshedder provides the tools to enable loadshedding forcasts for Port Elizabeth Metro
package webshedder

import (
	"encoding/json"

	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	// Filename1 const for schedules.json - schedules for loadshedding.
	Filename1 string = "data/schedules.json"
	// Filename2 const for schedules.json - areas for loadshedding.
	Filename2 string = "data/areas.json" // A
	// Layout - standard go date format.
	Layout string = "Monday, 2 January 2006" // S
)

// Declare the start date of the schedule
var day time.Time = time.Date(2020, 6, 24, 0, 0, 0, 0, time.UTC)

//Schedule struct
type Schedule struct {
	Date  time.Time `json:"date"`
	Stage string    `json:"stage"`
	Group []Group   `json:"group"`
}

// Group struct
type Group struct {
	Group string
	Times []string
}

// Area struct
type Area struct {
	Group    string   `json:"group"`
	AreaName []string `json:"areaname"`
}

// AreaM maps area to group
var AreaM map[string][]string

// BuildMap builds a map for area lookup
func BuildMap(aX []Area) map[string][]string {
	// Create map variable
	AreaMap := map[string][]string{}
	for _, v := range aX {
		//loop through all area names for this group
		for _, name := range v.AreaName {
			if value, ok := AreaMap[name]; !ok {
				AreaMap[name] = append(AreaMap[name], v.Group)
			} else {
				delete(AreaMap, name) //Delete the key to recreate it
				value = append(value, v.Group)
				AreaMap[name] = value
			}
		}
	}
	return AreaMap
}

// ReadJSON reads the json files for schedules and areas.
func ReadJSON(filename1, filename2 string) ([]Schedule, []Area) {
	// Open the schedule JSON file
	f1, err := os.OpenFile(filename1, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Cannot open json file!: %v\n", err)
	}
	defer f1.Close()
	// Open the area JSON file
	f2, err := os.OpenFile(filename2, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Cannot open json file!: %v\n", err)
	}
	defer f2.Close()
	// Deal with schedules
	scheduleJ := []Schedule{}
	bX := []byte{}
	bX, err = ioutil.ReadAll(f1)
	if err != nil {
		log.Fatalf("Cannot read from json file!: %v\n", err)
	}
	err = json.Unmarshal(bX, &scheduleJ)
	if err != nil {
		log.Fatalf("Cannot unmarshal from schedule json file!: %v\n", err)
	}
	// Deal with areas
	areaJ := []Area{}
	aBX := []byte{}
	aBX, err = ioutil.ReadAll(f2)
	if err != nil {
		log.Fatalf("Cannot read from json file!: %v\n", err)
	}
	err = json.Unmarshal(aBX, &areaJ)
	if err != nil {
		log.Fatalf("Cannot unmarshal from area json file!: %v\n", err)
	}
	return scheduleJ, areaJ
}

// unique removes duplicates from slice
func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// cleanTSlice cleans up the time slice for areas
func cleanTSlice(sX []string) []string {
	tX := make([]string, 0)
	t := make([]string, 0)
	// remove all new lines and add to slice
	for _, time := range sX {
		tmp := strings.Split(time, "\n")
		// Take out leading & trailing spaces
		for _, v := range tmp {
			t = append(t, strings.Trim(v, " "))
		}
		tX = append(tX, t...)
	}
	return tX
}

// SearchTimes finds the times in the schedule per group
func SearchTimes(d *time.Time, st *string, g []string, s []Schedule) []Group {
	sX := make([]string, 0)
	groupX := []Group{}
	group := Group{}
	for _, v := range s {
		if v.Date == *d && v.Stage == *st {
			for _, gr := range v.Group {
				switch len(g) {
				case 1:
					if gr.Group == g[0] {
						sX = append(sX, gr.Times...)
						sX = cleanTSlice(sX)
						group.Group = g[0]
						group.Times = append(group.Times, sX...)
						sX = sX[:0] //clear sX to re-use
						groupX = append(groupX, group)
						return groupX
					}
				case 2:
					if gr.Group == g[0] {
						sX = append(sX, gr.Times...)
						sX = cleanTSlice(sX)
						group.Group = g[0]
						group.Times = append(group.Times, sX...)
						sX = sX[:0] //clear sX to re-use
						groupX = append(groupX, group)
					} else if gr.Group == g[1] {
						sX = append(sX, gr.Times...)
						sX = cleanTSlice(sX)
						group.Group = g[1]
						group.Times = append(group.Times, sX...)
						sX = sX[:0] //clear sX to re-use
						groupX = append(groupX, group)
						return groupX //This implies that g has group entries from small to large
					}
				}
			}
		}
	}
	return groupX
}
