package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type TogglClient struct {
	client    *http.Client
	api_token string
}

type Me struct {
	DefaultWorkspaceID int `json:"default_workspace_id"`
}

type TimeEntry struct {
	ID          int      `json:"id"`
	ProjectID   int      `json:"project_id"`
	Start       string   `json:"start"`
	Stop        string   `json:"stop"`
	Duration    int      `json:"duration"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Project struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func NewTogglClient(api_token string) *TogglClient {
	client := &http.Client{}
	return &TogglClient{client, api_token}
}

func (tc *TogglClient) GetMe() Me {
	log.Println("Fetching user data")
	req, err := http.NewRequest(http.MethodGet, "https://api.track.toggl.com/api/v9/me", nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(tc.api_token, "api_token")

	resp, err := tc.client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to perform request: %v", err))
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get time entries: %s", resp.Status)
	}

	var result Me
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode response body: %v", err))
	}

	return result
}

func (tc *TogglClient) GetProjects() []Project {
	me := tc.GetMe()
	log.Println("Fetching projects")

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects", me.DefaultWorkspaceID), nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(tc.api_token, "api_token")

	resp, err := tc.client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to perform request: %v", err))
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get time entries: %s", resp.Status)
	}

	var result []Project
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode response body: %v", err))
	}

	return result
}

func (tc *TogglClient) GetTimeEntries(start *string, end *string) []TimeEntry {
	log.Println(fmt.Sprintf("Fetching time entries from %s to %s", *start, *end))
	params := url.Values{}
	params.Add("start_date", *start)
	params.Add("end_date", *end)

	fullUrl := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", params.Encode())
	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(tc.api_token, "api_token")

	resp, err := tc.client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to perform request: %v", err))
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get time entries: %s", resp.Status)
	}

	var entries []TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode response body: %v", err))
	}

	return entries
}
