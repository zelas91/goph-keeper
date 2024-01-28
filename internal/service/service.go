package service

type Service struct {
	*auth
}

func New(services ...func(c *Service)) *Service {
	ctl := &Service{}
	for _, serv := range services {
		serv(ctl)
	}
	return ctl
}
