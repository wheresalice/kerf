# Kerf

Kerf is a Kafka performance testing tool.  The main, and currently only, functionality is to publish lots of messages to a Kafka broker and then calculate number of messages per second.

There's currently very little in the way of configuration options.

The goal is a non-Java version of kafka-*-perf-test but with custom messages