package main

import (
	"fmt"
	api "test/Api"
	storage "test/Storage"
)

func main() {
	fmt.Println("Start")
	sql := storage.New()
	s := api.New(sql)
	s.Run()
}
