package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aimichen/qubic-prime-bind-server/creator/admin"

	"github.com/getamis/sirius/log"
)

type SavePrimeFn func(ctx context.Context, memberId, prime string) error
type GetPrimeFn func(ctx context.Context, memberId string) (string, error)

type Handler struct {
	serverUrl   string
	apiKey      string
	apiSecret   string
	savePrimeFn SavePrimeFn
	getPrimeFn  GetPrimeFn
}

func NewHandler(serverUrl, apiKey, apiSecret string,
	savePrimeFn SavePrimeFn, getPrimeFn GetPrimeFn,
) *Handler {
	return &Handler{
		serverUrl:   serverUrl,
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		savePrimeFn: savePrimeFn,
		getPrimeFn:  getPrimeFn,
	}
}

func setHeaders(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Handler) PrimeBind(w http.ResponseWriter, req *http.Request) {
	setHeaders(w, req)

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
		return
	}
	bindTicket := req.FormValue("bindTicket")
	if bindTicket == "" || len(bindTicket) > 32 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: invalid bindTicket [%s]\n", bindTicket)
		return
	}

	ctx := req.Context()
	adminClient := admin.NewClient(h.serverUrl, h.apiKey, h.apiSecret)
	primeResp, err := adminClient.PrimeGet(ctx, bindTicket)
	if err != nil {
		log.Warn("admin prime failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: PrimeGet api failed, err: %s\n", err)
		return
	}
	log.Trace("admin prime", "primeResp", *primeResp, "user address", primeResp.User.Address)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true}`)

	err = h.savePrimeFn(ctx, memberId, primeResp.Prime)
	if err != nil {
		log.Warn("savePrimeFn Failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: savePrimeFn failed, err: %s\n", err)
		return
	}
	log.Info("prime saved", "memberId", req.FormValue("memberId"), "prime", primeResp.Prime)
}

func (h *Handler) CredentialIssue(w http.ResponseWriter, req *http.Request) {
	setHeaders(w, req)

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
		return
	}

	ctx := req.Context()
	prime, err := h.getPrimeFn(ctx, memberId)
	if err != nil {
		log.Warn("getPrimeFn Failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: getPrimeFn failed, err: %s\n", err)
		return
	}

	adminClient := admin.NewClient(h.serverUrl, h.apiKey, h.apiSecret)
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
	log.Info("cridential get", "ticket", credential.IdentityTicket, "memberId", req.FormValue("memberId"), "prime", prime)
}
