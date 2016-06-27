package slack

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestNewSink(t *testing.T) {
	Convey("Given a test context", t, func() {
		cc := &core.ContextConfig{}
		ctx := core.NewContext(cc)
		ioParams := &bql.IOParams{}
		Convey("When create a sink with only required parameters", func() {
			params := data.Map{
				"hook": data.String("hoge"),
			}
			s, err := NewSink(ctx, ioParams, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Close(ctx)
			})
			Convey("Then the sink should set up with default values", func() {
				wb, ok := s.(*webHook)
				So(ok, ShouldBeTrue)
				So(wb.hookURL, ShouldEqual, "hoge")
				So(wb.channel, ShouldEqual, "")
				So(wb.username, ShouldEqual, "")
				So(wb.iconURL, ShouldEqual, "")
				So(wb.iconEmoji, ShouldEqual, "")
			})
		})
		Convey("When create a sink with customized parameters", func() {
			params := data.Map{
				"hook":               data.String("hoge"),
				"default_channel":    data.String("channel"),
				"default_username":   data.String("username"),
				"default_icon_url":   data.String("icon_url"),
				"default_icon_emoji": data.String("icon_emoji"),
			}
			s, err := NewSink(ctx, ioParams, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Close(ctx)
			})
			Convey("Then the sink should set up with the parameters", func() {
				wb, ok := s.(*webHook)
				So(ok, ShouldBeTrue)
				So(wb.hookURL, ShouldEqual, "hoge")
				So(wb.channel, ShouldEqual, "channel")
				So(wb.username, ShouldEqual, "username")
				So(wb.iconURL, ShouldEqual, "icon_url")
				So(wb.iconEmoji, ShouldEqual, "")
			})
		})
		Convey("When create a sink with 'icon_emoji'", func() {
			params := data.Map{
				"hook":               data.String("hoge"),
				"default_icon_emoji": data.String("icon_emoji"),
			}
			s, err := NewSink(ctx, ioParams, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Close(ctx)
			})
			Convey("Then the sink should set up with 'icon_emoji'", func() {
				wb, ok := s.(*webHook)
				So(ok, ShouldBeTrue)
				So(wb.hookURL, ShouldEqual, "hoge")
				So(wb.channel, ShouldEqual, "")
				So(wb.username, ShouldEqual, "")
				So(wb.iconURL, ShouldEqual, "")
				So(wb.iconEmoji, ShouldEqual, "icon_emoji")
			})
		})
	})
}

func TestNewSinkWithError(t *testing.T) {
	Convey("Given a test context", t, func() {
		cc := &core.ContextConfig{}
		ctx := core.NewContext(cc)
		ioParams := &bql.IOParams{}
		Convey("When create a sink without web hook address", func() {
			params := data.Map{}
			Convey("Then an error should be occurred", func() {
				_, err := NewSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a sink with error channel name", func() {
			params := data.Map{
				"hook":            data.String("hoge"),
				"default_channel": data.False,
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a sink with error user name", func() {
			params := data.Map{
				"hook":             data.String("hoge"),
				"default_username": data.Int(55),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a sink with error icon URL", func() {
			params := data.Map{
				"hook":             data.String("hoge"),
				"default_icon_url": data.Float(0.1),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a sink with error icon emoji", func() {
			params := data.Map{
				"hook":               data.String("hoge"),
				"default_icon_emoji": data.Blob([]byte("hoge")),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
