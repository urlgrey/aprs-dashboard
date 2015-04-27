APRS Dashboard [![Circle CI](https://circleci.com/gh/urlgrey/aprs-dashboard.png?style=badge)](https://circleci.com/gh/urlgrey/aprs-dashboard)
==============

Service API to record and query for Automated Position Reporting System (APRS) messages.

Installation
------------
APRS Dashboard is designed to run in a [Docker](https://www.docker.com/) container.  Hopefully you've come to know and love the flexibility and ease that comes with using Docker.  [APRS Dashboard Docker images](https://registry.hub.docker.com/u/urlgrey/aprs-dashboard/) are available on DockerHub.

APRS messages are stored in a Redis database.  The [Dynamic Redis](https://matt.sh/dynamic-redis) fork is used so that geo searches by latitude-longitude can be performed on Redis data.  These instructions make use of a [Docker image for the Dynamic Redis server](https://registry.hub.docker.com/u/urlgrey/dynamic-redis/).

```shell
docker pull urlgrey/dynamic-redis
docker pull urlgrey/aprs-dashboard

# Run Redis container
sudo docker run --name aprs_db -d -p 6379:6379 -v /home/skidder/git/docker-dynamic-redis/redis.conf:/usr/local/etc/redis/redis.conf urlgrey/dynamic-redis:latest redis-server /usr/local/etc/redis/redis.conf

# Run APRS Dashboard, linking it to the Redis container
sudo docker run -d --link aprs_db:db -p 3000:3000 urlgrey/aprs-dashboard:latest
```

**Note:** To limit access to the PUT API, optionally set the `APRS_API_TOKENS` environment variable with a comma-separated list of API tokens.  Default behavior is to allow access without the use of a token.
```shell
sudo docker run -d --link aprs_db:db -e APRS_API_TOKENS="secret123" -p 3000:3000 urlgrey/aprs-dashboard:latest
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
