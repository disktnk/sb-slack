package slack

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
				So(wb.iconEmoji, ShouldEqual, "icon_emoji")
			})
		})
	})
}

func TestWrite(t *testing.T) {
	Convey("Given a HTTP server for test and default sink", t, func() {
		actualText := "error"
		actualChannel := "error"
		actualUsername := "error"
		actualIconURL := "error"
		actualIconEmoji := "error"
		var actualAttachments data.Array
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				var p payload
				json.Unmarshal(b, &p)
				actualText = p.Text
				actualChannel = p.Channel
				actualUsername = p.Username
				actualIconURL = p.IconURL
				actualIconEmoji = p.IconEmoji
				actualAttachments = p.Attachments
			},
		)
		ts := httptest.NewServer(mux)
		Reset(func() {
			ts.Close()
		})

		cc := &core.ContextConfig{}
		ctx := core.NewContext(cc)
		ioParams := &bql.IOParams{}
		params := data.Map{
			"hook":               data.String(ts.URL),
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

		Convey("When write a empty tuple", func() {
			d := data.Map{}
			t := core.NewTuple(d)
			err := s.Write(ctx, t)
			So(err, ShouldBeNil)
			Convey("Then actual tuple should be overwritten in default values", func() {
				So(actualText, ShouldEqual, "")
				So(actualChannel, ShouldEqual, "channel")
				So(actualUsername, ShouldEqual, "username")
				So(actualIconURL, ShouldEqual, "icon_url")
				So(actualIconEmoji, ShouldEqual, "icon_emoji")
				So(actualAttachments, ShouldResemble, data.Array(nil))
			})
		})

		Convey("When write a tuple with text and bot information", func() {
			d := data.Map{
				"text":       data.String("homhom"),
				"channel":    data.String("_channel"),
				"username":   data.String("_username"),
				"icon_url":   data.String("_icon_url"),
				"icon_emoji": data.String("_icon_emoji"),
			}
			t := core.NewTuple(d)
			err := s.Write(ctx, t)
			So(err, ShouldBeNil)
			Convey("Then actual tuple should be overwritten in default values", func() {
				So(actualText, ShouldEqual, "homhom")
				So(actualChannel, ShouldEqual, "_channel")
				So(actualUsername, ShouldEqual, "_username")
				So(actualIconURL, ShouldEqual, "_icon_url")
				So(actualIconEmoji, ShouldEqual, "_icon_emoji")
				So(actualAttachments, ShouldResemble, data.Array(nil))
			})
		})

		Convey("When write a tuple with attachment", func() {
			att := data.Array{
				data.Map{
					"pretext": data.String("pretext"),
					"text":    data.String("sub_text"),
				},
			}
			d := data.Map{
				"attachments": att,
			}
			t := core.NewTuple(d)
			err := s.Write(ctx, t)
			So(err, ShouldBeNil)
			Convey("Then actual tuple should be overwritten in default values", func() {
				So(actualText, ShouldEqual, "")
				So(actualChannel, ShouldEqual, "channel")
				So(actualUsername, ShouldEqual, "username")
				So(actualIconURL, ShouldEqual, "icon_url")
				So(actualIconEmoji, ShouldEqual, "icon_emoji")
				So(actualAttachments, ShouldResemble, att)
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
