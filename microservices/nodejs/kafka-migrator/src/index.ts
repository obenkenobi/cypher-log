import { Kafka, SASLOptions } from 'kafkajs'
import * as dotenv from 'dotenv'


if (process.env.IGNORE_ENV_FILE !== "true") {
  dotenv.config()
}

const { KAFKA_USERNAME: username, KAFKA_PASSWORD: password } = process.env;
const sasl: SASLOptions | undefined = username && password ? { username, password, mechanism: 'plain' } : undefined;
const ssl = !!sasl;

console.log((process.env.KAFKA_BOOTSTRAP_SERVERS || "").split(","))
const kafka = new Kafka({
  clientId: 'npm-slack-notifier',
  brokers: (process.env.KAFKA_BOOTSTRAP_SERVERS || "").split(","),
  ssl,
  sasl
});

const admin = kafka.admin();

(async ()=> {
  await admin.connect()
  try {
    console.log("Hello world!!!")
  } finally {
    await admin.disconnect()
  }
})()

