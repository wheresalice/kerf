# Kerf

Kerf is a Kafka performance testing tool.  The main, and currently only, functionality is to publish lots of messages to a Kafka broker and then calculate number of messages per second.

There's currently very little in the way of configuration options.

The goal is a non-Java version of kafka-*-perf-test but with custom messages

## Usage

Build it with:

```
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.14 go build -v
```

Kerf can then be run in the following fashion - sending 1000000 records to partition zero of the TEST_ALICE Kafka topic.

```
./kerf publish -b localhost:9092 -t TEST_ALICE -u $APIKEY -k "$SECRET" -v message.json -s SASL_SSL -m 1000000 -p 0
```

It should be noted that Kerf does not validate the message in any way, so it is important to do this separately - sending a million invalid messages is bad.
