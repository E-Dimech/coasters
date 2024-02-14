package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type Park struct {
    Name string `json:"name"`
}

type Coaster struct {
    Name         string `json:"name"`
    Manufacturer string `json:"manufacturer"`
    Park         Park   `json:"park"`
    Height       int    `json:"height"`
    Speed        int    `json:"speed"`
}

type APIResponse struct {
    Coasters []Coaster `json:"hydra:member"`
}

func queryCoasters(query string) (*APIResponse, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://captaincoaster.com/api/coasters", nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return nil, err
    }

    // Adjust the header to use X-Auth-Token
    req.Header.Add("X-Auth-Token", "76489403-869f-46b6-87e0-bc72e410225e")

    q := req.URL.Query()
    q.Add("name", query)
    req.URL.RawQuery = q.Encode()

    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error executing request:", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Improved error handling: check the response status code
    if resp.StatusCode != http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err == nil {
            log.Printf("API responded with status code %d: %s\n", resp.StatusCode, string(bodyBytes))
        } else {
            log.Printf("API responded with status code %d, but the body could not be read\n", resp.StatusCode)
        }
        return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
    }

    var apiResp APIResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        log.Println("Error decoding response:", err)
        return nil, err
    }

    return &apiResp, nil
}

func main() {
    http.HandleFunc("/search-coasters", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        query := r.URL.Query().Get("query")
        if query == "" {
            http.Error(w, "Query parameter 'query' is missing", http.StatusBadRequest)
            return
        }

        apiResp, err := queryCoasters(query)
        if err != nil {
            http.Error(w, "Failed to query CaptainCoaster API", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(apiResp)
    })

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}