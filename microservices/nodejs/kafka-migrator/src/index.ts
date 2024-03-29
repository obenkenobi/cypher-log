import { Kafka, SASLOptions, Admin } from 'kafkajs'
import * as dotenv from 'dotenv'

// Todo: add args to specify what task to do when migrating

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

const migrateTask1 = async (admin: Admin) => {
  console.log("Begin migrate task 1")
  const userChange1Topic = "user-change-1"
  const noteService = "note-service"
  const keyService = "key-service"
  const retry = "retry"
  const deadLetter = "dead-letter"
  await admin.createTopics({
    validateOnly: false,
    waitForLeaders: true,
    timeout: 10000,
    topics: [
      {
      topic: userChange1Topic,
      numPartitions: 6,
      replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${noteService}-${retry}-1`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${noteService}-${retry}-2`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${noteService}-${retry}-3`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${noteService}-${retry}-4`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${noteService}-${deadLetter}`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${keyService}-${retry}-1`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${keyService}-${retry}-2`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${keyService}-${retry}-3`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${keyService}-${retry}-4`,
        numPartitions: 6,
        replicationFactor: 2
      },
      {
        topic: `${userChange1Topic}-${keyService}-${deadLetter}`,
        numPartitions: 6,
        replicationFactor: 2
      },
    ]
  })
  console.log(`Created topic ${userChange1Topic} and associated retry topics and dead-letters`)
  console.log("End migrate task 1")
}


const task = async () => {
  const admin = kafka.admin();

  await admin.connect()
  console.log("Connected to admin")
  await migrateTask1(admin)
  console.log("Beginning migration")

  console.log("Ending migration")
  console.log("topics created")

  await admin.disconnect()
}

task().then(() => console.log('done'));
