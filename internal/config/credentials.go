package config

type Credentials struct {
	AccessToken  string
	ClientID     string
	ClientSecret string
	PublicKey    string
	PrivateKey   string
	BaseURL      string
	RealmBaseURL string
}

func (c *Credentials) AuthMethod() AuthMethod {
	if c.AccessToken != "" {
		return AccessToken
	}
	if c.ClientID != "" || c.ClientSecret != "" {
		return ServiceAccount
	}
	if c.PublicKey != "" || c.PrivateKey != "" {
		return Digest
	}
	return Unknown
}
