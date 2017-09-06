# Acquire driver

This package allows to build a driver to play Acquire games using the Sackson server.

## Usage

The provided Dockerfile specifies a container that will build the Sackson plugin in this directory when
you run it.

`docker-compose up --build`

Once `acquire.so` is built, just copy it to the same machine / container where Sackson server is running, placing it in `/usr/lib/sackson`, and restart Sackson server to load it. Take into account that this will
only work in Linux systems, as Go plugin feature is only supported in Linux ATM.
