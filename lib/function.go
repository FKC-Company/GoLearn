package lib

import (
	"fmt"
	"strconv"
)

func parseInt(arg string) int {
	i1, err := strconv.Atoi(arg)
	if err == nil {
		fmt.Println(i1)
	}

	return i1
}
