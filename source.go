package slack

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	apiHeaderPath = data.MustCompilePath("api_header")
	portPath      = data.MustCompilePath("port")
)

// NewSource returns a HTTP listener, which has started to listening when
// create.
//
// "api_header": option value, default: "/"
//
// "port": option value, default: 15619
func NewSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Source, error) {
	apiHeader := "/"
	if ah, err := params.Get(apiHeaderPath); err == nil {
		if apiHeader, err = data.AsString(ah); err != nil {
			return nil, err
		}
	}
	port := ":15619"
	if po, err := params.Get(portPath); err == nil {
		if portInt, err := data.AsInt(po); err == nil {
			port = fmt.Sprintf("%d", portInt)
		} else if port, err = data.AsString(po); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(port, ":") {
			port = fmt.Sprintf(":%s", port)
		}
	}

	l := listener{
		msgCh:  make(chan message),
		stopCh: make(chan struct{}),
		errCh:  make(chan error),
	}

	nl, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	l.netListner = nl
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc(apiHeader, l.bind)
		l.errCh <- http.Serve(nl, mux)
	}()
	ctx.Log().Infof("listening server has started, port%v", port)

	return &l, nil
}

type message struct {
	token       string
	teamID      string
	channelID   string
	channelName string
	timestamp   string
	userID      string
	userName    string
	text        string
	triggerWord string
}

type listener struct {
	msgCh      chan message
	stopCh     chan struct{}
	errCh      chan error
	netListner net.Listener
}

func (l *listener) bind(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	teamID := r.FormValue("team_id")
	channelID := r.FormValue("channel_id")
	channelName := r.FormValue("channel_name")
	timestamp := r.FormValue("timestamp")
	userID := r.FormValue("user_id")
	userName := r.FormValue("user_name")
	text := r.FormValue("text")
	triggerWord := r.FormValue("trigger_word")

	msg := message{
		token:       token,
		teamID:      teamID,
		channelID:   channelID,
		channelName: channelName,
		timestamp:   timestamp,
		userID:      userID,
		userName:    userName,
		text:        text,
		triggerWord: triggerWord,
	}

	l.msgCh <- msg
}

func (l *listener) GenerateStream(ctx *core.Context, w core.Writer) error {
	for {
		select {
		case msg := <-l.msgCh:
			m := data.Map{
				"token":        data.String(msg.token),
				"team_id":      data.String(msg.teamID),
				"channel_id":   data.String(msg.channelID),
				"channel_name": data.String(msg.channelName),
				"timestamp":    data.String(msg.timestamp),
				"user_id":      data.String(msg.userID),
				"user_name":    data.String(msg.userName),
				"text":         data.String(msg.text),
				"trigger_word": data.String(msg.triggerWord),
			}
			now := time.Now()
			t := &core.Tuple{
				Data:          m,
				Timestamp:     now, // TODO: should use message timestamp
				ProcTimestamp: now,
				Trace:         []core.TraceEvent{},
			}
			if err := w.Write(ctx, t); err != nil {
				return err
			}
		case err := <-l.errCh:
			return err
		case <-l.stopCh:
			return core.ErrSourceStopped
		}
	}
	return nil
}

func (l *listener) Stop(ctx *core.Context) error {
	close(l.stopCh)
	return l.netListner.Close()
}
