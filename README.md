pault.ag/go/config
==================

This package allows you define structs which both define the configuration
file format, and command line flags.

```go
package main

import (
	"fmt"
	"os"

	"pault.ag/go/config"
)

type Example struct {
	Option string `flag:"option" description:"This sets all sorts of things"`
	Value  int    `flag:"value"  description:"This is an integer value!"`
}

func main() {
	conf := Example{
		Option: "default",
	}
	flags, err := config.LoadFlags("example", &conf)
    // This will load the RFC822 formatted config from ~/.examplerc and
    // return a flag.FlagSet. The Flags will be populated by defaults from
    // the struct given.
	if err != nil {
		panic(err)
	}
	flags.Parse(os.Args[1:])
	fmt.Printf("option %s, value %d\n", conf.Option, conf.Value)
}
```
