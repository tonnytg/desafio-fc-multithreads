package main

import (
	"context"
	"fmt"
	"github.com/tonnytg/desafio-fc-multithreads/pkg/webclient"
	"log"
	"sync"
	"time"
)

func MakeRequest(cep string) error {

	ctx, cancel := context.WithCancel(context.Background())

	firstAnswer := make(chan bool)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup, cep string) {

		method := "GET"
		url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
		log.Println("Request 1 - URL", url)

		bodyByte, err := webclient.Request(ctx, method, url, nil)
		if err != nil {
			log.Println(err)
		}

		log.Println("Request 1", string(bodyByte))
		firstAnswer <- true

		wg.Done()

	}(&wg, cep)

	go func(wg *sync.WaitGroup, cep string) {

		method := "GET"
		url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
		log.Println("Request 2 - URL", url)

		bodyByte, err := webclient.Request(ctx, method, url, nil)
		if err != nil {
			log.Println(err)
		}

		log.Println("Request 2", string(bodyByte))
		firstAnswer <- true

		wg.Done()
	}(&wg, cep)

	go func() {
		log.Println("Check if needs cancel")

		for {
			if <-firstAnswer {
				cancel()
			}
			log.Println("End Check if needs cancel")
		}

	}()

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

			log.Println("CEP", cep)
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
