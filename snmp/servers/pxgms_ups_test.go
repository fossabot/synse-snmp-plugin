package servers

import (
	"fmt"
	"testing"
)

// TestPxgmsUps is the first PxgmsUps test.
func TestPxgmsUps(t *testing.T) {
	fmt.Printf("TestPxgmUps start\n")

	pxgmsUps, err := NewPxgmsUps()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("pxgmsUps: %+v\n", pxgmsUps)
}
