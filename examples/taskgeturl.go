package main

import (
	//"github.com/giahuyng98/threadpool/task"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TaskGetUrl struct {
	url string
}

func (task *TaskGetUrl) Process() error {
	start := time.Now()
	resp, err := http.Get(task.url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	str := string(body)
	stop := time.Now()

	elapsed := stop.Sub(start)
	fmt.Printf("Get URL: %v => %v KB, elapsed: %v\n", task.url, len(str)/1024, elapsed)

	return nil
}
