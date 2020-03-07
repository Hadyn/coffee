package main

import (
	cli "github.com/hadyn/coffee/cli/coffee"
	"github.com/hadyn/coffee/jagex/dbj2"
)

func main() {
	cli.Execute()

	println(dbj2.Sum([]byte("123456789")))
}
