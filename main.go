package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/henriquegmendes/go-expert-multithreading-exercise/client"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	// YOU MAY CHANGE THESE PARAMS TO TEST DIFFERENT RESULTS IN CEP RESULT RESPONSE
	responseTimeoutDelaySeconds = 1
	apiCEPResponseDelaySeconds  = 0
	viaCEPResponseDelaySeconds  = 0
)

const (
	cepRegexPattern   = "^[0-9]{5}-[0-9]{3}$"
	apiCEPUrlTemplate = "https://cdn.apicep.com/file/apicep/%s.json"
	viaCEPUrlTemplate = "http://viacep.com.br/ws/%s/json"
	resultTemplate    = "response received by provider %s. Result: %v"
	timeoutMessage    = "[error] provider cep response timed out"
)

// See README instructions before running this code
func main() {
	cep, err := getAndValidateCepArg()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()
	cepClient := client.NewClient(http.Client{}, apiCEPUrlTemplate, viaCEPUrlTemplate)
	apiCEPChan := make(chan client.Response[client.ApiCEPResponse])
	viaCEPChan := make(chan client.Response[client.ViaCEPResponse])

	go getApiCEPResults(ctx, cepClient, cep, apiCEPChan)
	go getViaCEPResults(ctx, cepClient, cep, viaCEPChan)

	result := ""
	select {
	case response := <-apiCEPChan:
		result = getResultString[client.ApiCEPResponse](response.ProviderName, response.Result)
	case response := <-viaCEPChan:
		result = getResultString[client.ViaCEPResponse](response.ProviderName, response.Result)
	case <-time.After(time.Second * responseTimeoutDelaySeconds):
		result = timeoutMessage
	}

	fmt.Println(result)
}

func getAndValidateCepArg() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return "", errors.New("[error] cep argument not informed")
	}

	cep := args[1]
	matchString, err := regexp.MatchString(cepRegexPattern, cep)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[error] could not validate cep %s: %s", cep, err.Error()))
	}
	if !matchString {
		return "", errors.New(fmt.Sprintf("[error] cep %s should be in format '12345-678'", cep))
	}

	return cep, nil
}

func getApiCEPResults(ctx context.Context, cepClient client.Client, cep string, apiCEPChan chan<- client.Response[client.ApiCEPResponse]) {
	time.Sleep(time.Second * apiCEPResponseDelaySeconds)

	result, err := cepClient.GetCepFromApiCEP(ctx, cep)
	if err != nil {
		fmt.Println(fmt.Sprintf("[error] could not get results from api cep: %s", err.Error()))
		return
	}

	apiCEPChan <- *result
}
func getViaCEPResults(ctx context.Context, cepClient client.Client, cep string, viaCEPChan chan<- client.Response[client.ViaCEPResponse]) {
	time.Sleep(time.Second * viaCEPResponseDelaySeconds)

	result, err := cepClient.GetCepFromViaCEP(ctx, cep)
	if err != nil {
		fmt.Println(fmt.Sprintf("[error] could not get results from via cep: %s", err.Error()))
		return
	}

	viaCEPChan <- *result
}

func getResultString[T any](providerName string, result T) string {
	return fmt.Sprintf(resultTemplate, providerName, result)
}
