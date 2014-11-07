APRS Dashboard
==============

Service API to record and query for Automated Position Reporting System (APRS) messages.

Record a sample message
=======================
Start the Redis server if not already running.

Run the APRS Dashboard server process:
```shell
$ ./aprs-dashboard
```

Issue a CURL to the server to record data in the ```examples/sample_message.json``` file provided.
```shell
$ curl -X PUT -H "Content-Type: application/json" http://127.0.0.1:3000/api/v1/message -d @examples/sample_message.json
```

Observe output resembling the following in the APRS Dashboard console output:
```shell
[martini] Started PUT /api/v1/message for 127.0.0.1:63695
[martini] Completed 200 OK in 297.553us
```

Observe the text ```OK``` in the terminal where the ```curl``` command was run.