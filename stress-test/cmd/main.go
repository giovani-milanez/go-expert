package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"errors"
	"flag"
)

func main() {
	urlFlag := flag.String("url", "", "URL para teste de estresse, exemplo: https://www.google.com/")
	requestFlag := flag.Int("requests", 1, "Número total de requisições a serem feitas")
	concurrencyFlag := flag.Int("concurrency", 1, "Número de rotinas concorrentes")
	flag.Parse()

	url := *urlFlag
	requests := *requestFlag
	concurrency := *concurrencyFlag

	if url == "" {
		log.Fatal("A URL para teste de estresse deve ser informada com a flag -url")
	}
	if requests <= 1 || concurrency <= 1 {
		log.Fatal("O número de requisições e concorrência deve ser maior que 1")
	}

	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	startTime := time.Now()
	log.Printf("Iniciando teste de estresse na URL: %s com %d requisições e %d concorrência", url, requests, concurrency)
	
	errorsCounter := int32(0)
	total  := int32(0)
	counters := map[int]int32{}
	
	var lock sync.RWMutex
	var wg sync.WaitGroup
	reqPerThread := requests / concurrency
	extra := requests % concurrency

	wg.Add(concurrency + 1)

	worker := func(reqs int) {
		defer wg.Done()
		for range reqs {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					lock.Lock()
					errorsCounter++
					lock.Unlock()
					log.Printf("Erro ao criar requisição: %v", err)
					continue
				}
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					lock.Lock()
					errorsCounter++
					lock.Unlock()
					log.Printf("Erro ao fazer requisição: %v", err)
					continue
				}
				return
			}
			defer resp.Body.Close()

			lock.Lock()
			total++
			val, ok := counters[resp.StatusCode]
			if !ok {
				val = 0
			}
			counters[resp.StatusCode] = val + 1
			lock.Unlock()
		}
	}

	for range concurrency {
		go worker(reqPerThread)
	}
	go worker(extra)

	go func() {
		sig := <-sigs
		log.Printf("Recebido sinal %s, finalizando o teste de estresse aguarde...", sig)
		cancel()
	}()

	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Println("Tempo total de execução:", duration)
	for key, value := range counters {
		log.Printf("[HTTP %d %s]: %d requisições\n", key, http.StatusText(key), value)
	}
	log.Println("Total de requisições realizadas:", total)
	log.Println("Total de requisições com falhas:", errorsCounter)
	log.Println("Total de requisições nao realizadas:", int32(requests) - total - errorsCounter)
}