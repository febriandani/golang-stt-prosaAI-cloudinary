package general

type RenewToken struct {
	Token        string `json:"token"`
	TokenExpired string `json:"token_expired"`
}

type JWTAccess struct {
	AccessToken        string      `json:"access"`
	AccessTokenExpired string      `json:"access_expired"`
	RenewToken         string      `json:"renew"`
	RenewTokenExpired  string      `json:"renew_expired"`
	Address            interface{} `json:"address"`
}
