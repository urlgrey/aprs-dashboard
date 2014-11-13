APRS Dashboard [![Circle CI](https://circleci.com/gh/urlgrey/aprs-dashboard.png?style=badge)](https://circleci.com/gh/urlgrey/aprs-dashboard)
==============

Service API to record and query for Automated Position Reporting System (APRS) messages.

Installation
------------

### Dynamic Redis
Dynamic Redis is needed for the Geo module.  More information about Dynamic Redis can be found on the project site:
https://matt.sh/dynamic-redis

```shell
mkdir -p ~/repos
cd ~/repos
git clone https://github.com/mattsta/redis
cd redis
git checkout dynamic-redis-unstable
make
make test
```

If the build succeeded, then run the Redis server:

```shell
~/repos/redis/src/redis-server &
```

### Redis Geo Module
The Redis Geo module makes it possible to store & query for geo-tagged data in Redis.

```shell
cd ~/repos
git clone https://github.com/mattsta/krmt
cd krmt
make -j
~/repos/redis/src/redis-cli config set module-add `pwd`/geo.so
```

You should see output like the following:

```shell
79865:M 12 Nov 07:04:45.783 * Loading new [/Users/scott/repos/krmt/geo.so] module.
79865:M 12 Nov 07:04:45.783 * Added command geoadd [/Users/scott/repos/krmt/geo.so]
79865:M 12 Nov 07:04:45.783 * Added command georadius [/Users/scott/repos/krmt/geo.so]
79865:M 12 Nov 07:04:45.783 * Added command georadiusbymember [/Users/scott/repos/krmt/geo.so]
79865:M 12 Nov 07:04:45.783 * Added command geoencode [/Users/scott/repos/krmt/geo.so]
79865:M 12 Nov 07:04:45.783 * Added command geodecode [/Users/scott/repos/krmt/geo.so]
79865:M 12 Nov 07:04:45.783 * Module [/Users/scott/repos/krmt/geo.so] loaded 5 commands.
79865:M 12 Nov 07:04:45.783 * Running load function of module [/Users/scott/repos/krmt/geo.so].
OK
```

### APRS Dashboard

#### Build
```shell
make build
```

#### Test
```shell
export APRS_REDIS_HOST=":6379"
make test
```

#### Benchmark
```shell
export APRS_REDIS_HOST=":6379"
make bench
```

#### Run

**Note:** To limit access to the PUT API, optionally set the `APRS_API_TOKENS` environment variable with a comma-separated list of API tokens.  Default behavior is to allow access without the use of a token.
```shell
export APRS_API_TOKENS="secret123"
export APRS_REDIS_HOST=":6379"
./aprs-dashboard
```

API
---

### Record an APRS Message
Send an APRS message in the following JSON format:

| Field  | Required?  | Type | Description  |
|---|---|---|---|
| data  | yes  | string | ASCII-encoded APRS message |
| is_ax25  |  yes | boolean  | Indicates whether the message is AX.25 encoded |

#### Sample Payload
```json
{
    "data": "WX4GSO-9>APN382,qAR,WD4LSS:!3545.18NL07957.08W#PHG5680/R,W,85NC,NCn Mount Shepherd Piedmont Triad NC",
    "is_ax25": false
}
```

#### Example
```shell
curl -X PUT -H "Content-Type: application/json" http://127.0.0.1:3000/api/v1/message -d @examples/sample_message.json
```

If access is limited by API token (see "Run" section), then include an `X-API-KEY` header:
```shell
curl -X PUT -H "X-API-KEY: secret123" -H "Content-Type: application/json" http://127.0.0.1:3000/api/v1/message -d @examples/sample_message.json
```

Observe output resembling the following in the APRS Dashboard console output:
```shell
[martini] Started PUT /api/v1/message for 127.0.0.1:63695
[martini] Completed 200 OK in 297.553us
```

### Retrieve Paginated List of Messages Sent by Callsign

#### Example
```shell
curl -H "Accept: application/json" http://127.0.0.1:3000/api/v1/callsign/WX4GSO?page=1
```

##### Response
| Field  | Type | Description  |
|---|---|---|
| page  | integer | Page associated with these results |
| number_of_pages  |  integer  | Total number of pages for this callsign |
| total_number_of_records  |  integer  | Total number of records for this callsign |
| total_number_of_records  |  integer  | Total number of records for this callsign |
| records  |  Array of ```APRS Message```  | An array of APRS Messages, can be empty |

##### APRS Message
| Field  | Type | Description  |
|---|---|---|
| timestamp  | integer | Unix epoch time when APRS message was received |
| src_callsign  |  string  | Callsign that sent the message |
| dst_callsign  |  string  | Callsign of message's intended recipient |
| latitude  |  float  | Latitude |
| longitude  |  float  | Longitude  |
| includes_position  |  boolean  | Indicates whether message contained location information |
| altitude  |  float  | Altitude (meters) |
| speed  |  float  | Speed (kilometers/hour) |
| course  |  unsigned integer  | Direction in degrees (0-360) |
| weather_report  |  Weather Report  | Weather data, can be null |
| raw_message  |  string  | Raw APRS messages received by service |
