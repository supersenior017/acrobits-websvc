# acrobits-websvc

[Acrobits Web Services](https://doc.acrobits.net/api/client/index.html). 
You must overwrite functions in `main.go` file.

[balance checker](https://doc.acrobits.net/api/client/balance_checker.html)
```go
func getBalance(
	ctx context.Context,
	account websvc.Account,
) (websvc.Balance, error) {
	return websvc.Balance{}, fmt.Errorf("NotImplemented")
}
```  

[rate checker](https://doc.acrobits.net/api/client/rate_checker.html)
```go
func getRate(
	ctx context.Context,
	rate websvc.RateParams,
) (websvc.Rate, error) {
	return websvc.Rate{}, fmt.Errorf("NotImplemented")
}
```

## Installation
```sh
$ make
$ make install
```
or
```sh
$ docker build -t acrobits-websvc "."
```

## Usage
```text
Usage of acrobits-websvc:
  -V	Print version and exit
  -c value
    	Path to configuration file
  -h	Print help and exit
```

## Example

To run localhost:
```sh
$ acrobits-websvc
```

To run with config:
```sh
$ acrobits-websvc -c /usr/local/etc/acrobits-websvc.json
```

To run in docker:
```sh
$ docker run -d -p 8080:8080 acrobits-websvc
```

To run in docker with config and env file (docker-compose.yaml):
```sh
$ docker-compose up -d
```
