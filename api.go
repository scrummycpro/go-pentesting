package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Define flags for JSON file path and port number
	jsonFilePath := flag.String("json", "data.json", "path to JSON file")
	port := flag.String("port", "8080", "port number")
	flag.Parse()

	// Load JSON file
	data, err := ioutil.ReadFile(*jsonFilePath)
	if err != nil {
		log.Fatal("Error loading JSON file:", err)
	}

	var records []map[string]interface{}
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Fatal("Error unmarshaling JSON data:", err)
	}

	// Define handler function for search endpoint
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q") // Get the search query parameter from the URL
		if query == "" {
			http.Error(w, "Missing search query", http.StatusBadRequest)
			return
		}

		// Perform full-text search
		searchResults := []map[string]interface{}{}
		for _, record := range records {
			if containsKeyword(record, query) {
				searchResults = append(searchResults, record)
			}
		}

		response, err := json.Marshal(searchResults)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	// Start HTTP server
	log.Printf("Server listening on port %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// containsKeyword checks if the search query is present in any field of the record
func containsKeyword(record map[string]interface{}, query string) bool {
	for _, value := range record {
		switch v := value.(type) {
		case string:
			if strings.Contains(strings.ToLower(v), strings.ToLower(query)) {
				return true
			}
		}
	}
	return false
}
