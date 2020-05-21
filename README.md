# acrobits-balance

Acrobits balance checker web service. You must overwrite 
`getBalance(username, password string) (float64, error)` function in `main.go` 
file.

[Doc](https://doc.acrobits.net/api/client/balance_checker.html)

## Installation
```sh
$ make
$ make install
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

To default run:
```sh
$ acrobits-balance -c ~/acrobits-balance.json
```
