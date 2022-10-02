package cipher_test

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipher"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestPasswordBasedEncryption(t *testing.T) {
	startTimeMilli := time.Now().UnixMilli()
	password := "123242432sgdfdffhgvdfhgdfghdfghetyehgbewrtynbertyuuryt43"

	key, keyDerivationSalt, errKeyDerivation := cipher.DeriveKey([]byte(password), nil)
	logWithTimestamp("Derived key", startTimeMilli)

	keyHash, errKeyHash := cipher.HashKey(key)
	logWithTimestamp("Generated key hash from password", startTimeMilli)

	cv.Convey("When using the derived encryption key and key hash from a password", t, func() {

		cv.Convey("Expect no error from deriving the key", func() { cv.So(errKeyDerivation, cv.ShouldBeNil) })
		cv.Convey("Expect no error from generating the key hash", func() { cv.So(errKeyHash, cv.ShouldBeNil) })

		cv.Convey("Expect the derived key can encrypt and decrypt a message", func() {
			testKeyCanEncryptAndDecrypt(key, startTimeMilli)
			logWithTimestamp("Tested key can encrypt & decrypt", startTimeMilli)
		})
		cv.Convey("Create a newly generated key from the password and generated salt", func() {
			newKey, _, err := cipher.DeriveKey([]byte(password), keyDerivationSalt)
			logWithTimestamp("Derived new key from password", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			cv.Convey("Expect the new key can encrypt and decrypt a message", func() {
				testKeyCanEncryptAndDecrypt(newKey, startTimeMilli)
				logWithTimestamp("Tested new key can encrypt & decrypt", startTimeMilli)

				cv.Convey("Expect the new key will match the key hash", func(c cv.C) {
					isVerified, err := cipher.VerifyKeyHash(keyHash, newKey)
					logWithTimestamp("Compared key hash and new key", startTimeMilli)
					c.So(err, cv.ShouldBeNil)

					c.So(isVerified, cv.ShouldBeTrue)
				})
			})
		})
		cv.Convey("Expect the wrong password will fail verification", func() {
			wrongPassword := "wrongPassword"
			wrongKey, _, err := cipher.DeriveKey([]byte(wrongPassword), nil)
			logWithTimestamp("Derived wrong key", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			isVerified, err := cipher.VerifyKeyHash(keyHash, wrongKey)
			logWithTimestamp("Compared key hash and wrong key", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			cv.So(isVerified, cv.ShouldBeFalse)
		})
	})

}

func testKeyCanEncryptAndDecrypt(key []byte, startTimeMilli int64) {
	messageToEncrypt := "Hello world"

	cypherText, err := cipher.Encrypt(key, []byte(messageToEncrypt))
	logWithTimestamp("Encrypted cypher text", startTimeMilli)
	cv.So(err, cv.ShouldBeNil)

	decrypted, err := cipher.Decrypt(key, cypherText)
	logWithTimestamp("Decrypted cypher text", startTimeMilli)
	cv.So(err, cv.ShouldBeNil)
	decryptedTxt := string(decrypted)

	cv.So(messageToEncrypt, cv.ShouldEqual, decryptedTxt)
}

func logWithTimestamp(msg string, startTimeMilli int64) {
	timeSinceStart := time.Now().UnixMilli() - startTimeMilli
	logger.Log.WithField("timeSinceStartMilli", timeSinceStart).Info(msg)
}
