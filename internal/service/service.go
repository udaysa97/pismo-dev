package service

type Service struct {
}

type ServiceOption func(*Service) error

func NewService(opts ...ServiceOption) *Service {
	service := &Service{}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *service as the argument
		err := opt(service)
		if err != nil {
			return nil
		}
	}

	return service
}
