package option_test

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNilMaybe(t *testing.T) {
	cv.Convey("When passing nil to option.Perhaps", t, func() {
		type dummyType struct {
			val string
		}
		nilMaybe := option.Perhaps[*dummyType](nil)
		cv.Convey("Expect calling IsEmpty to be true and IsPresent to be false", func() {
			cv.So(nilMaybe.IsEmpty(), cv.ShouldBeTrue)
			cv.So(nilMaybe.IsPresent(), cv.ShouldBeFalse)
		})
		cv.Convey("Expect a nil Maybe to be the same as calling None()", func() {
			cv.So(nilMaybe, cv.ShouldResemble, option.None[*dummyType]())
		})

		cv.Convey(
			"Expect a flatmap and a map returning a same type Maybe will return a resembling empty Maybe",
			func() {
				mappedMaybe := option.FlatMap(nilMaybe, func(v1 *dummyType) option.Maybe[*dummyType] {
					return option.Map(nilMaybe, func(v2 *dummyType) *dummyType {
						return &dummyType{val: v1.val + v2.val + "123"}
					})
				})
				cv.So(nilMaybe, cv.ShouldResemble, mappedMaybe)
			},
		)
		cv.Convey(
			"Expect a flatmap and a map returning a string type Maybe will return a not equal empty Maybe",
			func() {
				mappedMaybeString := option.Map(nilMaybe, func(v1 *dummyType) string {
					return v1.val
				})
				cv.So(nilMaybe, cv.ShouldNotEqual, mappedMaybeString)
			},
		)
		cv.Convey(
			"Expect a filter on an empty maybe returns a same type NoneQuery",
			func() {
				maybeFiltered := nilMaybe.Filter(func(d *dummyType) bool { return false })
				cv.So(nilMaybe, cv.ShouldResemble, maybeFiltered)
			},
		)
		cv.Convey("Expect calling orElse returns the value specified in the expected struct", func() {
			expected := dummyType{"Hello World"}
			other := expected
			result := nilMaybe.OrElse(&other)
			cv.So(&other, cv.ShouldResemble, result)
			cv.So(other, cv.ShouldResemble, *result)
		})
		cv.Convey("Expect calling orElseGet returns the value specified in the expected struct", func() {
			expected := dummyType{"Hello World"}
			result := nilMaybe.OrElseGet(func() *dummyType {
				other := expected
				return &other
			})
			cv.So(&expected, cv.ShouldResemble, result)
			cv.So(expected, cv.ShouldResemble, *result)
		})
	})
}
func TestHasValueMaybe(t *testing.T) {
	cv.Convey("When passing string to option.Perhaps", t, func() {
		defaultValue := "Hello"
		countStr := "1234"
		maybe := option.Perhaps(defaultValue)
		cv.Convey("Expect when held value is to be retrieved called", func() {
			cv.Convey(
				"When orElseGet is called on a new valuable, the result equals the value in the original Maybe",
				func() {
					cv.So(defaultValue, cv.ShouldEqual, maybe.OrElseGet(func() string { return countStr }))
				},
			)
			cv.Convey(
				"When orElse is called on a new valuable, the result equals the value in the original Maybe",
				func() {
					cv.So(defaultValue, cv.ShouldEqual, maybe.OrElse(countStr))
				},
			)
		})

		cv.Convey("Expect when calling flatMap can return a Maybe with a new value", func() {
			expectedStrTransform := "Hello1234"
			cv.Convey(fmt.Sprintf("Calling OrElse returns a string of value: %v", expectedStrTransform), func() {
				mappedMaybe := option.FlatMap(maybe, func(v1 string) option.Maybe[string] {
					return option.Perhaps(v1 + countStr)
				})
				cv.So(expectedStrTransform, cv.ShouldEqual, mappedMaybe.OrElse("."))
			})
			cv.Convey(fmt.Sprintf("Calling OrElseGet returns a string of value: %v", expectedStrTransform), func() {
				mappedMaybe := option.Map(maybe, func(v1 string) string {
					return v1 + countStr
				})
				cv.So(expectedStrTransform, cv.ShouldEqual, mappedMaybe.OrElse("."))
			})
		})
		cv.Convey("Expect calling IsEmpty to be false and IsPresent to be true ", func() {
			cv.So(maybe.IsEmpty(), cv.ShouldBeFalse)
			cv.So(maybe.IsPresent(), cv.ShouldBeTrue)
		})

		cv.Convey("Expect running filter with a failing predicate will return an empty maybe", func() {
			filered := maybe.Filter(func(v string) bool {
				return v != defaultValue
			})
			cv.So(filered.IsPresent(), cv.ShouldBeFalse)
			cv.So(filered.IsEmpty(), cv.ShouldBeTrue)
		})

		cv.Convey("Expect running filter with a passing predicate will return a non-empty maybe", func() {
			filered := maybe.Filter(func(v string) bool {
				return v == defaultValue
			})
			cv.So(filered.IsPresent(), cv.ShouldBeTrue)
			cv.So(filered.IsEmpty(), cv.ShouldBeFalse)
		})

	})
}
