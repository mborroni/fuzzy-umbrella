package cache

//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=cache cache

type cache interface {
	Get(string) ([]byte, error)
	Save(string, []byte, ...int32) error
}

type Repository struct {
	ttl int32
	cache cache
}

func NewRepository(ttl int32, c cache) *Repository {
	return &Repository{
		ttl:   ttl,
		cache: c,
	}
}

func (r *Repository) Get(id string) (string, error) {
	value, err := r.cache.Get(id)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (r *Repository) Save(id string, value string) error {
	return r.cache.Save(id, []byte(value), r.ttl)
}