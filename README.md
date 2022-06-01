# checkproxy

`checkproxy` is a proof of concept to allow arbitrary scripts to be monitored via
HTTP checkers, eg. Pingdom.

`checkproxy` is a tiny web server that exposes two endpoints: 

The POST endpoint at `http://{BASE_URL}/checks/{check_name}` reads the body of the request
and stores it. This endpoint can be called from any script, for example via curl:

```
curl -XPOST --data-binary '{"ok": true}' 'http://{BASE_URL}/checks/{check_name}?status=200' -H 'X-Auth-Token:{token}' -v
```

The GET endpoint at `http://{BASE_URL}/checks/{check_name}` returns the exact body that was
previously written by the POST request. This endpoint can be configured as a check endpoint
in the HTTP checker.

`checkproxy` can easily be deployed as Azure Function or AWS Lambda.

## Configuration

`checkproxy` takes a few arguments:

```
# checkproxy -h
Usage of checkproxy:

(alternatively the -c parameter can be specified via environment variable 'CONFIG_FILE')

  -a string
    	ip:port to listen on when run locally
  -c string
    	name of the config file (default "./config.yaml")
  -m string
    	mode, can me either 'local', 'azurefunc' or 'awslambda' (default "local")
  -v	print version information
```

Note that the configuration file (`-c` or environment variable `CONFIG_FILE`) can be provided using different locations:

* Use a simple file path to indicate that a _local file_ is used.
* Use a path of the pattern `s3://[bucket_name]/[object_path]` to indicate that file is stored on s3 is used. All [official authentication methods](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/) are allowed.
* Use  a path of the pattern `blob://[storage_account]/[container_name]/[object_path]` to indicate that file is stored on Azure Blob storage. The `AZURE_STORAGE_ACCOUNT_KEY` environment variable must be set in order to grant access.

The configuration file holds should look something like this:

```
---
presistance_base: blob://stsaccheckproxypers/persistence/checks/
auth:
  rw_tokens:
    - ThisIsATokenThatAllowsRWAccess
    - ThisIsAnotherTokenThatAllowsRWAccess
  r_tokens:
    - ThisIsATokenThatAllowsRAccess
  w_tokens:
    - ThisIsATokenThatAllowsWAccess
```

The `presistance_base` path specified a directory/prefix where states are stored. Note that the same rules as for the `-c` parameter apply
in regards of the path specification.

The `rw_tokens`, `r_tokens` and `w_tokens` specify the usable tokens to access the end points. `r_tokens` allow to perform `GET` request,
`w_tokens` allow for `POST` requests, `rw_tokens` allow for both.

## Usage

If the service is up and running use any HTTP Client to set and get states. To set a state use a `POST` request:

```
curl -XPOST --data-binary '{"ok": true}' 'http://{BASE_URL}/checks/{check_name}?status=200&validity_duration=30m' -H 'X-Auth-Token:{token}' -v
```

Note the that the required query parameter `status` allows to specify the status code that will be provided on a `GET` request.
The optional `validity_duration` query parameter allows to specify a duration during which the result will be stored. After this duration, a
`Gateway Timeout` will be returned in order to prevent a state to be returned to infinity. Default is 10m. Refer to the [official Go documentation](https://pkg.go.dev/time#ParseDuration)
to learn in which format the duration can be provided.
