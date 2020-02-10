package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	citations := readFile()
	for _, v := range citations {
		r := query(v)
		fmt.Printf("\r\n%s\r\n", r)
	}

}

func query(q string) string {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN env variable missing")
	}

	c := &http.Client{}
	websiteEscaped := url.QueryEscape(q)
	base := "https://autocite.citation-api.com/index/json?url="
	query := fmt.Sprintf("%s%s", base, websiteEscaped)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		fmt.Println("Error executing the query")
		return ""
	}
	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", authHeader)

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error from api got %d", resp.StatusCode)
		return ""

	}
	data, _ := ioutil.ReadAll(resp.Body)
	// bodyStr := string(data)

	gt := &getCitation{}

	err = json.Unmarshal(data, gt)
	if err != nil {
		fmt.Println("Error parsing the response")
		return ""
	}

	gtdd := gt.Data.Data
	citation := fmt.Sprintf("%s. %s. %s. Published %s. Accessed %s", gtdd.Website.Title, gtdd.Pubonline.Title, gtdd.Pubonline.URL, fmt.Sprintf("%s %s, %s", gtdd.Pubonline.Month, gtdd.Pubonline.Day, gtdd.Pubonline.Year), fmt.Sprintf("%s %s, %s", gtdd.Pubonline.Monthaccessed, gtdd.Pubonline.Dayaccessed, gtdd.Pubonline.Yearaccessed))
	return citation

}

func readFile() []string {
	urls := make([]string, 1)
	file, err := os.Open("citations.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return urls
}

type getCitation struct {
	Status string `json:"status"`
	Data   struct {
		Data struct {
			Pubonline struct {
				Title         string `json:"title"`
				Day           string `json:"day"`
				Month         string `json:"month"`
				Year          string `json:"year"`
				Inst          string `json:"inst"`
				Dayaccessed   string `json:"dayaccessed"`
				Monthaccessed string `json:"monthaccessed"`
				Yearaccessed  string `json:"yearaccessed"`
				URL           string `json:"url"`
			} `json:"pubonline"`
			Website struct {
				Title string `json:"title"`
			} `json:"website"`
			Contributors []interface{} `json:"contributors"`
			Autocite     struct {
				URL string `json:"url"`
			} `json:"autocite"`
			Pubtype struct {
				Main string `json:"main"`
			} `json:"pubtype"`
			Source string `json:"source"`
		} `json:"data"`
		Display struct {
			PageTitle string `json:"page_title"`
		} `json:"display"`
	} `json:"data"`
}
