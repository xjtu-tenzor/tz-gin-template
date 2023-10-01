package service

type Service struct {
	Hello
}

func New() *Service {
	service := &Service{}
	return service
}
