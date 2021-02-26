# aws-ip

Sync/set AWS security group entry (by description) with current user's local IP address. I.e. if the access to AWS resources is limited by the IP address, and the user/machine doesn't have 

## Syntax
Possible syntax considered

```
aws-ip --group-name <security group name> --port <port> --location <free text/code of the current user's location, like home, office, coworking, etc.>
```

Additionaly can take into the account `$AWS_IP_DEFAULT_SECURITY_GROUP_NAME`, so `--group-name` can be omitted, and set only once.

Alternative to `--port` could be `--service` option with values like `mysql`, `postgresql`, etc., and so you default defined port for these services.

## Commands

Might be good to add the following additional commands:

- `list-locations` - lists user's available locations
- `ip` - prints user's current local IP address
- `rm` - removes the corresponding entry

## Ideas

Might able auto-delete of the entry (after timeout), if user assigns IP of the public wi-fi, like cafe, library, airport, etc., while needs only temporal access from this location. And this could be the default behaviour, unless `--durable` flag is set.
