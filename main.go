package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("500.jsonl")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//fmt.Println(string(byteValue))

	type CategoryItemJson struct {
		Url        string   `json:"url"`
		Categories []string `json:"categories"`
	}

	type CategoryItem struct {
		Url        string
		Categories []string
	}

	//todo: make size
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
	//fmt.Println(categoriesMap)
	for category, urls := range categoriesMap {
		fmt.Println("category: ", category)
		fileName := fmt.Sprintf("%s.tsv", category)
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		for _, url := range urls {
			fmt.Println("url: ", url)
			result := fmt.Sprintf("%s \t %s \t %s \n", url, "header", "desc")
			_, err = f.Write([]byte(result))
			if err != nil {
				log.Fatal(err)
			}
		}
		_ = f.Close()
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
