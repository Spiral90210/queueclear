# queueclear

Simple project to read an SQS queue and clear it, writing each message to stdout (disable with `-quiet`)

To run it:

```shell
go run ./ -queue <queuename> [-quiet]
```

This uses your current AWS_PROFILE, or other env vars for any and all AWS configuration. You can use this to execute against a local environment like:

```ini
[profile myprofile]
region = us-west-2
services = my-services

[services my-services]
sqs =
  endpoint_url = http://localhost:9324
```

```shell
AWS_PROFILE=myprofile go run ./ -queue localqueuename
```
