# ViiV Veeva Integration Service / CLI

#### This service and command-line tool is designed to select Veeva-related data from a Luckie database, pre-process that data according to GSK specifications, and ultimately pass that data to GSK via their SOAP web service interfaces.

### Configuration

This tool can be run in a Docker container automatically, or manually at the command-line.

Most options are specified in a configuration file called ***veeva.toml***, that is in an AWS S3 bucket called ***luckie-veeva***, and accessible via the security credentials of the ***veeva*** IAM user.

The Luckie source database and its credentials are configured in ***veeva.toml***, along with the configuration information for the relevant GSK SOAP web service interfaces.

The code herein makes specific SQL calls, so the database schemas and the SOAP WSDLs must remain immutable, unless future changes to this code base are acceptable.

### Docker Deployment - Automated Operation

The tool comes with its own Dockerfile, and should usually run as an automated service that is initiated each day by its built-in clock:

When run automatically within a Docker container, the following environmental variables must be set:

- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_REGION
- HOUR
- MINUTE

The HOUR and MINUTE environmental variables are used to set the start clock for each day's integration run.  HOUR must be an integer between 0 and 23, and MINUTE must be an integer between 0 and 59.

### Command-Line Interface - Manual Operation

You must either be in the same directory as the binary, or have the binary included in your PATH:

##### e.g. .bashrc entry on a Mac
```
PATH=$PATH:/path/to/veeva/binary
```

##### ~/.aws/credentials
```
[veeva]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```
(Credentials above are merely an example, and are not valid.)

##### ~/.aws/config
```
[profile veeva]
region = us-east-1
output = table
```

Current ***credentials*** and ***config*** files are ready and available to use in the AWS S3 bucket called ***luckie-veeva***.

Assuming that the correctly-configured ***veeva.toml*** file is in the AWS S3 bucket called ***luckie-veeva***, and is accessible via the security credentials of the ***veeva*** IAM user specified in the ***~/.aws/credentials*** file. then you merely type the following command:
```
veeva now veeva
```

or more generically:

```
veeva now <aws profile>
```

### Manual Binary Builds

Docker Deployment performs automatic builds upon the Linux-based Docker cluster, so the follow information in this section only applied to manual binary builds for use within the command-line interface.

The following instructions assume that one is building the binary on a Mac, regardless of the platform being deployed to.  For other platforms, please refer to the official Go documentation for the equivalent Go commands for those platforms.

It is also assumed that Go has already been successfully installed on the Mac performing the build, prior to issuing the following commands.

It also assumes that once Go has been successfully installed, that the following command has been issued from your GOPATH:
```
go get -u github.com/chrisbenson/viiv-veeva-integration
```

When building on a Mac for Mac deployment, please use this command from your GOPATH:
```
go install github.com/chrisbenson/viiv-veeva-integration
```
The path to the resulting binary is:
```
GOPATH/bin/veeva
```
When building on a Mac for Linux deployment, please use this command from your GOPATH:
```
env GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/veeva -v src/github.com/chrisbenson/viiv-veeva-integration
```
The path to the resulting binary is:
```
GOPATH/bin/linux_amd64/veeva
```

When building on a Mac for Windows deployment, please use this command from your GOPATH:
```
env GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64/veeva -v src/github.com/chrisbenson/viiv-veeva-integration
```
The path to the resulting binary is:
```
GOPATH/bin/windows_amd64/veeva
```
