package main

import "os"

func mulfunc(i int) (int, error) {
	return i * 2, nil
}
func main() {
	_, err := mulfunc(2)
	if err != nil {
		os.Exit(1)
	}
}
