package main

import (
	"esh-cli/pkg/utils"
	"fmt"
)

func main() {
	// Test the specific case from the user's issue
	result := utils.IncrementTag("dev_0.1.0", false)
	fmt.Printf("IncrementTag(\"dev_0.1.0\", false) = \"%s\"\n", result)

	// Test other cases
	result2 := utils.IncrementTag("dev_0.1.0-0", false)
	fmt.Printf("IncrementTag(\"dev_0.1.0-0\", false) = \"%s\"\n", result2)

	result3 := utils.IncrementTag("dev_0.1.0", true)
	fmt.Printf("IncrementTag(\"dev_0.1.0\", true) = \"%s\"\n", result3)
}
