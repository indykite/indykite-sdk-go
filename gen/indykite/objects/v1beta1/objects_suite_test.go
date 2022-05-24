package objectsv1beta1_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestObjects(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Objects Suite")
}
