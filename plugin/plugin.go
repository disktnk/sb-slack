package plugin

import (
	slack "github.com/disktnk/sb-slack"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("slack", bql.SourceCreatorFunc(
		slack.NewSource))
	bql.MustRegisterGlobalSinkCreator("slack", bql.SinkCreatorFunc(
		slack.NewSink))
}
