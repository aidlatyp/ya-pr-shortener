package usecase

type ServicePinger interface {
	Ping() error
}

type Liveliness struct {
	service ServicePinger
}

func NewLiveliness(p ServicePinger) *Liveliness {
	return &Liveliness{
		service: p,
	}
}

func (l *Liveliness) Do() error {
	err := l.service.Ping()
	if err != nil {
		// figure out what to do with error
		// deletionListener somehow
		return err
	}
	return nil
}
