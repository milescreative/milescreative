package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	url_, err := url.Parse("https://www.google.com?hello=world")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url_.String())
	q := url_.Query()
	fmt.Println(q)

	q.Set("hello2", "world2")
	fmt.Println(q)

	fmt.Println(url_.String())

	url_.RawQuery = q.Encode()
	fmt.Println(url_.String())

	r, err := http.NewRequest("POST", url_.String(), nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println(r.URL.String())
	q2 := r.URL.Query()
	fmt.Println(q2)
	q2.Set("hello3", "world3")
	r.URL.RawQuery = q2.Encode()
	fmt.Println(r.URL.String())
	fmt.Println(r.Header)
}
