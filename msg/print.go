package msg

import (
	"fmt"
)

func Println(title string, msg interface{}) {
	fmt.Println("==========" + title + "==========")
	fmt.Printf("%+v", msg)
	fmt.Println("")
	fmt.Println("==========" + title + "==========")
}
