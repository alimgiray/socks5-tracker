package main

type staticCredentials map[string]string

func NewStaticCredentials() staticCredentials {
	return map[string]string{
		"alim":  "pass",
		"giray": "pass",
	}
}

func (s staticCredentials) Valid(user, password, _ string) bool {
	pass, ok := s[user]
	return ok && password == pass
}
