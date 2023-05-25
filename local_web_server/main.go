package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimichen/qubic-prime-bind-server/handler"
	"github.com/getamis/sirius/log"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/urfave/negroni"
)

const (
	serverUrlFlag      = "serverURL"
	adminAPIKeyFlag    = "adminAPIKey"
	adminAPISecretFlag = "adminAPISecret"
)

var (
	adminAPIKey    string
	adminAPISecret string
	serverURL      string

	memberIdToPrime map[string]string
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "run web server",
	Long:  `run primeBind and credentialIssue web server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initService(cmd)
		if err != nil {
			log.Error("Failed to initialize", "err", err)
			return err
		}
		h := handler.NewHandler(serverURL, adminAPIKey, adminAPISecret, savePrime, getPrime)

		mux := http.NewServeMux()
		// $ curl  -X POST --data "bindTicket=mock-bind-ticket&memberId=999" localhost:80/primeBind
		mux.HandleFunc("/primeBind", h.PrimeBind)
		// $ curl  -X POST --data "memberId=999" localhost:80/credentialIssue
		mux.HandleFunc("/credentialIssue", h.CredentialIssue)

		allowHeaders := cors.New(cors.Options{AllowOriginFunc: func(origin string) bool { return true }})

		n := negroni.New(allowHeaders)
		n.Use(negroni.NewLogger())
		n.UseHandler(mux)

		go func() {
			http.ListenAndServe(":80", n)
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		defer signal.Stop(sigs)

		select {
		case sig := <-sigs:
			fmt.Println("Shutting down", "signal", sig)
			return fmt.Errorf("shutdown with signal")
		}
	},
}

func savePrime(ctx context.Context, memberId, prime string) error {
	memberIdToPrime[memberId] = prime
	return nil
}

func getPrime(ctx context.Context, memberId string) (string, error) {
	prime, ok := memberIdToPrime[memberId]
	if !ok {
		return "", fmt.Errorf("prime not found of memberId [%s]", memberId)
	}
	return prime, nil
}

func init() {
	Cmd.Flags().String(serverUrlFlag, "https://creator.qubic.app/admin/graphql", "qubic creator admin api endpoint")
	Cmd.Flags().String(adminAPIKeyFlag, "", "qubi creator admin api key")
	Cmd.Flags().String(adminAPISecretFlag, "", "qubi creator admin api secret")
}

func initService(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	serverURL = viper.GetString(serverUrlFlag)
	adminAPIKey = viper.GetString(adminAPIKeyFlag)
	adminAPISecret = viper.GetString(adminAPISecretFlag)

	memberIdToPrime = make(map[string]string)
	memberIdToPrime["1"] = "UWKHWbtIcimszXMrjokyKsCogFpNXNTA"
	return nil
}

func main() {
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
