package main

import mantle ".."

func main() {
	// Setting a value
	// mantle.Handle(`["CuminLevel", 0, 1, 0, [0]]`)

	// Invoking a method. Note that the logging methods aren't 100% because of type issues, hence the array
	// mantle.Handle(`["NewID", 10, 11, 0, []]`)

	mantle.Handle(`["NewApp", 10, 11, 12345, []]`)
	mantle.Handle(`["SetState", 10, 11, 12345, [1]]`)
}
