# MoinApp Server

[![Build Status](https://travis-ci.org/MoinApp/moinapp-server.svg?branch=feature/go-rewrite)](https://travis-ci.org/MoinApp/moinapp-server)

The server for Moin.
The server is deployed to http://moinapp.herokuapp.com/ (this link is funny).

## Build

To build the server binary, run `build.sh`. The result will be an executable called `moinapp-server`. This is the entire server with all of its dependencies included.

For improved build times run `go install` once so that all dependencies are cached.

## Running

Start `moinapp-server`.
You can modify the port via the `PORT` environment variable. Setting this to `0` will cause the server to use a system-defined port.

## Bugs

Bug reports are welcome. Post them [here](https://github.com/MoinApp/moinapp-server/issues).
