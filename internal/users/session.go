package users

type Session struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func newSession(id int, token string) *Session {
	return &Session{
		ID: id,
		Token: token,
	}
}