package client

type AccessTokenResponse struct {
	Value    string `json:"access_token"`
	ExpireIn int    `json:"expires_in"`
}
