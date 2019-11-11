# NetworkWatcherCheck

A GoLang usage reference of NetWorkWatcher Connectivity Check Service API, which is a
service that Azure provides as part of the Network Watcher. What it does is, bascically
check if the source and the destination you specify are able to communicate. To get more
information about the network watcher and the services it provides, you can check this
[link](https://docs.microsoft.com/en-us/azure/network-watcher/).

For creating this project, we are consuming the REST API for the Connectivity Check, that
is shown in this [link](https://docs.microsoft.com/en-us/rest/api/network-watcher/networkwatchers/checkconnectivity).

We created this to perform network tests, this way we are always sure that the
infrastructure we are creating is behaving as expected.

## Installing dependencies

To run this project you need some dependencies installed. The code depencies are managed
by dep, so after you get that there is not too much to be concerned about.
So, dependencies are?

- **go:** golang programming language. [Get it](https://golang.org/)
- **dep:** dependency manager for golang, reads the Gopkg.toml file and installs all the
dependencies for a project.
[Get it](https://github.com/golang/dep)
- **mage:** a make-like build tool using Go. Similar to Makefiles, but using golang.
[Get it](https://github.com/magefile/mage)

After you install all the requirements, go to the root of the project and run:

``` shell
dep ensure
```

This will download all the dependencies that are inside the ```Gopkg.toml```.

## Setting the environment

To run the application correctly, you will need to export some environment variables,
related to service principal in Azure. These are:

- **ARM_CLIENT_ID:** client Id from service principal.
- **ARM_CLIENT_SECRET:** client password from service principal.
- **ARM_SUBSCRIPTION_ID:** subscription id of Azure subscription.
- **ARM_TENANT_ID:** tenant Id from service principal.

## Running Tests with Mage

As said before, ```mage``` is a make-like build tool and we are using that to test the
application. In order to do that, you have four methods written in mage (that are the
name of the methods). To run them, you need to do:

``` shell
mage <option>
```

and the options are:

- **UnitTest:** run all the unit tests (prefix: TestUN_)
- **Integration:** run all the integration tests (prefix: TestIT_)
- **Full:** runs both the unit and the integration tests
- **Format:** format all the files following the best practices of golang.

## How to use the code

To execute the connection, all you need to do is create a new connection test and run it.
To do that, inside the ```pkg/networkcheck```, get a new instance o testing with
the method:

``` go
conTest := NewConnectionTest()
```

then, run it like:

``` go
conTest.CheckVMConnection()
```

that will return true, in case of connection and false otherwise.