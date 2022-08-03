package conf_test

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPort(t *testing.T) {
	cv.Convey("Given a mock ServerConf", t, func() {
		mockServerConf := conf.MockServerConf{}
		mockServerConf.On("GetAppServerPort").Return("8080")
		expectedPort := "8080"
		cv.Convey(fmt.Sprintf("Calling GetAppServerPort returns %v", expectedPort), func() {
			cv.So(mockServerConf.GetAppServerPort(), cv.ShouldEqual, expectedPort)
		})
	})
}
