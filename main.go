package main

import (
	"fmt"

	"github.com/johnsudaar/freebox/client"
)

func main() {

	client, err := client.NewFreeboxClient("mafreebox.free.fr", "com.johnsudaar.test", "0.0.1")
	if err != nil {
		panic(err)
	}

	token, c, err := client.RequestAppToken("test")
	if err != nil {
		panic(err)
	}

	fmt.Println(token)

	a := <-c
	fmt.Println(a)
}
