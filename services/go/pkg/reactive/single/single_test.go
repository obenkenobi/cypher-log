package single_test

import (
	"context"
	"fmt"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSingleFromSupplierAsync(t *testing.T) {
	cv.Convey("When creating an observable with a supplier returns 1 that runs asynchronously", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
		defer cancel()
		oneSrc := single.FromSupplierAsync(func() (int, error) { return 1, nil })
		cv.Convey("The same source is then mapped to other Singles twoSrc and threeSrc,\n"+
			"incrementing the source value", func() {
			twoSrc := single.Map(oneSrc, func(v int) int { return v + 1 })
			threeSrc := single.Map(oneSrc, func(v int) int { return v + 2 })
			twoExpected, threeExpected := 2, 3
			tupleExpected := tuple.T2[int, int]{twoExpected, threeExpected}

			cv.Convey(
				fmt.Sprintf(
					"Expect zipping twoSrc and threeSrc successfully results in a tuple of (%v, %v)",
					twoExpected,
					threeExpected,
				),
				func() {
					zipped := single.Zip2(twoSrc, threeSrc)
					value, err := single.RetrieveValue(ctx, zipped)
					cv.So(err, cv.ShouldBeNil)
					cv.So(tupleExpected, cv.ShouldResemble, value)
				},
			)

			cv.Convey(fmt.Sprintf("Expect twoSrc returns %v when the value is retrieved", twoExpected), func() {
				value, err := single.RetrieveValue(ctx, twoSrc)
				cv.So(err, cv.ShouldBeNil)
				cv.So(twoExpected, cv.ShouldResemble, value)
			})

			cv.Convey(
				fmt.Sprintf("Expect threeSrc returns %v when the value is retrieved", threeExpected),
				func() {
					value, err := single.RetrieveValue(ctx, threeSrc)
					cv.So(err, cv.ShouldBeNil)
					cv.So(threeExpected, cv.ShouldResemble, value)
				},
			)
		})
	})
}
