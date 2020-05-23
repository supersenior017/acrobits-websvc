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

To test run:
```sh
$ acrobits-balance
```

To config run:
```sh
$ acrobits-balance -c /usr/local/etc/acrobits-balance.json
```

To docker run:
```sh
$ docker run -d -p 8080:8080 acrobits-balance
```

To docker config run (in folder):
```sh
$ docker-compose up -d
```
