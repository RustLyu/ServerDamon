## Makefile
This microservice relies heavily on a Makefile that has multiple tools for developing, testing, and building the project. This works very well for rapid development and is compatible with Windows Subsystem for Linux (WSL), MacOS, and most other Linux variants. Here is the full list of supported make targets. 

Environmental Stuff:
* `make run` -> lints, builds and tests the service. Starts a local container running the template
* `make run local` -> lints, builds and tests the service. Starts running locally without docker
* `make version` -> Returns the git version/tag of the running service.

Build Tools:
* `make` -> lints, builds, and tests the service
* `make build` -> same as `make`
* `make build-only` -> Builds the binary only (should only be used for pipelines)
* `make docker` -> Builds and tags a local docker image

Test Tools:
* `make test` -> Executes all tests
* `make test-bench` -> Executes any benchmark tests
* `make test-coverage` -> Executes tests and provides an html coverage report
* `make test-default` -> same as `make test`
* `make test-race` -> Tests for race conditions
* `make test-verbose` -> Verbose version of `make-test`

## Configuration
### GoMods
This solution uses go mods. 

### Go Mods
This template uses [Go Mods](https://github.com/golang/go/wiki/Modules) to manage dependencies. All external dependencies are in the `go.mod` file. There is no need to run the command `go get`, just import any modules you would like to use into the project, and the package will need to be added to the `go.mod` file in the root directory. 

- Unix Command
    ```
    $ export GO111MODULE=on
    $ go mod download
    ```

### Environment
The environment configuration is initialized in the `internal/config` package. Generally, we use local, dev, test, and prod, but any configuration is supported as long as the environment variable is set to the name used in the configuration file.

### Secrets
Secrets should not be stored anywhere in source control. There are two methods for using secrets: 
- You may create a file and use the TOML format to put secret in the file (these will be excluded from source control.)
- You can set environment variables to store the secrets.

Here are the current secrets for this project. If this file is not updated, secrets are defined in `./internal/config/config.go` in the struct named `Secrets`

```
//Secrets is used to store credential data
type secrets struct {
	FedExCredentials struct {
		ClientID     string
		ClientSecret string
	}
```
### Strings
There should be no hardcoded strings used for words or sentences in this service outside of the messages package. This is to support internationalization and localization. Here is an example of how to extend the package in the `./internal/messages` directory to support your needs.

```
type SomeString int
const (
	Foo SomeString = iota
)

var SomeStrings = map[SomeString]string{
    Foo: "Bar",
}
func (s SomeString) String() string {
	return SomeStrings[s]
}

The result is messages.SomeString.Foo.String() = "Bar"
```


### Logging
The logging implements the [Logrus](https://github.com/sirupsen/logrus) package.

`[INFO]   2022-05-24 08:33:09 [pkg/logging/loggingConfig.go] [logging.(*LoggingConfig).Initialize: 31] - Logging initialized`

Here is an example:
```
// Log a warning
Log.Warn("This is a warning")

// Log an error example, err is of type error here
Log.Errorf("This is an error: %v", err)
```

In addition to logging messages, Logrus has the ability to log data with key-value pairs.
Keys needs to be strings and the value is accepted as `interface{}`
Here is an example:
```
// Log a single data field
logging.Log.withField("key", "value").Info("The log message")

// Log multiple data fields
dataFields = map[string]interface{}{
    "key1": "value1",
    "key2": 2,
}
logging.Log.WithFields(dataFields).Debug("The log message")
```

### Testing
In Go, tests belong in the same directory as file they are testing, and append the filename with `_test.go`.
For example, to write tests for a file `pkg/example/example.go`, the test file would be `pkg/example/example_test.go`

In order to test all test files in a certain directory manually, execute the following
command from the desired directory to be tested:
```
$ make test
```

Here is the full list of testing targets:
* `make test`
* `make test-bench`
* `make test-coverage`
* `make test-coverage-tools`
* `make test-default`
* `make test-race`
* `make test-short`
* `make test-verbose`

## Endpoints
There are no endpoints as this is a headless container with no ingress. 

## Pipeline
TBD 


### Docker
The microservice is hosted in a [Docker](https://docs.docker.com) container. In order to run the microservice locally, Docker must first be
installed and set up. There is a target in the `Makefile` for building the image: `make docker`. This will create the
image for the Docker container. 
```
$ make docker
...
...
$ docker images
REPOSITORY               TAG                                        IMAGE ID            CREATED             SIZE
u-label-api   a3d76dd96b3c822d471da438a613cd4d3b897c83   92e0dae9fa46        7 seconds ago       21.2MB
```