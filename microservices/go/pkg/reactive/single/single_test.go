package single_test

import (
	"context"
	"fmt"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSingleFromSupplierAsync(t *testing.T) {
	cv.Convey("When creating an observable with a supplier returns 1 that runs asynchronously", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
		defer cancel()
		oneSrc := single.FromSupplier(func() (int, error) { return 1, nil }).ScheduleEagerAsync(ctx)
		cv.Convey("The same source is then mapped to other Singles twoSrc and threeSrc,\n"+
			"incrementing the source value", func() {
			twoSrc := single.Map(oneSrc, func(v int) int { return v + 1 })
			threeSrc := single.Map(oneSrc, func(v int) int { return v + 2 })
			twoExpected, threeExpected := 2, 3
			tupleExpected := stream.Tuple2[int, int]{V1: twoExpected, V2: threeExpected}

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
