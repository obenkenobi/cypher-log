package cipherutils_test

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestPasswordBasedEncryptionWithAES(t *testing.T) {
	startTimeMilli := time.Now().UnixMilli()
	password := "123242432sgdfdffhgvdfhgdfghdfghetyehgbewrtynbertyuuryt43"

	key, keyDerivationSalt, errKeyDerivation := cipherutils.DeriveAESKeyFromPassword([]byte(password), nil)
	logWithTimestamp("Derived key", startTimeMilli)

	keyHash, errKeyHash := cipherutils.HashKeyBcrypt(key)
	logWithTimestamp("Generated key hash from password", startTimeMilli)

	cv.Convey("When using the derived encryption key and key hash from a password", t, func() {

		cv.Convey("Expect no error from deriving the key", func() { cv.So(errKeyDerivation, cv.ShouldBeNil) })
		cv.Convey("Expect no error from generating the key hash", func() { cv.So(errKeyHash, cv.ShouldBeNil) })

		cv.Convey("Expect the derived key can encrypt and decrypt a message", func() {
			testAESKeyCanEncryptAndDecrypt(key, startTimeMilli)
			logWithTimestamp("Tested key can encrypt & decrypt", startTimeMilli)
		})
		cv.Convey("Create a newly generated key from the password and generated salt", func() {
			newKey, _, err := cipherutils.DeriveAESKeyFromPassword([]byte(password), keyDerivationSalt)
			logWithTimestamp("Derived new key from password", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			cv.Convey("Expect the new key can encrypt and decrypt a message", func() {
				testAESKeyCanEncryptAndDecrypt(newKey, startTimeMilli)
				logWithTimestamp("Tested new key can encrypt & decrypt", startTimeMilli)

				cv.Convey("Expect the new key will match the key hash", func(c cv.C) {
					isVerified, err := cipherutils.VerifyKeyHashBcrypt(keyHash, newKey)
					logWithTimestamp("Compared key hash and new key", startTimeMilli)
					c.So(err, cv.ShouldBeNil)

					c.So(isVerified, cv.ShouldBeTrue)
				})
			})
		})
		cv.Convey("Expect the wrong password will fail verification", func() {
			wrongPassword := "wrongPassword"
			wrongKey, _, err := cipherutils.DeriveAESKeyFromPassword([]byte(wrongPassword), nil)
			logWithTimestamp("Derived wrong key", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			isVerified, err := cipherutils.VerifyKeyHashBcrypt(keyHash, wrongKey)
			logWithTimestamp("Compared key hash and wrong key", startTimeMilli)
			cv.So(err, cv.ShouldBeNil)

			cv.So(isVerified, cv.ShouldBeFalse)
		})
	})

}

func TestRandomKeyEncryptionWithAES(t *testing.T) {
	startTimeMilli := time.Now().UnixMilli()
	cv.Convey("When given an randomly generated AES key", t, func() {
		key, err := cipherutils.GenerateRandomKeyAES()
		cv.So(err, cv.ShouldBeNil)
		cv.Convey("Expect the key can encrypt and decrypt", func() {
			testAESKeyCanEncryptAndDecrypt(key, startTimeMilli)
		})
	})
}

func TestVerifyHashSHA256WithSalt(t *testing.T) {
	msg := "New Message"
	msgBytes := []byte(msg)
	msgHash, hashErr := cipherutils.HashWithSaltSHA256(msgBytes)
	cv.Convey("When a message was hashed with SGA256 with an embedded salt", t, func() {
		cv.Convey("Expect no errors when hashing", func() { cv.So(hashErr, cv.ShouldBeNil) })
		cv.Convey("Expect hash will be verified with the correct message", func() {
			correctBytes := []byte(msg)
			verified, err := cipherutils.VerifyHashWithSaltSHA256(msgHash, correctBytes)
			cv.So(err, cv.ShouldBeNil)
			cv.So(verified, cv.ShouldBeTrue)
		})
		cv.Convey("Expect hash will not be verified with the wrong message", func() {
			verified, err := cipherutils.VerifyHashWithSaltSHA256(msgHash, []byte("wrong message"))
			cv.So(err, cv.ShouldBeNil)
			cv.So(verified, cv.ShouldBeFalse)
		})
	})
}

func testAESKeyCanEncryptAndDecrypt(key []byte, startTimeMilli int64) {
	messageToEncrypt := "Hello world"

	cypherText, err := cipherutils.EncryptAES(key, []byte(messageToEncrypt))
	logWithTimestamp("Encrypted cypher text", startTimeMilli)
	cv.So(err, cv.ShouldBeNil)

	decrypted, err := cipherutils.DecryptAES(key, cypherText)
	logWithTimestamp("Decrypted cypher text", startTimeMilli)
	cv.So(err, cv.ShouldBeNil)
	decryptedTxt := string(decrypted)

	cv.So(messageToEncrypt, cv.ShouldEqual, decryptedTxt)
}

func logWithTimestamp(msg string, startTimeMilli int64) {
	timeSinceStart := time.Now().UnixMilli() - startTimeMilli
	logger.Log.WithField("timeSinceStartMilli", timeSinceStart).Info(msg)
}
