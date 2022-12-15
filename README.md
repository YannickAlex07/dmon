# dmon - Google Dataflow Monitor

`dmon` is a CLI based application to monitor Google Dataflow jobs in a GCP project and send notifications if a job fails or times-out.

### Usage

To use `dmon` download the newest release from the `Release`-tab of this GitHub repository.

Then create a config file (structure is documented [below](#config)) and run `dmon` like this:

```bash
dmon -config path/to/config.yml
```

### Why was it developed?

Prior to `dmon` we used Google Cloud Monitoring and Alarming to receive Slack messages when our Dataflow jobs failed. The issue with this was mainly that the alerts that we received were often not very verbose and not helpful without digging deeper. To improve this, we wanted to develop an application that is more flexible in alerting and monitoring.

### How does it work?

`dmon` works by periodically listing all Dataflow jobs for a specific GCP project. It then checks the update time of the status for each job and when the update happend after the last time we ran the check, it will react to the status update by notifiying so-called `handlers` about the update. It will also calculate the total runtime of each job and will notify `handlers` if the job exceeds a configured timeout.

`Handlers` are structs that follow the `handler`-interface and can therefore receive updates about jobs from the monitor. Currently there is only a `SlackHandler` that is used to send Slack messages when jobs timeout or fail.

## Config

`dmon` offers quite a few configuration options that are available through the config file. This config file is a simple `yaml`-file that gets read during startup of the monitor. Here is an example of a full config:

```yml
request_interval: 2 # request every 2 minutes

logging:
  verbose: true # enable debug logging

timeout:
  max_timeout_duration: 10 # jobs that run longer than 10 min are considered timeouted
  expire_timeout_duration: 1440 # timeouts are cleared after 1440 minutes (24 hours)

project:
  id: my-google-project # GCP project id
  location: europe-west4 # GCP location that the Dataflow Jobs run in

slack:
  token: secret-slack-token # Token with permissions to post messages
  channel: my-error-channel # The channel that messages will be posted in
  include_error_section: true # If true, the latest error message of the job will be included
  include_dataflow_button: true # If true, a button that links to the Dataflow UI will be included
```

This section lists all the different options that are available in the config.

### Request Interval

```yaml
request_interval: 10
```

The request interval controls how much minutes pass between requesting jobs from the Dataflow API.

### Logging

#### Verbose

```yaml
logging: 
  verbose: true
```

If set to `true` messages with the level `DEBUG` are also logged.

### Timeout

#### Max Timeout Duration

```yaml
timeout:
  max_timeout_duration: 10
```

Controls the maximal timout in minutes for a job. If a job runs for longer than the specified amount, a timeout notification will be triggered for that job. This does not apply to streaming jobs!

#### Expire Timeout Duration

```yaml
timeout:
  expire_timeout_duration: 10
```

dmon keeps a list of all the job ID that triggered a timeout notification - this is done to not send out a timout notification for every check cycle. This setting controls after how many minutes a job is removed from this list - meaning a new timeout notification will be send out on the next check.

To always send a notification on each check cycle, set this lower than the `request_interval`.

### Project

#### ID

```yaml
project:
  id: my-google-project
```

The GCP project that dmon should monitor Dataflow jobs in.

#### Location

```yaml
project:
  location: europe-west4
```

The location that the Dataflow jobs run in. 

### Slack

#### Token

```yaml
slack:
  token: my-secret-token
```

The Slack token that dmon will use for authentication. Be aware that the Token needs permission to send messages - checkout the Slack documentation about this.

#### Channel

```yaml
slack:
  channel: my-slack-channel
```

The Slack channel that dmon will its messages into.

#### Include Error Section

```yaml
slack:
  include_error_section: true
```

If this is enabled, the last error message of the error will be attached to the
slack message.

#### Include Dataflow Button

```yaml
slack:
  include_dataflow_button: true
```

If this is enabled, a "Open in Dataflow"-button will be attached to the message. This button
will open the Dataflow UI of the job.