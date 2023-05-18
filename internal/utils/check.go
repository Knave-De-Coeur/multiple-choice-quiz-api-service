package utils

import "fmt"

// Check takes in a method returning error and logs that error, primarily used for deferred funcs
func Check(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error:", err)
	}
}
