package cloudfunction

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/aimichen/qubic-prime-bind-server/handler"
	"github.com/getamis/sirius/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var (
	adminAPIKey    string
	adminAPISecret string
	serverURL      string
	gcpProjectId   string

	logClient *logging.Client
)

// curl -m 70 -X POST https://function-1-oad7xixtwq-as.a.run.app?memberId=qqq\&bindTicket=b -H "Authorization: bearer $(gcloud auth print-identity-token)"
func init() {
	adminAPIKey = os.Getenv("adminAPIKey")
	adminAPISecret = os.Getenv("adminAPISecret")
	serverURL = os.Getenv("serverURL")
	gcpProjectId = os.Getenv("gcpProjectId")

	var err error
	logClient, err = logging.NewClient(context.Background(), gcpProjectId)
	if err != nil {
		log.Error("Failed to create client: %v", err)
		return
	}
	defer logClient.Close()
	logger := logClient.Logger("qubic-api-log").StandardLogger(logging.Info)

	h := handler.NewHandler(serverURL, adminAPIKey, adminAPISecret, savePrime, getPrime)
	functions.HTTP("PrimeBind", h.PrimeBind)
	functions.HTTP("CredentialIssue", h.CredentialIssue)

	logger.Println("Server start with key", adminAPIKey)
}

type MemberWithPrime struct {
	MemberId string `firestore:"memberId"`
	Prime    string `firestore:"prime"`
}

func savePrime(ctx context.Context, memberId, prime string) error {
	client, err := firestore.NewClient(ctx, gcpProjectId)
	if err != nil {
		log.Warn("Failed to create firestore client", "err", err)
		return err
	}
	primeDoc := client.Doc(fmt.Sprintf("prime/member_%s", memberId))
	writeResult, err := primeDoc.Create(ctx, &MemberWithPrime{
		MemberId: memberId,
		Prime:    prime,
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil
		}
		log.Warn("Failed to create prime doc", "err", err)
		return err
	}
	log.Info("Successfully store prime", "memberId", memberId, "prime", prime, "time", writeResult.UpdateTime)
	return nil
}

func getPrime(ctx context.Context, memberId string) (string, error) {
	client, err := firestore.NewClient(ctx, gcpProjectId)
	if err != nil {
		log.Warn("Failed to create firestore client", "err", err)
		return "", err
	}
	primeDoc := client.Doc(fmt.Sprintf("prime/member_%s", memberId))
	docsnap, err := primeDoc.Get(ctx)
	if err != nil {
		log.Debug("Failed to get prime", "err", err, "member", memberId)
		return "", err
	}
	var prime MemberWithPrime
	if err := docsnap.DataTo(&prime); err != nil {
		log.Warn("Failed to read data from docsnap", "err", err)
		return "", err
	}
	return prime.Prime, nil
}
