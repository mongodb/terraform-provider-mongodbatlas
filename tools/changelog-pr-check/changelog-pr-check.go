package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("PR_TITLE", os.Getenv("PR_TITLE"))
	fmt.Println("PR_NUMBER", os.Getenv("PR_NUMBER"))
	fmt.Println("PR_LABELS", os.Getenv("PR_LABELS"))
}
