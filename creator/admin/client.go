package admin

import (
	"context"
	"time"

	"github.com/aimichen/qubic-prime-bind-server/creator/admin/model"
	"github.com/aimichen/qubic-prime-bind-server/utils"

	"github.com/getamis/graphql-client"
)

type client struct {
	qc *utils.QubicGraphqlClient
}

func NewClient(svrURL string, key string, secret string, opts ...graphql.ClientOption) *client {
	return &client{
		qc: utils.NewQubicGraphqlClient(svrURL, key, secret, nil, opts...),
	}
}

func (c *client) PrimeGet(ctx context.Context, bindTicket string) (*model.Prime, error) {
	if bindTicket == "mock-bind-ticket" {
		return &model.Prime{
			Prime: "mock prime",
			User: &model.User{
				ID:      "123456789",
				Address: "0x123456789",
			},
		}, nil
	}

	result := &struct {
		PrimeGet *model.Prime
	}{}
	err := c.qc.SendReqVars(ctx, primeGetStr, result, map[string]interface{}{
		"bindTicket": bindTicket,
	})
	if err != nil {
		return nil, err
	}
	return result.PrimeGet, nil
}

func (c *client) CredentialIssue(ctx context.Context, prime string) (*model.Credential, error) {
	if prime == "mock prime" {
		return &model.Credential{
			User: &model.User{
				ID:      "user1234",
				Address: "0x1234",
			},
			IdentityTicket: "mock identity ticket",
			ExpiredAt:      time.Now(),
		}, nil
	}

	result := &struct {
		CredentialIssue *model.Credential
	}{}
	err := c.qc.SendReqVars(ctx, credentialIssueStr, result, map[string]interface{}{
		"prime": prime,
	})
	if err != nil {
		return nil, err
	}
	return result.CredentialIssue, nil
}
