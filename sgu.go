package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type SguEntry struct {
	Date     string
	Project  string
	Category string
	Activity string
	CardKey  string
	Hours    float64
	UserName string
}

type entryKey struct {
	project     string
	date        string
	description string
}

func FromTogglTimeEntry(toggl_entry *TimeEntry, project_map map[int]string) SguEntry {
	date, err := time.Parse(time.RFC3339, toggl_entry.Start)
	project := project_map[toggl_entry.ProjectID]

	if err != nil {
		panic(fmt.Sprintf("Failed to parse date: %v", err))
	}

	return SguEntry{
		Date:     date.Format("02/01/2006"),
		Project:  project,
		Category: toggl_entry.Tags[0],
		Activity: toggl_entry.Description,
		CardKey:  "",
		Hours:    float64(toggl_entry.Duration) / 3600,
		UserName: os.Getenv("SGU_USER"),
	}
}

func WriteCsv(entries *[]SguEntry) {
	filename := fmt.Sprintf("report-%s.csv", time.Now().Format("20060102"))
	log.Println("Writing to file: ", filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Comma = ';'

	header := []string{"DATA", "PROJETO", "CATEGORIA", "ATIVIDADE", "CARD_KEY", "HORAS", "USERNAME"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}

	for _, entry := range *entries {
		record := []string{
			entry.Date,
			entry.Project,
			entry.Category,
			entry.Activity,
			entry.CardKey,
			strings.Replace(fmt.Sprintf("%.2f", entry.Hours), ".", ",", 1),
			entry.UserName,
		}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Failed to write record: %s", err)
		}
	}
}
