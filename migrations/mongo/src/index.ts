import { mongoMigrateCli } from 'mongo-migrate-ts';
import * as dotenv from 'dotenv'

dotenv.config()

console.log(process.env.MIGRATION_DIRECTORY)
mongoMigrateCli({
    uri: process.env.MONGO_URI || "mongodb://username:passwordf@localhost:27017",
    database: process.env.MONGO_DB_NAME || "db",
    migrationsDir: `${__dirname}/${process.env.MIGRATION_DIRECTORY}`,
    migrationsCollection: 'migrations',
});