package admin

const (
	primeGetStr = `
	mutation PrimeGet($bindTicket: String!) {
		primeGet(bindTicket: $bindTicket) {
			prime
			user {
				address
				id
				createdAt
				updatedAt
			}
		}
	}`
	credentialIssueStr = `
	mutation CredentialIssue($prime: String!) {
		credentialIssue (prime: $prime) {
			identityTicket
			expiredAt
			user {
				address
				id
				createdAt
				updatedAt
			}
		}
	}`
)
