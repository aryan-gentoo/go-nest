package user

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetAllUsers() []string {
	return []string{"Alice", "Bob", "Charlie"}
}
