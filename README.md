# Go Expert Multithreading Exercise

## To test
- Clone this repo
- Run `main.go` file passing CEP information as argument
- CEP must be in following format: `12345-678`
- Example:
```shell
go run main.go 00000-000
```
- `main.go` file has some variables you may change in order to either simulate higher delays in provider responses or increase timeout delay:
  - `responseTimeoutDelaySeconds`: how many seconds application will wait before printing timeout error (default: 1)
  - `apiCEPResponseDelaySeconds`: delay to send results to channel after receiving APICEP response (default: 0)
  - `viaCEPResponseDelaySeconds`: delay to send results to channel after receiving VIACEP response (default: 0)

### Have fun ;-)
*Made with S2 by Henrique G Mendes*
