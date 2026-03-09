package dto

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	User        CurrentUser
}

type RefreshResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
}
