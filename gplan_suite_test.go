package gplan_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGplan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gplan Suite")
}
