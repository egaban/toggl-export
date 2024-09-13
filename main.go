package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Maps a project id into the project name.
func getProjectMap(client *TogglClient) map[int]string {
	projects := client.GetProjects()
	result := make(map[int]string)
	for _, project := range projects {
		result[project.ID] = project.Name
	}
	return result
}

func groupSguEntries(entries *[]SguEntry) []SguEntry {
	group := make(map[entryKey]SguEntry)

	for _, entry := range *entries {
		key := entryKey{entry.Project, entry.Date, entry.Activity}

		var updated_entry SguEntry
		if sgu_entry, found := group[key]; !found {
			updated_entry = entry
		} else {
			sgu_entry.Hours = sgu_entry.Hours + entry.Hours
			updated_entry = sgu_entry
		}

		group[key] = updated_entry
	}

	result := make([]SguEntry, 0, len(group))
	for _, entry := range group {
		result = append(result, entry)
	}
	return result
}

func main() {
	toggl_key := os.Getenv("TOGGL_KEY")

	if toggl_key == "" {
		log.Fatal("Environment variable TOGGL_KEY is not set")
	}

	start_date := flag.String("start-date", "", "Start date in YYYY-MM-DD format")
	end_date := flag.String("end-date", "", "End date in YYYY-MM-DD format")

	// Parse the flags
	flag.Parse()

	// Validate that both flags are provided
	if *start_date == "" || *end_date == "" {
		fmt.Println("Both --start-date and --end-date must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	client := NewTogglClient(toggl_key)
	project_map := getProjectMap(client)
	toggl_entries := client.GetTimeEntries(start_date, end_date)

	sgu_entries := make([]SguEntry, 0, len(toggl_entries))
	for _, entry := range toggl_entries {
		sgu_entry := FromTogglTimeEntry(&entry, project_map)
		sgu_entries = append(sgu_entries, sgu_entry)
	}

	grouped_entries := groupSguEntries(&sgu_entries)
	WriteCsv(&grouped_entries)
}
