# snacks

A small collection of basically (poorly) implemented tools for pen testing.

This project is more of a learning experience in golang so will likely fall short of a lot of go best practices

### Clone
```bash
go get github.com/andrewflbarnes/snacks
# or
git clone git@github.com:andrewflbarnes/snacks
```

### Build
```bash
go build
# or
go install
```

### Run

##### Loris

For a full list of options, defaults and what they do
```bash
./snacks loris -h
```

The below command will
- send 1000000 arbitrary bytes to hold the connection open (not including HTTP POST headers)
- wait 10ms between sending each segment of bytes
- send 7 bytes in every segment (including the initial HTTP POST headers)
- set the path to `/boom` in the HTTP POST request
- attempt to open port `8888` on the target (defaults to `localhost`)
- enable trace logging
```bash
./snacks loris -size 1000000 -sd 10 -sb 7 -path /boom -port 8888 -vv
```

At the moment the Slow Loris implementation is geared towards HTTP POST requests with application/json.
