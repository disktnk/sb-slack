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
	attachmentsPath      = data.MustCompilePath("attachments")
	iconURLPath          = data.MustCompilePath("icon_url")
	iconEmojiPath        = data.MustCompilePath("icon_emoji")
)

// NewSink returns a sink to POST slack messages.
//
// "hook": required value, use incoming webhooks integration address.
//
// "default_channel": option value, default: "". If empty, tuples are emit
// to the address set by incoming webhooks. Channel name can be overwritten
// by tuples' "channel" key.
//
// "default_username": option value, default: "". If empty, displayed user
// name is set by incoming webhooks. User name can be overwritten by tuples'
// "username" key.
//
// "default_icon_url": option value, default: "". If empty, displayed icon
// is set by incoming webhooks. Icon URL can be overwritten by tuples'
// "icon_url" key.
//
// "default_icon_emoji": option value, default: "". Spec is same as
// "default_icon_url". Overrides "icon_url".
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
	Channel     string     `json:"channel,omitempty"`
	Username    string     `json:"username,omitempty"`
	Text        string     `json:"text"`
	Attachments data.Array `json:"attachments,omitempty"`
	IconURL     string     `json:"icon_url,omitempty"`
	IconEmoji   string     `json:"icon_emoji,omitempty"`
}

func (h *webHook) Write(ctx *core.Context, t *core.Tuple) error {
	text := ""
	if txt, err := t.Data.Get(textPath); err == nil {
		if text, err = data.AsString(txt); err != nil {
			return err
		}
	}
	p := payload{
		Text: text,
	}
	attachments := data.Array{}
	if att, err := t.Data.Get(attachmentsPath); err == nil {
		if attachments, err = data.AsArray(att); err != nil {
			return err
		}
		p.Attachments = attachments
	}

	if ch, err := t.Data.Get(channelPath); err != nil {
		if h.channel != "" {
			p.Channel = h.channel
		}
	} else {
		p.Channel, err = data.AsString(ch)
		if err != nil {
			return err
		}
	}
	if un, err := t.Data.Get(usernamePath); err != nil {
		if h.username != "" {
			p.Username = h.username
		}
	} else {
		p.Username, err = data.AsString(un)
		if err != nil {
			return err
		}
	}

	if iu, err := t.Data.Get(iconURLPath); err != nil {
		if h.iconURL != "" {
			p.IconURL = h.iconURL
		}
	} else {
		p.IconURL, err = data.AsString(iu)
		if err != nil {
			return err
		}
	}
	if ie, err := t.Data.Get(iconEmojiPath); err != nil {
		if h.iconEmoji != "" {
			p.IconEmoji = h.iconEmoji
		}
	} else {
		p.IconEmoji, err = data.AsString(ie)
		if err != nil {
			return err
		}
	}

	return post(h.hookURL, p)
}

func post(url string, p payload) error {
	jsonByte, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
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
