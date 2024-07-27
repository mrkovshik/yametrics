package main

import "os"

func mulfunc(i int) (int, error) {
	return i * 2, nil
}
func main() {
	_, err := mulfunc(2)
	if err != nil {
		os.Exit(1) // want "usage of os.Exit in main at 11:3"
	}
}

func positive() {
	_, err := mulfunc(2)
	if err != nil {
		os.Exit(1)
	}
}
