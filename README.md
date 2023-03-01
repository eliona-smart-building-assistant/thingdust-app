# Thingdust APP
The [Thingdust app](https://github.com/eliona-smart-building-assistant/thingdust-app) enables the communication of data from [Thingdust spaces](https://thingdust.com/) to the Eliona environment.

The app collects data from configurable endpoints. Each thingdust space is automatically assigned to an eliona asset when a new endpoint is configured. Data from each space (e.g temperature, humidity and occupancy) can then be transferred to and displayed in eliona.


## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.


### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app has to be registered in Eliona. In order to register the app, an entry in the database table `public.eliona_app` must be made.


### Environment variables

- `APPNAME`: must be set to `thingdust`. Some resources use this name to identify the app inside an Eliona environment.

- `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started. (e.g. `postgres://user:pass@localhost:5432/iot`)

- `API_ENDPOINT`:  configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started. (e.g. `http://api-v2:3000/v2`)

- `API_TOKEN`: defines the secret to authenticate the app and access the API. 

- `API_SERVER_PORT`(optional): define the port the API server listens. The default value is Port `3000`. <mark>Todo: Decide if the app needs its own API. If so, an API server have to implemented and the port have to be configurable.</mark>

- `DEBUG_LEVEL`(optional): defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-utils/tree/develop/log). Not defined the default level is `info`.

### Database tables ###

<mark>Todo: Describe the database objects the app needs for configuration</mark>

<mark>Todo: Decide if the app uses its own data and which data should be accessible from outside the app. This is always the case with configuration data. If so, the app needs its own API server to provide access to this data. To define the API use an openapi.yaml file and generators to build the server stub.</mark>

The app requires configuration data that remains in the database. In order to do this, the app creates its own database schema `thingdust` during initialization. To modify and handle the configuration data the app provides an API access. Take a look at the [API specification](https://github.com/eliona-smart-building-assistant/thingdust-app/blob/develop/openapi.yaml) to see how the configuration tables should be used.

- `thingdust.config`: contains the thingdust API endpoints. Each row contains the specification of one endpoint(i.e config id, url, key, polling intervals etc.)

- `thingdust.spaces`: contains the mapping from each space uniquely defined by its configuration and project id to an eliona asset. Each row contains the specification of one endpoint(i.e config id, url, key, polling intervals etc.) The app collects and writes data separately for each configured project. The mapping is created automatically by the app.


**Generation**: to generate access method to database see Generation section below.


## References

### App API ###

The thingdust app provides its own API to access configuration data and other functions. The full description of the API is defined in the `openapi.yaml` OpenAPI definition file.

- [API Reference](https://github.com/eliona-smart-building-assistant/thingdust-app/blob/develop/openapi.yaml) shows details of the API.

**Generation**: to generate api server stub see Generation section below.


### Eliona assets###

The app creates necessary asset types and attributes during initialization. See [eliona/asset-type-thingdust_space.json](eliona/asset-type-thingdust_space.json) for details.

The thindust app writes data for each thingdust space to the eliona database. Each thingdust space is mapped to an asset with attributes: temperature, humidity and occupancy. All attributes are of the eliona subtype `Input`.


## Tools

### Generate API server stub ###

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used to generate a server stub. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

### Generate Database access ###

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation. Please note that the database connection in the `sqlboiler.toml` file have to be configured.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

