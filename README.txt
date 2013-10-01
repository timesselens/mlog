mlog
====

WARNING: proof of concept - learning go

project scope
-------------
- mux i/o with pipes (break) or listeners (wait)
- one binary, deploy: scp /usr/bin/mlog host.example.com:~ && ssh -N -L1980:localhost:1980 ./mlog
- set initial config from http://localhost:1980
- usable from the command line as oneliner

example usage
-------------
- mlog [-http:1980] # basic, connect to http://localhost:1980 to configure

io
--
- stdin
- named pipes
- unix domain sockets
- syslog
- websocket
