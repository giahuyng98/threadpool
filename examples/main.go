package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/giahuyng98/threadpool"
)

func main() {
	filePath := flag.String("f", "test.txt", "file path to read from")
	poolSize := flag.Int("n", 1, "pool size")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	pool := threadpool.New(*poolSize, 100)
	defer pool.Join()

	s := bufio.NewScanner(file)
	for s.Scan() {
		link := s.Text()
		pool.AddTask(&TaskGetUrl{url: link})
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
