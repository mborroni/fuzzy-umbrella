package users

//go:generate mockgen -source=users.go -destination=users_mock.go -package=users authenticator,repository

type authenticator interface {
	Authenticate(*User) (*Session, error)
}

type repository interface{
	Get(*User) error
	Save(*User) (int, error)
}

type Service struct {
	auth authenticator
	repository repository
}

func NewService(a authenticator, r repository) *Service {
	return &Service{
		auth: a,
		repository: r,
	}
}

func (users *Service) Login(user *User) (*Session, error) {
	if err := users.repository.Get(user); err != nil {
		return nil, err
	}
	session, err := users.auth.Authenticate(user)
	if err != nil {
		return nil, err
	}
	return session, nil
}


func (users *Service) Register(user *User) (int, error) {
	return users.repository.Save(user)
}

