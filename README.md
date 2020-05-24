# acrobits-balance

Acrobits balance checker web service. You must overwrite 
```go
func getBalance(
	ctx context.Context,
	username string,
	password string,
) (float64, error) {
	return 0, fmt.Errorf("NotImplemented")
}
``` 
in `main.go` file.

[Doc](https://doc.acrobits.net/api/client/balance_checker.html)

## Installation
```sh
$ make
$ make install
```
or
```sh
$ docker build -t acrobits-balance "."
```

## Usage
```text
Usage of acrobits-balance:
  -V	Print version and exit
  -c value
    	Path to configuration file
  -h	Print help and exit
```

## Example

To run localhost:
```sh
$ acrobits-balance
```

To run with config:
```sh
$ acrobits-balance -c /usr/local/etc/acrobits-balance.json
```

To run in docker:
```sh
$ docker run -d -p 8080:8080 acrobits-balance
```

To run in docker with config (docker-compose.yaml):
```sh
$ docker-compose up -d
```
