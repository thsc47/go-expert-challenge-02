package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response interface{}

type ApiResponse struct {
	Api string `json:"Api"`
}

type ViaCEP struct {
	ApiResponse
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApi struct {
	ApiResponse
	Cep        string `json:"cep"`
	Estado     string `json:"state"`
	Uf         string `json:"city"`
	Bairro     string `json:"neighborhood"`
	Logradouro string `json:"street"`
}

func main() {
	channelViaCEP := make(chan ViaCEP)
	channelBrasilApi := make(chan BrasilApi)

	go GetViaCEP(channelViaCEP)
	go GetBrasilApi(channelBrasilApi)

	select {
	case res := <-channelViaCEP:
		fmt.Println(res)

	case res := <-channelBrasilApi:
		fmt.Println(res)

	case <-time.After(time.Second):
		fmt.Println("TimeOut")

	}

}

func GetViaCEP(chApi chan ViaCEP) {
	var cep ViaCEP
	err := getCep("https://viacep.com.br/ws/01153000/json/", &cep)
	if err != nil {
		panic(err)
	}
	cep.Api = "ViaCEP"
	chApi <- cep
}

func GetBrasilApi(chVia chan BrasilApi) {
	var cep BrasilApi
	err := getCep("https://brasilapi.com.br/api/cep/v1/01153000", &cep)
	if err != nil {
		panic(err)
	}
	cep.Api = "BrasilApi"
	chVia <- cep
}

func getCep(uri string, res Response) error {
	req, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil
}
