package main

import (
	"context"
	"fmt"
	"github.com/tonnytg/desafio-fc-multithreads/pkg/webclient"
	"log"
	"sync"
	"time"
)

//Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
//As duas requisições serão feitas simultaneamente para as seguintes APIs:
//https://brasilapi.com.br/api/cep/v1/01153000 + cep
//http://viacep.com.br/ws/" + cep + "/json/
//Os requisitos para este desafio são:
//- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
//- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
//- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

func MakeRequest(cep string) error {
	// Criar o contexto com cancelamento
	ctx, cancel := context.WithCancel(context.Background())

	// Canal para receber a primeira resposta
	firstResponse := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Primeira gorrotina
	go func() {
		defer wg.Done()

		method := "GET"
		url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

		bodyByte, err := webclient.Request(ctx, method, url, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Request 1:", string(bodyByte))

		// Sinalizar que a primeira resposta foi recebida
		firstResponse <- struct{}{}
	}()

	// Segunda gorrotina
	go func() {
		defer wg.Done()

		method := "GET"
		url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

		bodyByte, err := webclient.Request(ctx, method, url, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Request 2:", string(bodyByte))

		// Se a primeira resposta já foi recebida, cancelar esta gorrotina
		select {
		case <-firstResponse:
			cancel()
		default:
			// Nada a fazer, continuar esperando
		}
	}()

	// Esperar pela finalização das gorrotinas
	wg.Wait()

	return nil
}

func main() {
	log.Println("Start Get CEP")

	wg := sync.WaitGroup{}
	chanRequestCep := make(chan string)

	ceps := []string{"01153000", "05541000", "08032330", "01550020", "03312006", "01415900", "05639050"}

	go func(wg *sync.WaitGroup) {

		wg.Add(1)

		for _, cep := range ceps {
			chanRequestCep <- cep
			time.Sleep(5 * time.Second)
		}
		close(chanRequestCep)
		wg.Done()
	}(&wg)

	func(wg *sync.WaitGroup) {

		wg.Add(1)

		for cep := range chanRequestCep {
			err := MakeRequest(cep)
			if err != nil {
				log.Println(err)
			}
		}

		wg.Done()
	}(&wg)
	log.Println("End Get CEP")
	wg.Wait()
}