package main

const (
	authWrite = "w"
	authRead  = "r"
)

func (s Server) Auth(mode string, headers map[string][]string) bool {
	xauthtoken, ok := headers["X-Auth-Token"]
	if !ok {
		return false
	}

	switch mode {
	case authWrite:
		wtokens := append(s.config.Auth.WTokens, s.config.Auth.RWTokens...)
		for _, token := range wtokens {
			if token == xauthtoken[0] {
				return true
			}
		}
	case authRead:
		rtokens := append(s.config.Auth.RTokens, s.config.Auth.RWTokens...)
		for _, token := range rtokens {
			if token == xauthtoken[0] {
				return true
			}
		}
	}
	return false
}
