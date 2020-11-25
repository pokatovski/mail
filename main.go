package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	jsonFile, err := os.Open("500.jsonl")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	type CategoryItemJson struct {
		Url        string   `json:"url"`
		Categories []string `json:"categories"`
	}

	categoriesMap := make(map[string][]string)

	scanner := bufio.NewScanner(jsonFile)
	Category := CategoryItemJson{}
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &Category); err != nil {
			log.Fatal("err Unmarshal ", err)
		}
		if len(Category.Categories) > 0 {
			item := Category.Categories[0]
			categoriesMap[item] = append(categoriesMap[item], Category.Url)
		}
	}

	for category, urls := range categoriesMap {
		fmt.Println("category: ", category)
		if category == "yellow" {
			fileName := fmt.Sprintf("%s.tsv", category)
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}

				fmt.Println("resp.Header", resp.StatusCode)
				body, err := ioutil.ReadAll(resp.Body)

				//fmt.Println(h.Body.Content)
				fmt.Println("url: ", url)
				result := fmt.Sprintf("%s \t %s \t %s \n", url, getTitle(string(body)), getDesc(string(body)))
				_, err = f.Write([]byte(result))
				if err != nil {
					log.Fatal(err)
				}
			}
			_ = f.Close()
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal("err scanner ", err)
	}

	//totalChannel := make(chan int)
	//scanner := bufio.NewScanner(os.Stdin)
	//maxGoroutines := 5
	//done := make(chan struct{})
	//go func(totalChannel <-chan int, done chan<- struct{}) {
	//	var total int
	//	for result := range totalChannel {
	//		total += result
	//	}
	//	fmt.Printf("Total: %d\n", total)
	//	done <- struct{}{}
	//}(totalChannel, done)
	//
	//jobs := make(chan struct{}, maxGoroutines)
	//var wg sync.WaitGroup
	//for scanner.Scan() {
	//	jobs <- struct{}{}
	//	wg.Add(1)
	//	go process(scanner.Text(), totalChannel, jobs, &wg)
	//}
	//if err := scanner.Err(); err != nil {
	//	fmt.Fprintln(os.Stderr, "error:", err)
	//	os.Exit(1)
	//}
	//wg.Wait()
	//close(jobs)
	//close(totalChannel)
	//<-done
}

func getTitle(HTMLString string) (title string) {

	r := strings.NewReader(HTMLString)
	z := html.NewTokenizer(r)

	var i int
	for {
		tt := z.Next()

		i++
		if i > 100 { // Title should be one of the first tags
			return
		}

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <title> tag
			if t.Data != "title" {
				continue
			}

			// fmt.Printf("%+v\n%v\n%v\n%v\n", t, t, t.Type.String(), t.Attr)
			tt := z.Next()

			if tt == html.TextToken {
				t := z.Token()
				title = t.Data
				return
				// fmt.Printf("%+v\n%v\n", t, t.Data)
			}
		}
	}
}
func getDesc(HTMLString string) (title string) {

	r := strings.NewReader(HTMLString)
	z := html.NewTokenizer(r)

	var i int
	for {
		tt := z.Next()

		i++
		if i > 100 { // Title should be one of the first tags
			return
		}

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <title> tag
			if t.Data != "script" {
				continue
			}

			// fmt.Printf("%+v\n%v\n%v\n%v\n", t, t, t.Type.String(), t.Attr)
			tt := z.Next()

			if tt == html.TextToken {
				t := z.Token()
				title = t.Data
				return
				// fmt.Printf("%+v\n%v\n", t, t.Data)
			}
		}
	}
}

//func process(url string, totalChannel chan int, jobs <-chan struct{}, wg *sync.WaitGroup) {
//	resp, err := http.Get(url)
//	if err != nil {
//		panic(err)
//	}
//	defer wg.Done()
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	bodyString := string(body)
//
//	fmt.Printf("Count for %s: %d\n", url, strings.Count(bodyString, "Go"))
//	totalChannel <- strings.Count(bodyString, "Go")
//	<-jobs
//}
