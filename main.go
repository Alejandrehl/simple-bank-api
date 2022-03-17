package main

import (
	"fmt"
	"time"
)

func main() {
    t := time.Now()
    fmt.Println(t.Month())
    fmt.Println(t.Day())
    fmt.Println(t.Year())
}