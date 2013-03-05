package main

import (
  "fmt"
  "log"
  "github.com/whee/adn"
)

func main() {
  fmt.Println("does this work?")

  app := &adn.Application{}
  post, err := app.GetPost("", "1") // unauthenticated request
  if err != nil {
      log.Fatal(err)
  }
  fmt.Printf("%s [%v]\n%s\n", post.User.Username, post.CreatedAt, post.Text)
}
