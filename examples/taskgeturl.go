package main

import (
  //"github.com/giahuyng98/threadpool/task"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TaskGetUrl struct {
	url string
}

func (task *TaskGetUrl) Process() error {
	resp, err := http.Get(task.url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  str := string(body)
	fmt.Printf("Get URL: %v => %v\n", task.url, len(str))

	return nil
}
