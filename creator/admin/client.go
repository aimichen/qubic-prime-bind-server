package admin

import (
	"context"
	"simple-creator-client/creator/admin/model"
	"simple-creator-client/utils"

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
