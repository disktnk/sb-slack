package plugin

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	slack "pfi/tanakad/sb-slack"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("slack", bql.SourceCreatorFunc(
		slack.NewSource))
	bql.MustRegisterGlobalSinkCreator("slack", bql.SinkCreatorFunc(
		slack.NewSink))
}
