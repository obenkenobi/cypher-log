package conf_test

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPort(t *testing.T) {
	cv.Convey("Given a mock ServerConf", t, func() {
		mockServerConf := conf.MockServerConf{}
		mockServerConf.On("GetPort").Return("8080")
		expectedPort := "8080"
		cv.Convey(fmt.Sprintf("Calling GetPort returns %v", expectedPort), func() {
			cv.So(mockServerConf.GetPort(), cv.ShouldEqual, expectedPort)
		})
	})
}
