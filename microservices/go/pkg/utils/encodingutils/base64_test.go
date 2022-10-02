package encodingutils_test

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"

	//"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encoding"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringCanEncodeAndDecode(t *testing.T) {
	cv.Convey("When a string is encoded then decoded", t, func() {
		strToEncode := "asbfsasdfngllsldkjfgklorwigjtposdfijsodfhogsdfughjosadfhjosdfjhgophsjdfovhjsdfohuvnbopsduhbfgi"
		encoded := encodingutils.EncodeBase64([]byte(strToEncode))
		logger.Log.WithField("encoded", string(encoded)).Info()
		decoded, err := encodingutils.DecodeBase64(encoded)
		cv.So(err, cv.ShouldBeNil)
		decodedStr := string(decoded)
		cv.Convey("Expect the decoded string value to match the encoded string", func() {
			cv.So(decodedStr, cv.ShouldEqual, strToEncode)
		})
	})
}

func TestDecodedBase64StrIs32Bytes(t *testing.T) {
	cv.Convey("When a supplied base64 string is decoded", t, func() {
		base64Str := "nOOL3edlo/HjLkKvvmOXp7aVUlQFvqv/Q2Dw9SCr63A="
		decoded, err := encodingutils.DecodeBase64([]byte(base64Str))
		cv.So(err, cv.ShouldBeNil)
		cv.Convey("Expect the decoded value to be 32 bytes", func() {
			cv.So(decoded, cv.ShouldHaveLength, 32)
		})
	})

}
