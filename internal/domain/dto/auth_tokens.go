package dto

// LoginResult — результат login usecase.
type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	User        CurrentUser
}

// RefreshResult — результат refresh usecase.
type RefreshResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
}
