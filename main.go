package main

const corpusSz = 250000

const groutines = 16

func main() {
  if len(os.Args) < 2  {
    log.Fatal("Give a server e.g. localhost:8080")
  }
  urls := make(chan string, groutines)
  for i := 0; i < groutines; i++ {
    go listen(urls)
  }
  if err := makeCorpus(); err != nil {
    log.Fatal(err)
  }
  for i := 0; i < corpusSz; i++ {
    c <- makeUrl(os.Args[1], filepath.Join("corpus", fmt.Stringf("%d.pdf", i)))
  }
  close(c)
}

func makeUrl(server, file string) string {
  return strings.Join("http://", server, "/identify/", url.QueryEscape(file), "?base64=true&format=json")
}

func makeCorpus() error {
  pdf := make([]byte,25000)
  copy(pdf,[]byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E})
  if err := os.MkdirAll("corpus", 0777); err != nil {
    return err
  }
  for i := 0; i < corpusSz; i++ {
    if i % 2 == 0 {
      if err := os.WriteFile(filepath.Join("corpus", fmt.Stringf("%d.pdf", i)), pdf, 0666); err != nil {
        return err
      }
    } else {
      f, err := os.Create(filepath.Join("corpus", fmt.Stringf("%d.pdf", i)))
      if err != nil {
        return err
      }
      f.Close()
    }
  }
}

func listen(c chan string) {
  for url := range c {
    resp, err := http.Get(url) 
    if err != nil {
      i, _ := io.Copy(ioutil.Discard, resp.Body)
      log.Printf("Got %s. Dumped %d bytes\n", url, i)
      resp.Body.Close()
    }
  }
}
