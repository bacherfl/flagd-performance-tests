package grpc_test

import (
	schemav1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/util/rand"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	pb "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
)

var _ = Describe("YourGRPCService", func() {

	It("should perform gRPC requests successfully", func() {

		numClients := getNumClients()
		wg := &sync.WaitGroup{}

		wg.Add(numClients)

		for i := 0; i < numClients; i++ {
			go func() {
				defer wg.Done()
				if useHTTP() {
					doHttpRequests(15 * time.Minute)
				} else {
					doRequests(15 * time.Minute)
				}
			}()
		}

		wg.Wait()

	})
})

func getNumClients() int {
	return getEnvVarOrDefault("NUM_CLIENTS", 1)
}

func getWaitTimeBetweenRequests() int {
	return getEnvVarOrDefault("WAIT_TIME_BETWEEN_REQUESTS_MS", 10)
}

func getEnvVarOrDefault(envVar string, defaultValue int) int {
	if envVarValue := os.Getenv(envVar); envVarValue != "" {
		parsedEnvVarValue, err := strconv.ParseInt(envVarValue, 10, 64)
		if err == nil && parsedEnvVarValue > 0 {
			defaultValue = int(parsedEnvVarValue)
		}
	}
	return defaultValue
}

func usePersistentConnection() bool {
	if os.Getenv("USE_PERSISTENT_CONN") == "false" {
		return false
	}
	return true
}

func useHTTP() bool {
	if os.Getenv("USE_HTTP") == "true" {
		return true
	}
	return false
}

func doRequests(duration time.Duration) {
	var conn *grpc.ClientConn
	var grpcClient pb.ServiceClient

	// if we use persistent connections, establish the client connection here and reuse it among the requests
	if usePersistentConnection() {
		conn, grpcClient = establishGrpcConnection()
	}

	waitTimeBetweenRequests := getWaitTimeBetweenRequests()

	end := time.Now().Add(duration)

	for {
		if time.Now().After(end) {
			break
		}
		// if we don't use persistent connections, reestablish the connection in each request
		if !usePersistentConnection() {
			conn, grpcClient = establishGrpcConnection()
		}
		randNumber := rand.Intn(5000)
		resp, err := grpcClient.ResolveString(context.Background(), &schemav1.ResolveStringRequest{
			FlagKey: fmt.Sprintf("color-%d", randNumber),
			Context: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"version": structpb.NewStringValue("1.0.0"),
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		if waitTimeBetweenRequests > 0 {
			<-time.After(time.Duration(waitTimeBetweenRequests) * time.Millisecond)
		}
	}
	conn.Close()
}

func doHttpRequests(duration time.Duration) {
	client := req.C()

	waitTimeBetweenRequests := getWaitTimeBetweenRequests()

	end := time.Now().Add(duration)

	for {
		if time.Now().After(end) {
			break
		}
		randNumber := rand.Intn(5000)

		resp, err := client.R().
			SetBody(map[string]interface{}{
				"flagKey": fmt.Sprintf("color-%d", randNumber),
				"context": map[string]interface{}{
					"version": "1.0.0",
				},
			}).
			Post("http://flagd.flagd-performance-test:80/schema.v1.Service/ResolveString")
		if err != nil {
			log.Fatal(err)
		}

		Expect(err).To(Not(HaveOccurred()))
		Expect(resp).NotTo(BeNil())
		Expect(resp.IsSuccessState()).To(BeTrue())
		if waitTimeBetweenRequests > 0 {
			<-time.After(time.Duration(waitTimeBetweenRequests) * time.Millisecond)
		}
	}
}

func establishGrpcConnection() (*grpc.ClientConn, pb.ServiceClient) {
	conn, err := grpc.Dial("flagd.flagd-performance-test:8013", grpc.WithTransportCredentials(insecure.NewCredentials()))
	Expect(err).NotTo(HaveOccurred())
	client := pb.NewServiceClient(conn)
	return conn, client
}

func TestGRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GRPC Suite")
}
