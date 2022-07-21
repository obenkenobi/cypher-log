package conf

import "github.com/stretchr/testify/mock"

type MockServerConf struct {
	mock.Mock
}

func (s *MockServerConf) GetPort() string {
	args := s.Mock.Called()
	return args.Get(0).(string)
}
