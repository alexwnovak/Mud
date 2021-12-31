package main

import "testing"

func TestThing(t *testing.T) {
	value := thing()

	if value != 5 {
		t.Errorf("Values don't match")
	}
}
