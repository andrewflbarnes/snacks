# snacks

A small collection of basically (poorly) implemented tools for pen testing.

This project is more of a learning experience in golang so will likely fall short of a lot of go best practices

### Build
```bash
go build
# or
go install
```

### Run
```bash
# help
./snacks loris -h
./snacks loris [-port <port>] [-host <host>] [-v[v]] [embed]
# e.g. the below will start an embedded server on port 8989, send a payload to it and log at trace level
./snacks loris -embed -port 8989 -vv
```

At the moment there is a single stubbed implementation for Slow Loris.

The program will generate a payload and send it to the server using a strategy which
determines the number of bytes to send per delay period and how long the delay period is.

To change the number of bytes sent and the delay between sends changes need to be made in
`main.go` and then rebuilt.
