# Visualization Tool for Akita

## System requirements

1. NodeJS and NPM. 
2. MySQL (can be installed on a remote machine). 

## How to use the visualization tool

### Build webpages

After cloning the vis repository, entery the `gcn3` directory. Run `npm install` and then run `npm run build`. 

### Collect Trace

Set environment variable `AKITA_TRACE_USERNAME`, `AKITA_TRACE_PASSWORD`. Optionally set environment variable `AKITA_TRACE_IP` and `AKITA_TRACE_PORT` is the MySQL server is on a remote machine or is not on the default port.

Run any simulation of MGPUSim with the commandline argument `-trace-vis`. Copy the database name as printed in a line like `Trace is collected in database: XXXXXXX`. 

### Build Server

In the `server` directory, run `go build`. 

### Start Server

In the `server` directory, run `./server gcn3 XXXXXXXX`, where `XXXXXXXX` is the database you copied from an earlier step.

## Develop

First, install the dependencies for the website by running `npm install` in the `gcn3` folder. Then run `npm run watch` to update client-side code automatically when you update any source code.

To compile the server execuable, go to the `server` folder and run `go build`. The `go build` command will automatically download the dependencies.

To start the kernel, run `./server gcn3 [your_database_name]`. Possible flags include changing port that vis runs on in localhost. Use the flag (per example) `http=:8888`. The flag can be used as such: `./server -http=:8888 gcn3 [your_db_name]`.
