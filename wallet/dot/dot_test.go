package dot

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestQueryStorage(t *testing.T) {
	convey.Convey("TestQueryStorage", t, func() {
		// resp, err := QueryStorage("0x2ea4353f9c2aa4da42c318efbf01e14bd768aa825b75f6cd7b90c184c38b92ab")
		// convey.So(resp, convey.ShouldNotBeZeroValue)
		// convey.So(err, convey.ShouldBeNil)
	})
}

func TestGetAccountNonce(t *testing.T) {
	convey.Convey("TestGetAccountNonce", t, func() {
		resp, err := GetAccountNonce([]string{"16ZL8yLyXv3V3L3z9ofR1ovFLziyXaN1DPq4yffMAZ9czzBD"})
		convey.So(resp, convey.ShouldNotBeZeroValue)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGetBlock(t *testing.T) {
	convey.Convey("TestGetBlock", t, func() {
		resp, err := GetBlock()
		convey.So(resp, convey.ShouldNotBeZeroValue)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGetRuntimeVersion(t *testing.T) {
	convey.Convey("TestGetRuntimeVersion", t, func() {
		resp, err := GetRuntimeVersion()
		convey.So(resp, convey.ShouldNotBeZeroValue)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSubmitExtrinsic(t *testing.T) {
	convey.Convey("TestSubmitExtrinsic", t, func() {
		// resp, err := SubmitExtrinsic([]string{"0x00000000"})
		// convey.So(resp, convey.ShouldNotBeZeroValue)
		// convey.So(err, convey.ShouldBeNil)
	})
}

func TestTx(t *testing.T) {
	convey.Convey("TestTx", t, func() {
		// resp, err := GetTx()
		// convey.So(resp, convey.ShouldNotBeZeroValue)
		// convey.So(err, convey.ShouldBeNil)
	})
}
