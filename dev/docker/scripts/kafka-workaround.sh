#!/bin/sh

# Disable zookeeper connection
sed -i '/KAFKA_ZOOKEEPER_CONNECT/d' /etc/confluent/docker/configure
sed -i 's/cub zk-ready/echo ignore zk-ready/' /etc/confluent/docker/ensure

# Required for KRaft
echo "kafka-storage format --ignore-formatted -t NqnEdODVKkiLTfJvqd1uqQ== -c /etc/kafka/kafka.properties" >> /etc/confluent/docker/ensure