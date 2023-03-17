# dmon - Google Dataflow Monitor

`dmon` is a CLI based application to monitor Google Dataflow jobs in a GCP project and send notifications if a job fails or times-out.

### Usage

To use `dmon` download the newest release from the `Release`-tab of this GitHub repository.

Then create a config file (structure is documented [here](./docs/config.md)) and run `dmon` like this:

```bash
dmon -c path/to/config.yml
```

### Why was it developed?

Prior to `dmon` we used Google Cloud Monitoring and Alarming to receive Slack messages when our Dataflow jobs failed or nothing if they timed out. The issue with this was mainly that the alerts that we received were often not very verbose and not helpful without digging deeper. To improve this, we wanted to develop an application that is more flexible in alerting and monitoring.

### How does it work?

`dmon` works by periodically listing all Dataflow jobs for a specific GCP project. It then checks the update time of the status for each job and when the update happend after the last time we ran the check, it will react to the status update by notifiying so-called `handlers` about the update. It will also calculate the total runtime of each job and will notify `handlers` if the job exceeds a configured timeout.

`Handlers` are structs that follow the `handler`-interface and can therefore receive updates about jobs from the monitor. Currently there is only a `SlackHandler` that is used to send Slack messages when jobs timeout or fail. You can implement your own handler if you want to.

### Further Documentation

To find more information on how to use dmon, check the following documents:

* [Config](./docs/config.md)
* [Release](./docs/release.md)