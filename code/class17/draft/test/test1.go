package main

import "time"

func main() {
	ticker := time.NewTicker(5 * time.Minute)
	for true {
		select {
		case <-ticker.C:
			// TODO
		}
	}
}
