# srtl
Link Shortener

## Build

Easiest way is run `make release` this will build for ubuntu 20.04 env, you can alter the docker container for a different env if you need. Otherwise you can install rocksdb libraries using the `install.sh` script, or by following [rocksdb intructions](https://github.com/facebook/rocksdb/blob/main/INSTALL.md).

If you have the required libraries installed you can run `go mod download` to download the required go modules and then either `make build` or `make run` to build/run locally.

## Try it

[It's live]("https://www.srtl.ie") so you can mess around with it