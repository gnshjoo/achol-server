package main

import (
	"fmt"
)

func main() {

	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context){
		fmt.Fprintln(c.ResponseWriter, "welcome")
	})


	s.HandleFunc("GET","/user/:id", logHandler(recoverHandler(func(c *Context){
		if c.Params["id"] == "0" {
			panic("id is zero")
		}
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v\n", c.Params["id"])
	})) )

	s.HandleFunc("POST", "/users", logHandler(recoverHandler(parseFormHandler(parseJsonBodyHandler(func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})))))

	s.Run(":8080")

}
