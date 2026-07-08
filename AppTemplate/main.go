package sweetjuice

type AppService struct{}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) SayHello(name string) (map[string]string, error) {
	if name == "" {
		name = "World"
	}

	return map[string]string{
		"message": "Hello, " + name + "!",
	}, nil
}

func (s *AppService) Ping() (map[string]string, error) {
	return map[string]string{
		"status": "alive",
	}, nil
}

func (s *AppService) RequestPermissions() (map[string]string, error) {
	return map[string]string{
		"status": "alive",
	}, nil
}
