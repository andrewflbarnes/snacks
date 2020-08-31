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
make
# or
make install
```

### Run

##### Judy

Judy launches a RUDY (r-u-dead-yet) attack using JSON as the content-type over more typical
MIME types found in other implementations.

For a full list of options, defaults and what they do
```bash
./snacks judy -h
```

The below command will
- send 1000000 arbitrary bytes to hold the connection open (not including HTTP POST headers)
- wait 10ms between sending each segment of bytes
- send 7 bytes in every segment (excluding the initial HTTP POST headers)
- set the path to `/boom` in the HTTP POST request and send to `localhost:8888`
- enable trace logging
```bash
./snacks judy -size 1000000 -sd 10ms -sb 7 -vv localhost:8888/boom
```


##### Loris

Loris launches a Slow Loris attack by repeatedly sending a duplicate HTTP header every fixed time period.
At the moment this will always send the whole header, a `-sb`/SendBytes option is not currently available.

For a full list of options, defaults and what they do
```bash
./snacks loris -h
```

The below command will
- send 1000000 repeat instances of the header (not most application servers limit the max. number of HTTP headers)
- wait 1s between sending each header
- set the header to be set to `x-slow: loris`
- set the path to `/boom` in the HTTP POST request and send to `localhost:8888`
- enable trace logging
```bash
./snacks loris -size 1000000 -sd 1s -head "x-slow: loris" -vv localhost:8888/boom
```

### Useful options

The `loris` attack will not use most of these options as it never completes sending the HTTP headers

##### Content type

If a specific content-type header is required use the `-type` flag. e.g.
```bash
./snacks -type application/x-www-form-urlencoded ...
```

For supported content types this will set a default payload prefix which may be overridden with `-prefix`

To override a payload prefix use the `-prefix` flag. e.g. for a default JSON content-type (which would otherwise
default to using `{"a":"` as the payload prefix)
```bash
./snacks -prefix '{"payload":"' ...
```

For custom content-types and prefix apyloads specify both `-type` and `-prefix` e.g.
```bash
./snacks -type application/xml -prefix '<payload>' ...
./snacks -type application/vnd.my.custom.type -prefix '1|string|payload|' ...
```

##### Authorization

If authorization is required use either `-basic` or `-bearer` flags e.g.
```bash
./snacks -basic tomcat:tomcat ...
./snacks -bearer 0123456789ABCDEF ...
```