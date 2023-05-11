package model

import "time"

type User struct {
	ID        string    `json:"id" db:"user_id"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Prime struct {
	//  **prime** is a token, and the client developer third party can use it to issue credential.
	Prime string `json:"prime"`
	//  **user** is the prime owner, who has authorized the client developer third party to issue crendential by prime
	User *User `json:"user"`
}

type Credential struct {
	//  **user** is the credential owner
	User *User `json:"user"`
	//  **identityTicket** is a short-term and single-use login credential
	IdentityTicket string `json:"identityTicket"`
	//  **expiredAt** is the expire time of the identityTicket
	ExpiredAt time.Time `json:"expiredAt"`
}
