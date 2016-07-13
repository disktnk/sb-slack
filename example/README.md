# SensorBee slack plugin example

This is a example of SensorBee slack plugin and mstranslator plugin.

## Flow

1. [Slack] Outgoing webhooks integration posts messages to SensorBee server.
1. [Sensorbee] Emit posted messages from slack.
1. [SensorBee] Calls MS TRANSLATOR API and posts messages.
1. [MS TRANSLATOR] Return translated messages.
1. [SensorBee] Post messages to slack via incoming webhooks.
1. [Slack] Incoming webhooks integration pushes message to specified channel (or DM).

## Setup

### Require

* SensorBee >= 0.4.2
* slack plugin `go get github.com/disktnk/sb-slack`
* mstranslator plugin `go get github.com/disktnk/sb-mstranslator`
* slack incoming webhooks integrator
    * use "webhooks address"
* slack outgoing webhooks integrator
    * set server address to "URL(s)"
* microsoft developer account
    * setup translator API
    * use "client\_id" and "client\_secret"

see: slack plugin README and mstranslator plugin README

### Build

```bash
$ build_sensorbee
```

### Make BQL

```bash
$ cp test.bql_sample test.bql
```

Edit test.bql

* `<client id>`: "Client ID" of microsoft developer account
* `<client secret>`: "Client Secret" of microsoft developer account
* `<hook address>`: "webhooks address" of slack incoming webhooks integrator

## Run

```bash
$ ./sensorbee run -c conf.yaml
```
