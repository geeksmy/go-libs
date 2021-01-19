package grpcclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ParseToHostPort(t *testing.T) {
	Convey("Test_ParseToHostPort", t, func() {
		Convey("normal", func() {
			host, port, err := ParseToHostPort("127.0.0.1:8081")
			So(err, ShouldBeNil)
			So(host, ShouldEqual, "127.0.0.1")
			So(port, ShouldEqual, 8081)
		})

		Convey("normal has Scheme", func() {
			host, port, err := ParseToHostPort("grpc://127.0.0.1:8081")
			So(err, ShouldBeNil)
			So(host, ShouldEqual, "127.0.0.1")
			So(port, ShouldEqual, 8081)
		})

		Convey("no port, parse fail", func() {
			host, _, err := ParseToHostPort("127.0.0.1:")
			So(err, ShouldBeError)
			So(host, ShouldEqual, "127.0.0.1")
		})

		Convey("no port 2, parse fail", func() {
			_, _, err := ParseToHostPort("127.0.0.1")
			So(err, ShouldBeError)
		})
	})
}
