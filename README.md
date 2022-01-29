This repo was created to stress the siegfried server in order to reproduce https://github.com/richardlehane/siegfried/issues/172

It creates a corpus of synthetic files, half of which are empty, and then sends requests from 16 goroutines at the sf server. The corpus files are made in a directory called "corpus" that is relative to wherever you run the program from. Please clean that directory up yourself after the run.

To run this program you need to give the server address for your siegfried server (e.g. localhost:8080) and the size of your corpus (e.g. 1024). The corpus size must be at least 8.

E.g. `stresssf localhost:8080 1024`

I've found I can reliably crash the server on corpus sizes even as low as 248.

This program will do multiple runs (starting from 8 and doubling each time) up until the corpus size you've provided. This is in order to measure the time taken so that I can assess the speed impacts of any fixes I apply to siegfried.

To install this program, do `go install github.com/richardlehane/stresssf@latest`.

