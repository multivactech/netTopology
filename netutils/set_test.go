package netutils

import "testing"

func TestSet_AddAll(t *testing.T) {
	set := &Set{}
	set.AddAll(
		[]string{
			"12",
			"123",
		},
		[]string{
			"fwaef",
			"fewffff",
		},
	)
	set.Print()
}
