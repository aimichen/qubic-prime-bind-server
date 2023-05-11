package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"simple-creator-client/creator/admin"
	"syscall"
	"time"

	"github.com/getamis/sirius/log"
	"github.com/urfave/negroni"
)

const (
	adminAPIKey    = "1263c330-fc03-425f-a3d2-b6a597ed31e4"
	adminAPISecret = "vbYVXtbXPySPTqeLbatfsKrQNZyCghRG"
	adminApiRoute  = "/admin/graphql"
	serverURL      = "https://creator.dev.qubic.app/admin/graphql"
)

func main() {
	mux := http.NewServeMux()
	// $ curl  -X POST --data "bindTicket=ABCbindEFG123ticket&memberId=999" localhost:3000/primeBind
	mux.HandleFunc("/primeBind", func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parseForm failed %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: ParseForm failed\n")
			return
		}

		memberId := req.FormValue("memberId")
		bindTicket := req.FormValue("bindTicket")
		if memberId == "" || len(memberId) > 32 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: invalid memberId [%s]\n", memberId)
		}
		if bindTicket == "" || len(bindTicket) > 32 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: invalid bindTicket [%s]\n", bindTicket)
		}

		ctx := req.Context()
		adminClient := admin.NewClient(serverURL, adminAPIKey, adminAPISecret)
		primeResp, err := adminClient.PrimeGet(ctx, "mock-bind-ticket")
		if err != nil {
			log.Warn("admin prime failed", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: PrimeGet api failed, err: %s\n", err)
			return
		}
		log.Trace("admin prime", "primeResp", *primeResp, "user address", primeResp.User.Address)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"success": true}`)

		log.Info("prime saved", "memberId", req.FormValue("memberId"), "prime", primeResp.Prime)
	})

	// $ curl  -X POST --data "memberId=999" localhost:3000/credentialIssue
	mux.HandleFunc("/credentialIssue", func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parseForm failed %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: ParseForm failed\n")
			return
		}

		memberId := req.FormValue("memberId")
		if memberId == "" || len(memberId) > 32 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: invalid memberId [%s]\n", memberId)
		}

		prime := "UWKHWbtIcimszXMrjokyKsCogFpNXNTA"
		ctx := req.Context()
		adminClient := admin.NewClient(serverURL, adminAPIKey, adminAPISecret)
		credential, err := adminClient.CredentialIssue(ctx, prime)
		if err != nil {
			log.Warn("admin CredentialIssue failed", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: CredentialIssue api failed, err: %s\n", err)
			return
		}
		log.Trace("admin prime", "CredentialIssue", *credential)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"identityTicket": "%s", "expiredAt": "%s", "address": "%s"}`,
			credential.IdentityTicket,
			credential.ExpiredAt.Format(time.RFC3339),
			credential.User.Address,
		)
		log.Info("prime saved", "memberId", req.FormValue("memberId"), "prime", prime)
	})

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	go func() {
		http.ListenAndServe(":3000", n)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigs)

	select {
	case sig := <-sigs:
		fmt.Println("Shutting down", "signal", sig)
		return
	}
}
