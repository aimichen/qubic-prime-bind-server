package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getamis/graphql-client"
)

type QubicGraphqlClient struct {
	*graphql.Client
	accToken *string
}

func NewQubicGraphqlClient(svrURL string, key string, secret string, accToken *string, opts ...graphql.ClientOption) *QubicGraphqlClient {
	buildHeaderFunc := func(body string) http.Header {
		header := BuildQubicSigHeader(key, secret, svrURL)(body)
		if accToken != nil {
			header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		}
		return header
	}

	opts = append([]graphql.ClientOption{graphql.WithHTTPClient(DefaultClient)}, opts...)
	opts = append(opts, graphql.WithBuildHeaderFunc(buildHeaderFunc))
	c := graphql.NewClient(svrURL, opts...)
	c.Log = func(s string) {
		fmt.Println(s)
	}

	return &QubicGraphqlClient{
		Client:   c,
		accToken: accToken,
	}
}

func (c *QubicGraphqlClient) SendReqVars(ctx context.Context, body string, out interface{}, vars map[string]interface{}) error {
	req := graphql.NewRequest(body)
	for k, v := range vars {
		req.Var(k, v)
	}
	return c.Client.Run(ctx, req, out)
}
