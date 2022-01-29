package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const corpusSz = 2048

const groutines = 16

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Give a server e.g. localhost:8080")
	}
	wg := &sync.WaitGroup{}
	urls := make(chan string, groutines)
	log.Println("Starting goroutines")
	for i := 0; i < groutines; i++ {
		go listen(urls, wg)
	}
	log.Println("Making corpus")
	if err := makeCorpus(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < corpusSz; i++ {
		wg.Add(1)
		urls <- makeUrl(os.Args[1], filepath.Join("corpus", fmt.Sprintf("%d.pdf", i)))
	}
	log.Println("Waiting")
	wg.Wait()
	close(urls)
}

func makeUrl(server, file string) string {
	return "http://" + server + "/identify/" + url.QueryEscape(file) + "?format=json"
}

func makeCorpus() error {
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
	for url := range c {
		resp, err := retryablehttp.Get(url)
		if err == nil {
			i, _ := io.Copy(ioutil.Discard, resp.Body)
			log.Printf("Got %s. Dumped %d bytes\n", url, i)
			resp.Body.Close()
		} else {
			log.Fatal()
		}
		wg.Done()
	}
}
