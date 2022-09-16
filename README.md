# dmon - Google Dataflow Monitor

`dmon` is a CLI based application to monitor Google Dataflow jobs in a GCP project and send notifications if a job fails or times-out.

### Why was it developed?

Prior to `dmon` we used Google Cloud Monitoring and Alarming to receive Slack messages when our Dataflow jobs failed. The issue with this was mainly that the alerts that we received were often not very verbose and not helpful without digging deeper. To improve this, we wanted to develop an application that is more flexible in alerting and monitoring.

### How does it work?

`dmon` works by periodically listing all Dataflow jobs for a specific GCP project. It then checks the update time of the status for each job and when the update happend after the last time we ran the check, it will react to the status update by notifiying so-called `handlers` about the update. It will also calculate the total runtime of each job and will notify `handlers` if the job exceeds a configured timeout.

`Handlers` are structs that follow the `handler`-interface and can therefore receive updates about jobs from the monitor. Currently there is only a `SlackHandler` that is used to send Slack messages when jobs timeout or fail.

## Config

`dmon` offers quite a few configuration options that are available through the config file. This config file is a simple `yaml`-file that gets read during startup of the monitor.

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

