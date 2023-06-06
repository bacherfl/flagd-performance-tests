package grpc_test_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GrpcTest Suite")
}
