package main

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestMyAction(t *testing.T) {
	// Register Gomega's testingT interface
	gomega.RegisterTestingT(t)
}
