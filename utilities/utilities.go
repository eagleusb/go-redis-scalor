package utilities

import (
	"fmt"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		panic(err)
	}
}
