package slack

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"net/http"
	"net/url"
	"testing"
)

func TestNewSource(t *testing.T) {
	Convey("Given a test context", t, func() {
		cc := &core.ContextConfig{}
		ctx := core.NewContext(cc)
		ioParams := &bql.IOParams{}
		Convey("When create a source with default value", func() {
			params := data.Map{}
			s, err := NewSource(ctx, ioParams, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})
			Convey("And when start to generate stream", func() {
				l, ok := s.(*listener)
				So(ok, ShouldBeTrue)
				var lt core.Tuple
				w := core.WriterFunc(func(ctx *core.Context, t *core.Tuple) error {
					lt = *t
					return nil
				})
				go func() {
					l.GenerateStream(ctx, w)
				}()
				Convey("Then the source listens on default address", func() {
					value := url.Values{
						"token":        {"token"},
						"team_id":      {"team_id"},
						"channel_id":   {"channel_id"},
						"channel_name": {"channel_name"},
						"timestamp":    {"timestamp"},
						"user_id":      {"user_id"},
						"user_name":    {"user_name"},
						"text":         {"text"},
						"trigger_word": {"trigger_word"},
					}
					res, err := http.PostForm("http://localhost:15619", value)
					So(res.StatusCode, ShouldEqual, 200)
					So(lt.Data["token"], ShouldEqual, "token")
					So(lt.Data["team_id"], ShouldEqual, "team_id")
					So(lt.Data["channel_id"], ShouldEqual, "channel_id")
					So(lt.Data["channel_name"], ShouldEqual, "channel_name")
					So(lt.Data["timestamp"], ShouldEqual, "timestamp")
					So(lt.Data["user_id"], ShouldEqual, "user_id")
					So(lt.Data["user_name"], ShouldEqual, "user_name")
					So(lt.Data["text"], ShouldEqual, "text")
					So(lt.Data["trigger_word"], ShouldEqual, "trigger_word")
					So(err, ShouldBeNil)
				})
			})
		})
		Convey("When create a source with customized value", func() {
			params := data.Map{
				"api_header": data.String("/v1"),
				"port":       data.Int(15620),
			}
			s, err := NewSource(ctx, ioParams, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})
			Convey("And when start to generate stream", func() {
				l, ok := s.(*listener)
				So(ok, ShouldBeTrue)
				var lt core.Tuple
				w := core.WriterFunc(func(ctx *core.Context, t *core.Tuple) error {
					lt = *t
					return nil
				})
				go func() {
					l.GenerateStream(ctx, w)
				}()
				Convey("Then the source listens on default address", func() {
					value := url.Values{
						"token":        {"token"},
						"team_id":      {"team_id"},
						"channel_id":   {"channel_id"},
						"channel_name": {"channel_name"},
						"timestamp":    {"timestamp"},
						"user_id":      {"user_id"},
						"user_name":    {"user_name"},
						"text":         {"text"},
						"trigger_word": {"trigger_word"},
					}
					res, err := http.PostForm("http://localhost:15620/v1", value)
					So(res.StatusCode, ShouldEqual, 200)
					So(lt.Data["token"], ShouldEqual, "token")
					So(lt.Data["team_id"], ShouldEqual, "team_id")
					So(lt.Data["channel_id"], ShouldEqual, "channel_id")
					So(lt.Data["channel_name"], ShouldEqual, "channel_name")
					So(lt.Data["timestamp"], ShouldEqual, "timestamp")
					So(lt.Data["user_id"], ShouldEqual, "user_id")
					So(lt.Data["user_name"], ShouldEqual, "user_name")
					So(lt.Data["text"], ShouldEqual, "text")
					So(lt.Data["trigger_word"], ShouldEqual, "trigger_word")
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
