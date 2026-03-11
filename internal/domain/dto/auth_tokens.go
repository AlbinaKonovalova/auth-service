package dto

type LoginResult struct {
	AccessToken string
	User        CurrentUser
}

type RefreshResult struct {
	AccessToken string
	User        CurrentUser
}
