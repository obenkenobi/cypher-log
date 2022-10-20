export IGNORE_ENV=true
export MONGO_URI=mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0

export MONGO_DB_NAME=users
export MIGRATION_DIRECTORY=users

#yarn run migrate:down:all
yarn run migrate:up

export MONGO_DB_NAME=keys
export MIGRATION_DIRECTORY=keys

#yarn run migrate:down:all
yarn run migrate:up