# Slack plugin

`go get pfi/tanakad/sb-slack`

## Source

Require

* Outgoing webhooks integration https://api.slack.com/outgoing-webhooks
    * set server's URL to "URL(s)"

```sql
CREATE SOURCE s TYPE slack WITH port=8090;
```

The source will run a HTTP server addressed "http://localhost:8090" . This source is not supported "Responding" from Outgoing webhooks.

## Sink

Require

* Incoming webhooks integration https://api.slack.com/incoming-webhooks
    * the sink use "Webhook URL"

```sql
CREATE SINK t TYPE slack WITH
    hook="https://hooks.slack.com/services/XXX...";
```

The sink writes tuples to the "hook" address.

### Stream tuples

The sink is supported "text" field and "attachments" field.

```
CREATE STREAM t_msg AS SELECT RSTREAM
    s:text1 AS text,
    [{"text": s:text2, "color": "#36a64f"}] AS attachments
FROM some_stream AS s [RANGE 1 TUPLES];

INSERT INTO t FROM t_msg;
```

"s:text1" is posted as "text" field, and "s:text2" is posted as "attachments" field. Below JSON sample represents "t_msg" stream.

```json
{
  "text": "<s:text1 contents>",
  "attachments": [
    {
      "text": "<s:text2 contents>",
      "color": "#36a64f"
    }
  ]
}
```

More detail, see https://api.slack.com/docs/message-attachments
