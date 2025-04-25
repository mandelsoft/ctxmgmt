package helper

import (
	"fmt"
)

func Output(key string, closure func()) {
	fmt.Printf("// --- begin %s ---\n", key)
	closure()
	fmt.Printf("// --- end %s ---\n", key)
}
