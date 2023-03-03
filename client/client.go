package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	providerApiCEP = "apicep"
	providerViaCEP = "viacep"
)

type Response[T any] struct {
	ProviderName string
	Result       T
}

type ApiCEPResponse struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	OK         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCEPResponse struct {
	Cep          string `json:"cep"`
	Address      string `json:"logradouro"`
	Complement   string `json:"complemento"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	IBGE         string `json:"ibge"`
	GIA          string `json:"gia"`
	DDD          string `json:"ddd"`
	SIAFI        string `json:"siafi"`
}

type Client interface {
	GetCepFromApiCEP(ctx context.Context, cep string) (*Response[ApiCEPResponse], error)
	GetCepFromViaCEP(ctx context.Context, cep string) (*Response[ViaCEPResponse], error)
}

type cepClient struct {
	client                http.Client
	apiCEPBaseURLTemplate string
	viaCEPBaseURLTemplate string
}

func (c *cepClient) GetCepFromApiCEP(ctx context.Context, cep string) (*Response[ApiCEPResponse], error) {
	responseBytes, err := c.doRequest(ctx, fmt.Sprintf(c.apiCEPBaseURLTemplate, cep), http.MethodGet)
	if err != nil {
		return nil, err
	}

	var response ApiCEPResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	return &Response[ApiCEPResponse]{
		ProviderName: providerApiCEP,
		Result:       response,
	}, nil
}

func (c *cepClient) GetCepFromViaCEP(ctx context.Context, cep string) (*Response[ViaCEPResponse], error) {
	responseBytes, err := c.doRequest(ctx, fmt.Sprintf(c.viaCEPBaseURLTemplate, cep), http.MethodGet)
	if err != nil {
		return nil, err
	}

	var response ViaCEPResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	return &Response[ViaCEPResponse]{
		ProviderName: providerViaCEP,
		Result:       response,
	}, nil
}

func (c *cepClient) doRequest(ctx context.Context, url string, method string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if c.checkClientWithErrorStatusCode(response.StatusCode) {
		return nil, c.parseClientErrorResponse(response.StatusCode, responseBytes)
	}

	return responseBytes, nil
}

func (c *cepClient) checkClientWithErrorStatusCode(statusCode int) bool {
	return statusCode >= 400 && statusCode <= 599
}

func (c *cepClient) parseClientErrorResponse(statusCode int, body []byte) error {
	return errors.New(fmt.Sprintf("client error - status: %v - body - '%s'", statusCode, string(body)))
}

func NewClient(client http.Client, apiCEPBaseURLTemplate string, viaCEPBaseURLTemplate string) Client {
	return &cepClient{
		client:                client,
		apiCEPBaseURLTemplate: apiCEPBaseURLTemplate,
		viaCEPBaseURLTemplate: viaCEPBaseURLTemplate,
	}
}
