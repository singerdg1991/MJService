package tikka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hoitek/Kit/retry"
	"github.com/hoitek/Maja-Service/tikka/protobuf"
	"google.golang.org/grpc"
	gcred "google.golang.org/grpc/credentials/insecure"
)

// Tikka Create a new tikka interface
type Tikka interface {
	SendEmail(in *protobuf.EmailRequest) (*protobuf.EmailResponse, error)
}

// Create a new tikka struct
type tikka struct {
	GrpcHost       string
	GrpcPort       int
	GrpcConnection *grpc.ClientConn
	Timeout        time.Duration
}

// Default Create a default tikka
var Default = Tikka(nil)

// Setup creates a new tikka instance
func Setup(tikkaGrpcHost string, tikkaGrpcPort int) (Tikka, error) {
	// Create a connection to the grpc server
	grpcAddress := fmt.Sprintf("%s:%d", tikkaGrpcHost, tikkaGrpcPort)
	grpcConn, err := retry.Get(func() (*grpc.ClientConn, error) {
		return grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(gcred.NewCredentials()))
	}, 2, 3)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Grpc of Tikka")

	// Create an tikka instance
	Default = &tikka{
		GrpcHost:       tikkaGrpcHost,
		GrpcPort:       tikkaGrpcPort,
		GrpcConnection: grpcConn,
		Timeout:        120 * time.Second,
	}
	return Default, nil
}

func (t *tikka) SendEmail(in *protobuf.EmailRequest) (*protobuf.EmailResponse, error) {
	if !strings.Contains(in.EmailEntry.Recipient, "@gmail.com") && !strings.Contains(in.EmailEntry.Recipient, "@yahoo.com") {
		return nil, nil
	}

	// Create client
	client := protobuf.NewEmailServiceClient(t.GrpcConnection)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	// Send email
	return client.SendEmailSMTP(ctx, in)
}
