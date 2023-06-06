package grpc_test

import (
	schemav1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/util/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	pb "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
)

var _ = Describe("YourGRPCService", func() {

	It("should perform gRPC requests successfully", func() {

		numClients := 1
		if numClientsEnv := os.Getenv("NUM_CLIENTS"); numClientsEnv != "" {
			if parsedNumClients, err := strconv.ParseInt(numClientsEnv, 10, 64); err != nil && parsedNumClients > 0 {
				numClients = int(parsedNumClients)
			}
		}
		wg := &sync.WaitGroup{}

		wg.Add(numClients)

		for i := 0; i < numClients; i++ {
			go func() {
				defer wg.Done()
				doRequests()
			}()
		}

		wg.Wait()

	})
})

func doRequests() {
	conn, err := grpc.Dial("flagd.flagd-performance-test:8013", grpc.WithTransportCredentials(insecure.NewCredentials()))
	Expect(err).NotTo(HaveOccurred())
	client := pb.NewServiceClient(conn)
	for i := 0; i < 10000; i++ {
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
		//<-time.After(10 * time.Millisecond)
	}
	conn.Close()
}

func TestGRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GRPC Suite")
}
