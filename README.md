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

## TODO

In order to make use of this some improvements must be made:

* Currently the checks received via POST endpoint are only stored in-memory. This data must be persisted, for example on an object storage.
* All endpoints are unprotected. At least some simple token-based authorization must be implemented.
