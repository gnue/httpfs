# tinyweb

tiny http static server (dir, zip, git)

## Installation

```sh
$ go get github.com/gnue/httpfs/examples/tinyweb
```

## Usage

```
tinyweb [OPTIONS] [dir...]

Application Options:
  -H, --host=      host (default: localhost)
  -p, --port=      port (default: 3000)
  -b, --branch=    git branch (default: master)
      --index=     directory index (default: index.html)
      --autoindex  directory auto index
      --cert=      TLS cert file
      --key=       TLS key file

Help Options:
  -h, --help       Show this help message

Arguments:
  dir:             directory or zip
```

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

