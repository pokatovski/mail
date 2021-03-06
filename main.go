package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	m "github.com/keighl/metabolize"
	"io"
	"log"
	"net/http"

	"os"
	"sync"
	"time"
)

// максимальное кол-во запущенных горутин
const maxGoroutines = 10

type MetaData struct {
	Title       string `meta:"og:title"`
	Description string `meta:"og:description,description"`
}
type SyncWriter struct {
	m      sync.Mutex
	Writer io.Writer
}

func (w *SyncWriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.Writer.Write(b)
}

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

	// пробегаемся по файлу и создаем мапу по категориям
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &Category); err != nil {
			log.Fatal("err Unmarshal ", err)
		}
		if len(Category.Categories) > 0 {
			// Берем всегда только первую категорию,
			// Потому что в тестовом файле представлены урлы только с одной категорией
			item := Category.Categories[0]
			categoriesMap[item] = append(categoriesMap[item], Category.Url)
		}
	}


	maxGoroutines := maxGoroutines
	jobs := make(chan struct{}, maxGoroutines)
	// Проходим по всем категориям. Ставим таймаут на ответ в 10 секунд.
	// Если ответа нет, статус не 200 или не можем распарсить содержимое - пропускаем урл.
	for category, urls := range categoriesMap {
		fileName := fmt.Sprintf("%s.tsv", category)
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		wr := &SyncWriter{sync.Mutex{}, f}
		wg := sync.WaitGroup{}
		wgc := sync.WaitGroup{}
		for _, url := range urls {
			wg.Add(1)
			wgc.Add(1)
			jobs <- struct{}{}
			go process(url, jobs, &wgc, &wg, wr)
		}
		wg.Wait()
		wgc.Wait()
		_ = f.Close()
	}
	close(jobs)
	if err := scanner.Err(); err != nil {
		log.Fatal("err scanner ", err)
	}
}

func process(url string, jobs <-chan struct{}, wgc *sync.WaitGroup, wg *sync.WaitGroup, wr *SyncWriter) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	defer wg.Done()
	defer wgc.Done()
	if err != nil {
		log.Println("failed for get response from url: ", url)
		<-jobs
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("bad status from url: ", url)
		<-jobs
		return
	}
	data := &MetaData{}
	// todo: беда с windows-1251
	err = m.Metabolize(resp.Body, data)
	if err != nil {
		log.Println("failed for reading response body, url: ", url)
		<-jobs
		return
	}
	result := fmt.Sprintf("%s \t %s \t %s \n", url, data.Title, data.Description)
	_, err = wr.Write([]byte(result))
	if err != nil {
		log.Fatal(err)
	}
	<-jobs
}
