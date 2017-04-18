package shed

import (
	"fmt"
	"log"

	cli "gopkg.in/urfave/cli.v1"
)

func rerr(c *cli.Context, err error) {
	if c.Bool("debug") {
		log.Fatal(err)
	} else {
		fmt.Printf("==> error: %s\n", err.Error())
		return
	}
}
