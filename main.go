package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const groutines = 16

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Give a server and corpus size e.g. stresssf localhost:8080 2048")
	}
	corpusSz, err := strconv.Atoi(os.Args[2])
	if err != nil || corpusSz < 8 {
		log.Fatal("Second argument (corpus size) must be an integer  of at least 8 e.g. stresssf localhost:8080 2048")
	}
	wg := &sync.WaitGroup{}
	urls := make(chan string, groutines)
	log.Println("Starting goroutines")
	for i := 0; i < groutines; i++ {
		go listen(urls, wg)
	}
	log.Println("Making corpus")
	if err := makeCorpus(corpusSz); err != nil {
		log.Fatal(err)
	}
	// Do multiple runs up to the size of the corpus
	for i := 8; i <= corpusSz; i = i * 2 {
		start := time.Now()
		for j := 0; j < corpusSz; j++ {
			wg.Add(1)
			urls <- makeUrl(os.Args[1], filepath.Join("corpus", fmt.Sprintf("%d.pdf", j)))
		}
		wg.Wait()
		elapsed := time.Since(start)
		log.Printf("Run with corpus size %d took %s", i, elapsed)
	}
	close(urls)
}

func makeUrl(server, file string) string {
	return "http://" + server + "/identify/" + url.QueryEscape(file) + "?format=json"
}

func makeCorpus(corpusSz int) error {
	pdf := make([]byte, 25000)
	copy(pdf, []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34})
	copy(pdf[len(pdf)-5:], []byte{0x25, 0x25, 0x45, 0x4F, 0x46})
	if err := os.MkdirAll("corpus", 0777); err != nil {
		return err
	}
	for i := 0; i < corpusSz; i++ {
		if i%2 == 0 {
			if err := os.WriteFile(filepath.Join("corpus", fmt.Sprintf("%d.pdf", i)), pdf, 0666); err != nil {
				return err
			}
		} else {
			f, err := os.Create(filepath.Join("corpus", fmt.Sprintf("%d.pdf", i)))
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}

func listen(c chan string, wg *sync.WaitGroup) {
	client := retryablehttp.NewClient()
	client.Logger = nil
	for url := range c {
		resp, err := client.Get(url)
		if err == nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		} else {
			log.Fatal(err)
		}
		wg.Done()
	}
}
