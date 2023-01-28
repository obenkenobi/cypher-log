import { Kafka, SASLOptions, Admin } from 'kafkajs'
import * as dotenv from 'dotenv'


if (process.env.IGNORE_ENV_FILE !== "true") {
  dotenv.config()
}

const { KAFKA_USERNAME: username, KAFKA_PASSWORD: password } = process.env;
const sasl: SASLOptions | undefined = username && password ? { username, password, mechanism: 'plain' } : undefined;
const ssl = !!sasl;

console.log("Configuring to connect to brokers ", (process.env.KAFKA_BOOTSTRAP_SERVERS || ""))
const kafka = new Kafka({
  clientId: 'kafka-migrator',
  brokers: (process.env.KAFKA_BOOTSTRAP_SERVERS || "").split(","),
  ssl,
  sasl
});


const task = async () => {
  const admin = kafka.admin();

  await admin.connect()
  console.log("Connected to admin")

  console.log("creating topics")
  await admin.createTopics({
    validateOnly: false,
    waitForLeaders: true,
    timeout: 10000,
    topics: [{
      topic: "user-0",
      numPartitions: 6,
      replicationFactor: 2
    }]
  })
  console.log("topics created")

  await admin.disconnect()
}

task().then(() => console.log('done'));