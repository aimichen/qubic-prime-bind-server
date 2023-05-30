package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// #nosec: G101: Potential hardcoded credentials
	MetadataQubicAPIKey string = "X-Qubic-Api-Key"
	// Request timestamp in milliseconds since epoch time
	MetadataQubicTimestamp string = "X-Qubic-Ts"
	// Signature for the request
	MetadataQubicSignature string = "X-Qubic-Sign"
)

func signature(secret, timestamp, httpMethod, resource, body string) string {
	msg := fmt.Sprintf("%s%s%s%s", timestamp, httpMethod, resource, body)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func BuildQubicSigHeader(key string, secret string, endpoint string) func(body string) http.Header {
	return func(body string) http.Header {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		ts := fmt.Sprintf("%d", now)
		parsed, _ := url.Parse(endpoint)
		sig := signature(secret, ts, http.MethodPost, parsed.RequestURI(), body)

		header := http.Header{}
		header.Set(MetadataQubicAPIKey, key)
		header.Set(MetadataQubicTimestamp, ts)
		header.Set(MetadataQubicSignature, sig)

		return header
	}
}
