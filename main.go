package main

func main() {
}

func makeCorpus(path string) error {
pdf := make([]byte,25000)
copy(pdf,[]byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E})
for i := 0; i < 250000; i++ {
  f, err := os.Create(
}
}
