package reactive

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	cv "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestObservableFromChannels(t *testing.T) {
	cv.Convey("When creating a channel for an observable that sends 3 numbers", t, func() {
		intCh := make(chan int)
		go func() {
			defer close(intCh)
			for i := range [3]int{} {
				intCh <- i + 1
			}
		}()
		ctx := context.Background()
		intObs := stream.FromChannel(intCh)
		resCh, err := stream.ToChannels(ctx, intObs)
		var intList []int
		for i := range resCh {
			intList = append(intList, i)
			if i == 3 {
				break
			}
		}
		cv.Convey("3 numbers are listened to for the observable", func() {
			cv.So(len(intList), cv.ShouldEqual, 3)
			cv.So(intList, cv.ShouldResemble, []int{1, 2, 3})
			cv.So(<-err, cv.ShouldBeNil)
		})
	})
}
