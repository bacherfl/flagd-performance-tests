package grpc_test

import (
	schemav1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/util/rand"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	pb "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
)

var _ = Describe("YourGRPCService", func() {
	var (
		conn   *grpc.ClientConn
		client pb.ServiceClient
	)

	BeforeSuite(func() {
		var err error
		conn, err = grpc.Dial("flagd.flagd-performance-test:80", grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())
		client = pb.NewServiceClient(conn)
	})

	AfterSuite(func() {
		conn.Close()
	})

	It("should perform a gRPC request successfully", func() {

		randNumber := rand.Intn(5000)

		resp, err := client.ResolveString(context.Background(), &schemav1.ResolveStringRequest{
			FlagKey: fmt.Sprintf("color-%d", randNumber),
			Context: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"version": structpb.NewStringValue("1.0.0"),
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
	})
})

func TestGRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GRPC Suite")
}
