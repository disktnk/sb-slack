package slack

import (
	"bytes"
	"encoding/json"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"net/http"
	"time"
)

var (
	hookPath             = data.MustCompilePath("hook")
	defaultChannelPath   = data.MustCompilePath("default_channel")
	defaultUsernamePath  = data.MustCompilePath("default_username")
	defaultIconURLPath   = data.MustCompilePath("default_icon_url")
	defaultIconEmojiPath = data.MustCompilePath("default_icon_emoji")
	channelPath          = data.MustCompilePath("channel")
	usernamePath         = data.MustCompilePath("username")
	textPath             = data.MustCompilePath("text")
	iconURLPath          = data.MustCompilePath("icon_url")
	iconEmojiPath        = data.MustCompilePath("icon_emoji")
)

func NewSink(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Sink, error) {
	addr, err := params.Get(hookPath)
	if err != nil {
		return nil, err
	}
	hookURL, err := data.AsString(addr)
	if err != nil {
		return nil, err
	}
	channel := ""
	if dc, err := params.Get(defaultChannelPath); err == nil {
		if channel, err = data.AsString(dc); err != nil {
			return nil, err
		}
	}
	username := ""
	if dun, err := params.Get(defaultUsernamePath); err == nil {
		if username, err = data.AsString(dun); err != nil {
			return nil, err
		}
	}
	iconURL := ""
	if diu, err := params.Get(defaultIconURLPath); err == nil {
		if iconURL, err = data.AsString(diu); err != nil {
			return nil, err
		}
	}
	iconEmoji := ""
	if diem, err := params.Get(defaultIconEmojiPath); err == nil {
		if iconEmoji, err = data.AsString(diem); err != nil {
			return nil, err
		}
	}
	if iconURL != "" && iconEmoji != "" {
		ctx.Log().Warnf(
			"cannot set both 'icon_url' and 'icon_emoji', '%s' is used as priority",
			iconURL)
		iconEmoji = ""
	}

	return &webHook{
		hookURL:   hookURL,
		channel:   channel,
		username:  username,
		iconURL:   iconURL,
		iconEmoji: iconEmoji,
	}, nil
}

type webHook struct {
	hookURL   string
	channel   string
	username  string
	iconURL   string
	iconEmoji string
}

type payload struct {
	channel   string `json:"channel"`
	username  string `json:"username"`
	text      string `json:"text"`
	iconURL   string `json:"icon_url,omitempty"`
	iconEmoji string `json:"icon_emoji,omitempty"`
}

func (h *webHook) Write(ctx *core.Context, t *core.Tuple) error {
	channel := h.channel
	if ch, err := t.Data.Get(channelPath); err != nil {
		if h.channel == "" {
			return err
		}
	} else {
		channel, err = data.AsString(ch)
		if err != nil {
			return err
		}
	}
	username := h.username
	if un, err := t.Data.Get(usernamePath); err != nil {
		if h.username == "" {
			return err
		}
	} else {
		username, err = data.AsString(un)
		if err != nil {
			return err
		}
	}
	text := ""
	if txt, err := t.Data.Get(textPath); err != nil {
		return err
	} else if text, err = data.AsString(txt); err != nil {
		return err
	}

	p := payload{
		channel:  channel,
		username: username,
		text:     text,
	}

	if iu, err := t.Data.Get(iconURLPath); err != nil {
		if h.iconURL != "" {
			p.iconURL = h.iconURL
		}
	} else {
		h.iconURL, err = data.AsString(iu)
		if err != nil {
			return err
		}
	}
	if ie, err := t.Data.Get(iconEmojiPath); err != nil {
		if h.iconEmoji != "" {
			p.iconEmoji = h.iconEmoji
		}
	} else {
		if h.iconURL != "" {
			ctx.Log().Warnf(
				"cannot set both 'icon_url' and 'icon_emoji', '%s' is used as priority",
				h.iconURL) // TODO it is possible to occur many warning log...
		} else {
			h.iconEmoji, err = data.AsString(ie)
			if err != nil {
				return err
			}
		}
	}

	return post(h.hookURL, p)
}

func post(url string, p payload) error {
	jsonByte, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonByte))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	res, err := cli.Do(req)
	defer res.Body.Close()

	return err
}

func (h *webHook) Close(ctx *core.Context) error {
	return nil
}
