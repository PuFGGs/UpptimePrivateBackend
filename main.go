package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	r.HandleFunc("/raw/{owner}/{repo}/{rest:.*}", rawHandler).Methods("GET")

	r.HandleFunc("/repos/{owner}/{repo}/{additional}", reposHandler).Methods("GET")

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]
	rest := vars["rest"]

	githubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", owner, repo, rest)

	req, err := http.NewRequest("GET", githubURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addAuthHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]
	additional := vars["additional"]

	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/%s", owner, repo, additional)

	queryParams := url.Values{}
	for key, values := range r.URL.Query() {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	githubURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequest("GET", githubURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addAuthHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func addAuthHeader(req *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	// Set User-Agent to identify your application
	req.Header.Set("User-Agent", "YourAppName")
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
