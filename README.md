# aws-ip

Sync/set AWS security group entry (by description) with current user's local public IP address. 

I.e. if the access to AWS resources is restricted by the IP address, and the user/developer doesn't have fixed static IP 
address or user works frequently from different locations with different IP addresses.

```
NAME:
   aws-ip - Sync/set AWS security group entry (by description) with current user's local IP address

USAGE:
   aws-ip     
        --group-name    <comma separated affected EC2 security groups> 
        --port          <port>  
        --location      <free text/code of the current user's location, like home, office, coworking, etc.>

VERSION:
   1.0.0

AUTHOR:
   (c) Aliaksandr Kazlou

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --group-name value  comma separated affected EC2 security groups
   --port value        port number (default: 0)
   --location value    free text/code of the current user's location, like home, office, coworking, etc.
   --help, -h          show help
   --version, -v       print the version
```

Table of Contents
=================

* [Idea](#idea)
* [Installation](#installation)
* [Usage](#usage)
* [AWS Connection](#aws-connection)
* [Improvements](#improvements)
* [Copyright](#copyright)

## Idea

Tool locates the entry within specified EC2 security group name(s) by description and updates it. The trick is around 
description property. Where description will be set to 

    <IAM username implicitly derived from the used AWS profile>-<location>

where `location` is something user provides, like home, office, coworking, etc.

It solves the problem if the access to the AWS resource (like DB, for example) is restricted by the IP address, and 
user/developer doesn't have fixed static IP address (which is changing frequently) or user works frequently from 
different locations with different IP addresses.

## Installation

Use the `go` command:

    $ go get github.com/zshamrock/aws-ip

## Usage

    $ aws-ip --group-name db --port 3306 --location home

## AWS Connection

Connection to the AWS is established using profile credentials. Currently it relies on the environment variables entirely,
i.e. `$AWS_DEFAULT_PROFILE` or `$AWS_PROFILE`.

## Improvements

Additionally, can take into the account `$AWS_IP_DEFAULT_SECURITY_GROUP_NAME`, so `--group-name` can be omitted, and set only once.

Alternative to `--port` could be `--service` option with values like `mysql`, `postgresql`, etc., and so you default defined port for these services.

### Commands

Might be good to add the following additional commands:

- `list-locations` - lists user's available locations
- `ip` - prints user's current local IP address
- `rm` - removes the corresponding entry

### Extra features

Might be able auto-delete of the entry (after timeout), if user assigns IP of the public wi-fi, like cafe, library, airport, etc., while needs only temporal access from this location. And this could be the default behaviour, unless `--durable` flag is set.

## Copyright

Copyright (C) 2021-2021 by Aliaksandr Kazlou.

`aws-ip` is released under MIT License.                                                                                                                       
See [LICENSE](https://github.com/zshamrock/aws-ip/blob/master/LICENSE) for details.
