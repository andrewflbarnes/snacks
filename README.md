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
./snacks
```

At the moment there is a single stubbed implementation for Slow Loris.

The program will generate a payload and send it to the server using a strategy which
determines the number of bytes to send per delay period and how long the delay period is.

By default an embedded server will be started to connect to. There are no runtime confiurable
options, changes need to be made in `main.go` and then rebuilt.
