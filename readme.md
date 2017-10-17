## Tamework - an go telegram bot api framework [![Go Report Card](https://goreportcard.com/badge/github.com/zhuharev/tamework)](https://goreportcard.com/report/github.com/zhuharev/tamework) [![Coverage Status](https://coveralls.io/repos/github/zhuharev/tamework/badge.svg?branch=master)](https://coveralls.io/github/zhuharev/tamework?branch=master) [![GoDoc](https://godoc.org/github.com/zhuharev/tamework?status.svg)](http://godoc.org/github.com/zhuharev/tamework)

## Instalation

First you need [go language](https://golang.org/dl/).

For install from command just use:

`go get -u github.com/zhuharev/tamework`

## Usage

```go
package main

import (
  "github.com/zhuharev/tamework"
)

func main() {
  // put your bot token here
  token := ""

  tw, err := tamework.New(token)
  if err!=nil {
    panic(err)
  }

  // handler for /start command
  tw.Text("/start", func(ctx *tamework.Context) {
    ctx.Text("hello")
  })

  // wait updates
  tw.Run()
}
```

## Killer feature

Now you don't need FSM and other patterns for storing user input state. With `tamework` you can receive multiple messages in single handler:

```go
func handler(ctx *tamework.Context) {

  var (
    name = ""
    age = "0"

    // The time after which the user's choice or input is reset
    inputTimeout = 30 * time.Second
  )

  ctx.Send("Input your name:")

  update, done := ctx.Wait(inputTimeout)
  if !done {
    // timeout or connection error
    c.Send("timeout! try again")
    return
  }

  name = update.Text()

  ctx.Send("Input your age:")

  update, done := ctx.Wait(inputTimeout)
  if !done {
    // timeout or connection error
    c.Send("timeout! try again")
    return
  }

  age = update.Text()

  ctx.Markdown(fmt.Sprintf("your name: *%s*\n\nyour age: *%s*", name, age))
}
```
