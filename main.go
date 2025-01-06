package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	brazilAPIURL = "https://brasilapi.com.br/api/cep/v1/"
	viaCepURL    = "http://viacep.com.br/ws/"
	timeout      = 1 * time.Second
)

type ApiResponse struct {
	Source string
	Data   string
}

func fetchFromBrazilAPI(cep string, ch chan<- ApiResponse) {
	resp, err := http.Get(brazilAPIURL + cep)
	if err != nil {
		ch <- ApiResponse{Source: "BrasilAPI", Data: fmt.Sprintf("Erro: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- ApiResponse{Source: "BrasilAPI", Data: fmt.Sprintf("Erro na resposta: %s", resp.Status)}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- ApiResponse{Source: "BrasilAPI", Data: fmt.Sprintf("Erro ao ler resposta: %v", err)}
		return
	}

	ch <- ApiResponse{Source: "BrasilAPI", Data: string(body)}
}

func fetchFromViaCep(cep string, ch chan<- ApiResponse) {
	resp, err := http.Get(viaCepURL + cep + "/json/")
	if err != nil {
		ch <- ApiResponse{Source: "ViaCep", Data: fmt.Sprintf("Erro: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- ApiResponse{Source: "ViaCep", Data: fmt.Sprintf("Erro na resposta: %s", resp.Status)}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- ApiResponse{Source: "ViaCep", Data: fmt.Sprintf("Erro ao ler resposta: %v", err)}
		return
	}

	ch <- ApiResponse{Source: "ViaCep", Data: string(body)}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scan(&cep)

	ch := make(chan ApiResponse, 2)

	go fetchFromBrazilAPI(cep, ch)
	go fetchFromViaCep(cep, ch)

	select {
	case response := <-ch:
		fmt.Printf("Resposta recebida da %s: %s\n", response.Source, response.Data)
	case response := <-ch:
		fmt.Printf("Resposta recebida da %s: %s\n", response.Source, response.Data)
	case <-time.After(timeout):
		fmt.Println("Erro: Timeout de 1 segundo atingido.")
	}
}
