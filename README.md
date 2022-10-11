# dmon - Google Dataflow Monitor

`dmon` is a CLI based application to monitor Google Dataflow jobs in a GCP project and send notifications if a job fails or times-out.

### Usage

To use `dmon` install it ...

```bash
<install here>
```

Then create a config file (structure is documented [below](#config)) and run `dmon` like this:

```bash
dmon -c path/to/config.yml
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

Controls the maximal timout in minutes for jobs. If a job runs for longer than the specified amount, a timeout notification will be triggered for that job. This does not apply to streaming jobs!

#### Expire Timeout Duration

```yaml
timeout:
  expire_timeout_duration: 10
```

...

### Project

#### ID

```yaml
project:
  id: my-google-project
```

...

#### Location

```yaml
project:
  location: europe-west4
```

...

### Slack

#### Token

```yaml
slack:
  token: my-secret-token
```

...

#### Channel

```yaml
slack:
  channel: my-google-project
```

...

#### Include Error Section

```yaml
slack:
  include_error_section: true
```

...

#### Include Dataflow Button

```yaml
slack:
  include_dataflow_button: true
```

...