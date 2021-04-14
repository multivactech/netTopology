package main

import (
	"testing"
)

func Test_getHTML(t *testing.T) {
	data := getHTML()
	t.Error(string(data))
}
