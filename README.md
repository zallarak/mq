# mq - get stock market quotes

0.0.1

### Summary

mq is a command line tool for getting market quotes

```sh
$ mq -s GOOG,AMZN,BTC
         |Price     |% Change today
GOOG     |718.07    |+1.96%
AMZN     |706.36    |+1.38%
BTCUSD=X |445.00    |+0.23%
```

### Installation

mq requires the Go Tools (v1.6+ tested, but others are probably fine). Instructions are [here](https://golang.org/doc/install) for installation.

`make` builds the binary and `make install` places it on your path.

```sh
$ git clone git@github.com:zallarak/mq.git
$ cd mq
$ make && make install
```

### Usage

```sh
mq -s <COMMA SEPERATED SYMBOLS> -f <FILE OF NEWLINE SEPERATED SYMBOLS>
```

### Other stuff

Contributions are welcome.

Future plans for development:
* Better error handling
* Add other market data providers (currently, Yahoo Finance is queried)
