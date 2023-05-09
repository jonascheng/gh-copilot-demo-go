package main

// Build a web service to respond Github Repos by searching for
// a query and preferred language
// 1. User inputs a query string and preferred language
// 2. Web service respond a table of the top 10 results in JSON format

// Import library
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Define a struct to store the data from Github API
type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
}

// Define a function to get top 10 stared repo data from Github API
// parameter query: indicates user input query string
// parameter lang: indicates user input preferred language
// return: a slice of Repo struct
func getTop10Repos(query string, lang string) []Repo {
	// Define a slice of Repo struct
	var repos []Repo

	// Define a Github API URL
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+language:%s&sort=stars&order=desc", query, lang)

	// Send a GET request to Github API
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Close the response body at the end of the function
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the response body into a map
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Get the items from the map
	items := data["items"].([]interface{})

	// Iterate over the items
	for i, item := range items {
		// Convert the item to a map
		itemMap := item.(map[string]interface{})

		// Get the name, description and stargazers_count from the itemMap
		name := itemMap["name"].(string)
		description := itemMap["description"].(string)
		stargazers := int(itemMap["stargazers_count"].(float64))

		// Create a Repo struct with the name, description and stargazers_count
		repo := Repo{
			Name:        name,
			Description: description,
			Stargazers:  stargazers,
		}

		// Append the repo to the repos slice
		repos = append(repos, repo)

		// Break the loop if we have 10 repos
		if i == 9 {
			break
		}
	}

	// Return the repos slice
	return repos
}

// Define a function to handle the web request request
// parameter query: indicates user input query string
// parameter lang: indicates user input preferred language
// return: a slice of Repo struct
func handler(w http.ResponseWriter, r *http.Request) {
	// Get the query and lang from the URL query string
	query := r.URL.Query().Get("query")
	lang := r.URL.Query().Get("lang")

	// Get the top 10 repos
	repos := getTop10Repos(query, lang)

	// Marshal the repos slice into JSON
	reposJSON, err := json.Marshal(repos)
	if err != nil {
		log.Fatal(err)
	}

	// Write the reposJSON to the response writer
	w.Write(reposJSON)
}

// Define a main function to spin up the web service, and listen on port 8000 by default
func main() {
	// Get the PORT environment variable
	port := os.Getenv("PORT")

	// If the PORT environment variable is empty, set it to 8000
	if port == "" {
		port = "8000"
	}

	// Print the port we are listening on
	fmt.Println("Listening on port", port)

	// Handle the / route with the handler function
	http.HandleFunc("/", handler)

	// Listen on port 8000 and handle requests
	http.ListenAndServe(":"+port, nil)
}
