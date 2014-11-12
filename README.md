APRS Dashboard
==============

Service API to record and query for Automated Position Reporting System (APRS) messages.

Installing Dynamic Redis
========================
Dynamic Redis is needed for the Geo module that allows for geo-based searches of data in Redis.  More information about Dynamic Redis can be found on the project site:
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

Build and Load the Geo Module
=============================

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

Building APRS Dashboard
=======================
```shell
make build
```

Running APRS Dashboard
=======================
```shell
APRS_REDIS_HOST=":6379"
./aprs-dashboard
```

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