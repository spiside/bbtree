package bbtree

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	n := newNode(nil, nil)

	if n == nil {
		t.Errorf("Could not initialize NewNode()")
	}
}

func TestIsBlueWorks(t *testing.T) {
	n := newNode(nil, nil)

	if !isBlue(n) {
		t.Errorf("Should instantiate to blue node")
	}

	n.parentColor = black
	if isBlue(n) {
		t.Errorf("Should change color to black node")
	}
}
