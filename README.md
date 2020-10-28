# gRPC Demo

gRPC client is implemented in https://github.com/kikeyama/gorilla-sfx-demo.

## Env var

ENV_VAR | Description | Default Value
--------|-------------|--------------
`MONGO_HOST` | Hostname of MongoDB | `localhost`
`MONGO_PORT` | Port number of MongoDB | `27017`
`GRPC_PORT` | Port number of gRPC | `50051`

## Methods

Method | Description
-------|------------
`ListAnimals` | List all animals in MongoDB
`GetAnimal` | Get an animal specified with `UUID`
`CreateAnimal` | Create an animal
`DeleteAnimal` | Delete an animal
