# gorush

A push notification server using [Gin](https://github.com/gin-gonic/gin) framework written in Go (Golang).

[![GoDoc](https://godoc.org/github.com/appleboy/gorush?status.svg)](https://godoc.org/github.com/appleboy/gorush) [![Build Status](https://travis-ci.org/appleboy/gorush.svg?branch=master)](https://travis-ci.org/appleboy/gorush) [![Coverage Status](https://coveralls.io/repos/github/appleboy/gorush/badge.svg?branch=master)](https://coveralls.io/github/appleboy/gorush?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/gorush)](https://goreportcard.com/report/github.com/appleboy/gorush) [![codebeat badge](https://codebeat.co/badges/0a4eff2d-c9ac-46ed-8fd7-b59942983390)](https://codebeat.co/projects/github-com-appleboy-gorush)

## Support Platform

* [APNS](https://developer.apple.com/library/ios/documentation/networkinginternet/conceptual/remotenotificationspg/Chapters/ApplePushService.html)
* [GCM](https://developer.android.com/google/gcm/index.html)

## Feature

* Support [Google Cloud Message](https://developers.google.com/cloud-messaging/) using [go-gcm](https://github.com/google/go-gcm) library for Android.
* Support [HTTP/2](https://http2.github.io/) Apple Push Notification Service using [apns2](https://github.com/sideshow/apns2) library.
* Support [YAML](https://github.com/go-yaml/yaml) configuration.
* Support command line to send single Android or iOS notification.
* Support Web API to send push notification.
* Support zero downtime restarts for go servers using [endless](https://github.com/fvbock/endless).
* Support [HTTP/2](https://http2.github.io/) or HTTP/1.1 protocol.
* Support notification queue and multiple workers.

See the [YAML config example](config/config.yml):

```yaml
core:
  port: "8088"
  worker_num: 8
  queue_num: 8192
  max_notification: 100
  mode: "release"
  ssl: false
  cert_path: "cert.pem"
  key_path: "key.pem"

api:
  push_uri: "/api/push"
  stat_go_uri: "/api/status"

android:
  enabled: true
  apikey: "YOUR_API_KEY"

ios:
  enabled: false
  pem_cert_path: "cert.pem"
  pem_key_path: "key.pem"
  production: false

log:
  format: "string" # string or json
  access_log: "stdout" # stdout: output to console, or define log path like "log/access_log"
  access_level: "debug"
  error_log: "stderr" # stderr: output to console, or define log path like "log/error_log"
  error_level: "error"
```

## Basic Usage

How to send push notification using `gorush` command? (Android or iOS)

### Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/gorush/releases).

### Send Android notification

Send single notification with the following command.

```bash
$ gorush -android -m="your message" -k="API Key" -t="Device token"
```

* `-m`: Notification message.
* `-k`: [Google cloud message](https://developers.google.com/cloud-messaging/) api key
* `-t`: Device token.

### Send iOS notification

Send single notification with the following command.

```bash
$ gorush -ios -m="your message" -i="API Key" -t="Device token"
```

* `-m`: Notification message.
* `-i`: Apple Push Notification Certificate path (`pem` file).
* `-t`: Device token.

The default endpoint is APNs development. Please add `-production` flag for APNs production push endpoint.

```bash
$ gorush -ios -m="your message" -i="API Key" -t="Device token" -production
```

## Run gorush web server

Please make sure your [config.yml](config/config.yml) exist. Default port is `8088`.

```bash
$ gorush -c config.yml
```

Test status of api server using [httpie](https://github.com/jkbrzt/httpie) tool:

```bash
$ http -v --verify=no --json GET https://localhost:8088/api/status
```

## Web API

gorush support the following API.

* **GET**  `/api/status` Golang cpu, memory, gc, etc information. Thanks for [golang-stats-api-handler](https://github.com/fukata/golang-stats-api-handler).
* **POST** `/api/push` push ios and android notifications.

### POST /api/push

Simple send iOS notification example, the `platform` value is `1`:

```json
{
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 1,
      "message": "Hello World iOS!"
    }
  ]
}
```

Simple send Android notification example, the `platform` value is `2`:

```json
{
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World Android!"
    }
  ]
}
```

Send multiple notifications as below:

```json
{
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 1,
      "message": "Hello World iOS!"
    },
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World Android!"
    },
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World!"
    },
    .....
  ]
}
```

See more example about [iOS](#ios-example) or [Android](#android-example).

### Request body

Request body must has a notifications array. The following is a parameter table for each notification.

|name|type|description|required|note|
|-------|-------|--------|--------|---------|
|tokens|string array|device tokens|o||
|platform|int|platform(iOS,Android)|o|1=iOS, 2=Android|
|message|string|message for notification|o||
|title|string|notification title|-||
|priority|string|Sets the priority of the message.|-|`normal` or `high`|
|content_available|bool|data messages wake the app by default.|-||
|sound|string|sound type|-||
|data|string array|extensible partition|-||
|api_key|string|Android api key|-|only Android|
|to|string|The value must be a registration token, notification key, or topic.|-|only Android|
|collapse_key|string|a key for collapsing notifications|-|only Android|
|delay_while_idle|bool|a flag for device idling|-|only Android|
|time_to_live|int|expiration of message kept on GCM storage|-|only Android|
|restricted_package_name|string|the package name of the application|-|only Android|
|dry_run|bool|allows developers to test a request without actually sending a message|-|only Android|
|notification|string array|payload of a GCM message|-|only Android. See the [detail](#android-notification-payload)|
|expiration|int|expiration for notification|-|only iOS|
|apns_id|string|A canonical UUID that identifies the notification|-|only iOS|
|topic|string|topic of the remote notification|-|only iOS|
|badge|int|badge count|-|only iOS|
|category|string|the UIMutableUserNotificationCategory object|-|only iOS|
|alert|string array|payload of a iOS message|-|only iOS. See the [detail](#ios-alert-payload)|

### iOS alert payload

|name|type|description|required|note|
|-------|-------|--------|--------|---------|
|action|string|The label of the action button. This one is required for Safari Push Notifications.|-||
|action-loc-key|string|If a string is specified, the system displays an alert that includes the Close and View buttons.|-||
|launch-image|string|The filename of an image file in the app bundle, with or without the filename extension.|-||
|loc-args|array of strings|Variable string values to appear in place of the format specifiers in loc-key.|-||
|loc-key|string|A key to an alert-message string in a Localizable.strings file for the current localization.|-||
|title-loc-args|array of strings|Variable string values to appear in place of the format specifiers in title-loc-key.|-||
|title-loc-key|string|The key to a title string in the Localizable.strings file for the current localization.|-||

See more detail about [APNs Remote Notification Payload](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/TheNotificationPayload.html#//apple_ref/doc/uid/TP40008194-CH107-SW1).

### Android notification payload

|name|type|description|required|note|
|-------|-------|--------|--------|---------|
|icon|string|Indicates notification icon.|-||
|tag|string|Indicates whether each notification message results in a new entry on the notification center on Android.|-||
|color|string|Indicates color of the icon, expressed in #rrggbb format|-||
|click_action|string|The action associated with a user click on the notification.|-||
|body_loc_key|string|Indicates the key to the body string for localization.|-||
|body_loc_args|string|Indicates the string value to replace format specifiers in body string for localization.|-||
|title_loc_key|string|Indicates the key to the title string for localization.|-||
|title_loc_args|string|Indicates the string value to replace format specifiers in title string for localization.|-||

See more detail about [GCM server reference](https://developers.google.com/cloud-messaging/http-server-ref#send-downstream).

### iOS Example

Send normal notification.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 1,
      "message": "Hello World iOS!",
      "title": "You got message"
    }
  ]
```

The app icon be badged with the number `9` and that a bundled alert sound be played when the notification is delivered.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 1,
      "message": "Hello World iOS!",
      "title": "You got message",
      "badge": 9,
      "sound": "bingbong.aiff"
    }
  ]
```

Add other fields which user defined via `data` field.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 1,
      "message": "Hello World iOS!",
      "title": "You got message",
      "data": {
        "key1": "welcome",
        "key2": 2
      }
    }
  ]
```

### Android Example

Send normal notification.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World Android!",
      "title": "You got message"
    }
  ]
```

Add `notification` payload.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World Android!",
      "title": "You got message",
      "notification" : {
        "icon": "myicon",
        "color": "#112244"
      }
    }
  ]
```

Add other fields which user defined via `data` field.

```json
  "notifications": [
    {
      "tokens": ["token_a", "token_b"],
      "platform": 2,
      "message": "Hello World Android!",
      "title": "You got message",
      "data": {
       "Nick" : "Mario",
       "body" : "great match!",
       "Room" : "PortugalVSDenmark"
      }
    }
  ]
```

### Response body

Error response message table:

|status code|message|
|-------|-------|
|400|Missing `notifications` field.|
|400|Notifications field is empty.|
|400|Number of notifications(50) over limit(10)|

Success response:

```json
{
  "success": "ok"
}
```

## Run gorush in Docker

Set up `gorush` in the cloud in under 5 minutes with zero knowledge of Golang or Linux shell using our [gorush Docker image](https://hub.docker.com/r/appleboy/gorush/).

```bash
$ docker pull appleboy/gorush
$ docker run -name gorush -p 80:8088 appleboy/gorush
```

Testing your gorush server.

```bash
$ http -v --verify=no --json GET http://your.docker.host/api/status
```

![statue screenshot](screenshot/status.png)

## License

Copyright 2016 Bo-Yi Wu [@appleboy](https://twitter.com/appleboy).

Licensed under the MIT License.
