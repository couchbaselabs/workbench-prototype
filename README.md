# workbench-prototype

Table of Contents
==================
- [workbench-prototype](#workbench-prototype)
- [Table of Contents](#table-of-contents)
  - [Setup](#setup)
    - [Running using Docker](#running-using-docker)
    - [Running natively](#running-natively)
      - [Installing dependencies](#installing-dependencies)
        - [Installing Go](#installing-go)
        - [Installing SQLCipher](#installing-sqlcipher)
        - [Installing npm](#installing-npm)
      - [Building the backend](#building-the-backend)
      - [Building the UI](#building-the-ui)
      - [Running the project](#running-the-project)
  - [Auto-Configuration from Prometheus](#auto-configuration-from-prometheus)
  - [Prometheus Monitoring](#prometheus-monitoring)
  - [Contributing](#contributing)
    - [Unit Testing](#unit-testing)


### Running using Docker

To build using a container: `docker build -t <tag> .`
This uses a multistage build to compile all the executables then transfer them to a runtime image.

This can then be run as a container: `docker run --rm -it -P <tag>`.
By default this will run the `entrypoint.sh` script which launches the `workbench-prototype` with some default parameters.
The container also includes a shell and is Alpine Linux based rather than being distroless currently.
Provide extra arguments to run a shell, e.g. `docker run --rm -it <tag> bash` will launch a Bash shell inside the container.

### Running natively

#### Installing dependencies

To run the project you will need the following three dependencies

- Go 1.16 or above.
- SQLCipher
- npm 6.X

If you wish to also use Fluent Bit directly on your laptop (not in a VM), you will also need `msgpack` and `mbedtls`.
To install these, run `brew install msgpack mbedtls`.

##### Installing Go

You can install Go either by getting the binary release from https://golang.org/dl/ or using your favourite package
manager. Once installed check the version is 1.16 or above by running:

```
> go version
go version go1.16 darwin/amd64
```

##### Installing SQLCipher

I recommend using your package manager of choice on MacOS you can use homebrew as follows:

```
> brew install sqlcipher
```

You can check it is installed and working by running
```
> sqlcipher
SQLite version 3.33.0 2020-08-14 13:23:32 (SQLCipher 4.4.2 community)
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> .exit
```

##### Installing npm

To install `npm` you will need to install `NodeJS`. You can get the binary release from https://nodejs.org/en/download/
or use a package manager in Mac you can use homebrew as follows:

```
> brew install node
```

To confirm `npm` is installed run (Make sure the version is 6.X):

```
> npm version
{
  npm: '6.14.11',
  ares: '1.16.1',
  brotli: '1.0.9',
  cldr: '37.0',
  icu: '67.1',
  llhttp: '2.1.3',
  modules: '88',
  napi: '7',
  nghttp2: '1.41.0',
  node: '15.0.1',
  openssl: '1.1.1g',
  tz: '2019c',
  unicode: '13.0',
  uv: '1.40.0',
  v8: '8.6.395.17-node.15',
  zlib: '1.2.11'
}
```

#### Building the backend

For native build and deployment the backend is written in Go and uses `go mod` for dependency management. To download the dependencies you will first
need to export the environmental variables for CGO (Note depending on your setup this may not be required, for most
Mac users it will be. You can try without exporting and if you hit the error export and retry):

*Note:* The path may vary between machines, and the quotes may need to be removed depending on the shell you use.
```
export CGO_ENABLED=1
export CGO_LDFLAGS="-L/usr/local/Cellar/openssl@1.1/1.1.1i/lib"
export CGO_CPPFLAGS="-I/usr/local/Cellar/openssl@1.1/1.1.1i/include"
export CGO_CFLAGS="-I/usr/local/Cellar/openssl@1.1/1.1.1i/include"
export CGO_CXXFLAGS="-I/usr/local/Cellar/openssl@1.1/1.1.1i/include"
```

If during the download/build process Go cannot find them you will see the following error:
```
sqlite3-binding.c:24328:10: fatal error: 'openssl/rand.h' file not found
```

Once the environmental variables are setup you can download the dependencies doing:
```
> go mod download
```

To build the backend use the following:
```
> mkdir build # if you already have a build directory inside workbench-prototype feel free to ignore
> go build -o ./build ./cluster-monitor/cmd/workbench-prototype
```

This will build the backend to `./build/workbench-prototype`. To check that it worked do:
```
> ./build/workbench-prototype
NAME:
   Couchbase Multi Cluster Manager - Starts up the Couchbase Multi Cluster Manager

USAGE:
   workbench-prototype [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --sqlite-key value   The password for the SQLiteStore (default: "password") [$CB_MULTI_SQLITE_PASSWORD]
   --sqlite-db value    The path to the SQLite file to use. If the file does not exist it will create it.
   --cert-path value    The certificate path to use for TLS [$CB_MULTI_CERT_PATH]
   --key-path value     The path to the key [$CB_MULTI_KEY_PATH]
   --log-level value    Set the log level, options are [error, warn, info, debug] (default: "info")
   --http-port value    The port to serve HTTP REST API (default: 7196)
   --https-port value   The port to serve HTTPS REST API (default: 7197)
   --ui-root value      The location of the packed UI (default: "./ui/dist/app")
   --max-workers value  The maximum number of workers used for heartbeats (defaults to 75% of the number of CPUs) (default: 0)
   --help, -h           show help (default: false)
   --version, -v        print the version (default: false)
```

#### Building the UI

The UI is written in Typescript, CSS and HTML and has to be built before use. For more about developing the UI see the
UI [README](./ui/README.md). To build the UI do the following.

```
> cd ui
> npm install
> npm run build
```

This will build the UI in `./ui/dist/app`. This is default location the backend server checks to serve the UI files.

#### Running the project

To run `workbench-prototype` you will need to give it certificates to use for the HTTPS server. For development purposes
you may wish to do the following to create  the certificates

```
> mkdir priv
> cd priv
> openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -nodes -days 365
  # You will need to fill in the requested information
```

You will also need a location to store the encrypted SQLite database. For development purposes you may wish to create
a directory called `data` inside the `workbench-prototype` directory.

Once you have everything setup to start up `workbench-prototype` you can run
```
> ./build/workbench-prototype --sqlite-db ./data/data.sqlite --sqlite-key password --cert-path ./priv/cert.pem \
  --key-path ./priv/key.pem --log-level debug --ui-root ./ui/dist/app
2021-03-04T15:27:05.128Z INFO (Main) Maximum workers set to 9
2021-03-04T15:27:05.302Z INFO (Manager) Starting {"frequencies": {"Heart":60000000000,"Status":300000000000,"Janitor":21600000000000}}
2021-03-04T15:27:05.307Z INFO (Heart Monitor) Starting monitor {"frequency": 60}
2021-03-04T15:27:05.307Z INFO (Status Monitor) Starting monitor
2021-03-04T15:27:05.307Z INFO (Status Monitor) (API) Starting monitor {"frequency": 300}
2021-03-04T15:27:05.307Z INFO (Manager) Started
2021-03-04T15:27:05.307Z INFO (Manager) (HTTPS) Starting HTTPS server {"port": 7197}
2021-03-04T15:27:05.307Z INFO (Manager) (HTTP) Starting HTTP server {"port": 7196}
2021-03-04T15:27:05.307Z DEBUG (Status Monitor) (API) API check tick
2021-03-04T15:27:05.308Z DEBUG (Status Monitor) (API) API check tick {"elapsed": "184.55µs"}
```

You can also the option `--log-dir` to give it a location to persist the logging to.

The REST endpoints are defined in [routes.go](./cluster-monitor/pkg/manager/routes.go).

## Auto-Configuration from Prometheus

If you have a Prometheus instance set up to monitor your Couchbase Server nodes, `workbench-prototype` can use it to automatically discover them.

For this you will need to pass the following command line parameters (or environment variables if you are running in a container):

* `--prometheus-url` (`CB_MULTI_PROMETHEUS_URL`): the base URL of your Prometheus instance (e.g. `http://localhost:9090`)
* `--prometheus-label-selector` (`CB_MULTI_PROMETHEUS_LABEL_SELECTOR`): specifies which labels your Couchbase Server targets have. Each selector must have a label name and a value, separated by an equals sign (e.g. `job=couchbase-server`). Multiple selectors can be specified, separated by spaces, in which case they will all need to match. Passing an empty string or omitting this parameter will match *all* Prometheus targets. (To see what labels your targets have, visit `http://your.prometheus.host/prometheus/targets`.)
* `--couchbase-user` and `--couchbase-password` (`CB_MULTI_COUCHBASE_USER` and `CB_MULTI_COUCHBASE_PASSWORD`): the username and password used to authenticate against found Couchbase Server nodes

Note that, when Prometheus auto-discovery is enabled, `workbench-prototype` will assume all your clusters are in Prometheus and stop monitoring any that are not.

## Prometheus Monitoring

workbench-prototype exports metrics to prometheus for monitoring - To set up monitoring, please refer to the wiki: [Setup](https://github.com/couchbaselabs/workbench-prototype/wiki/Setup#prometheus-setup).

For a full build in a container set up - check out the [observability project](https://github.com/couchbaselabs/observability).

## Contributing

For full details of the contribution process, please see [CONTRIBUTING.md](https://github.com/couchbaselabs/workbench-prototype/blob/master/CONTRIBUTING.md).

### Unit Testing

As much of the code as reasonably possible should be covered by unit tests. We use the standard Go `testing` library, with [testify](https://pkg.go.dev/github.com/stretchr/testify) for assertions and mocking.

Mocks for the various interfaces are auto-generated using Mockery and `go generate`. To update the mocks after changing an interface, install [Mockery](https://github.com/vektra/mockery), then run:

```shell
go generate ./...
```

If you add a new interface, make sure it has a `//go:generate` line to ensure mocks are automatically updated, for example:

```go
package foo

//go:generate mockery --name FooIFace

type FooIFace interface {
	// ...
}
```