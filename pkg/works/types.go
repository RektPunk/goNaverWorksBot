package works

const UserPostURL = "https://www.worksapis.com/v1.0/bots/%s/users/%s/messages"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Source struct {
	UserID    string `json:"userId"`
	ChannelID string `json:"channelId"`
	DomainID  int    `json:"domainId"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type WebhookRequest struct {
	Type       string  `json:"type"`
	Source     Source  `json:"source"`
	IssuedTime string  `json:"issuedTime"`
	Content    Content `json:"content"`
}
