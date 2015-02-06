CLIPR
=============

CLIPR is a server that the [CloudFoundry CLI](https://github.com/cloudfoundry/cli) 
can interact with to browse and install plugins. This documentation covers how to run CLIPR itself,
as well as how to create your own implementation if desired.

Running CLIPR
=============
1. You need to have [git](http://git-scm.com/downloads) installed
1. Clone this repo `git clone https://github.com/cloudfoundry-incubator/cli-plugin-repo`
1. Build the project by running the build script `./bin/build`
1. Invoke the binary `./main.exe` with the following options
  - `-n`: IP Address for the server to listen on, default is `0.0.0.0`
  - `-p`: Port number for the server to listen on, default is `8080`
1. Add the running server to your CLI via the `cf add-plugin-repo` command
1. Browse & install plugins!
  
Creating your own Plugin Repo Server
=============
Alternatively, you can create your own plugin repo implementation. The server must meet the requirements:
- server must have a `/list` endpoint which returns a json object that lists the plugin info in the correct form
```
{"plugins":
  [{"name":"echo",
  "description":"echo repeats input back to the terminal",
  "version":"0.1.4",
  "date":"0001-01-01T00:00:00Z",
  "company":"",
  "author":"",
  "contact":"feedback@email.com",
  "homepage":"http://github.com/johndoe/plugin-repo",
  "binaries":
    [{"platform":"osx",
    "url":"https://github.com/johndoe/plugin-repo/raw/master/bin/osx/echo",
    "checksum":"2a087d5cddcfb057fbda91e611c33f46"},
    {"platform":"windows64",
    "url":"https://github.com/johndoe/plugin-repo/raw/master/bin/windows64/echo.exe",
    "checksum":"b4550d6594a3358563b9dcb81e40fd66"}]}]
}
```
