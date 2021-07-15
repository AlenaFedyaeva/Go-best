package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// максимально допустимое число ошибок при парсинге
	errorsLimit = 100000

	// число результатов, которые хотим получить
	resultsLimit = 10000

	// На сколько увеличить глубину поиска, при приеме USR1
	increaseDepth=2
)

var (
	// адрес в интернете (например, https://en.wikipedia.org/wiki/Lionel_Messi)
	url string

	// насколько глубоко нам надо смотреть (например, 10)
	depthLimit int
)

// Как вы помните, функция инициализации стартует первой
func init() {
	// задаём и парсим флаги
	flag.StringVar(&url, "url", "https://en.wikipedia.org/wiki/Main_Page", "url address")
	flag.IntVar(&depthLimit, "depth", 3, "max depth for run")
	flag.Parse()

	// Проверяем обязательное условие
	if url == "" {
		log.Print("no url set by flag. Used default value")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	started := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	go watchSignals(cancel)
	defer cancel()

	crawler := newCrawler(depthLimit)

	// создаём канал для результатов
	results := make(chan crawlResult)


	// создаём канал для увеличения глубины
	// depth := make(chan int)
	//запускаем горутину для чтения сигнала USR1
	go watchUSR1(crawler)

	// запускаем горутину для чтения из каналов
	done := watchCrawler(ctx, results, errorsLimit, resultsLimit)

	// запуск основной логики
	// внутри есть рекурсивные запуски анализа в других горутинах
	crawler.run(ctx, url, results, 0)

	// ждём завершения работы чтения в своей горутине
	<-done

	log.Println(time.Since(started))
}

//task2: увеличить глубину поиска на 2
func watchUSR1(c *crawler) {
	osSigUser1 := make(chan os.Signal)
	signal.Notify(osSigUser1, syscall.SIGUSR1)
	
	for _ = range osSigUser1 {
		c.IncreaseMaxDepth(increaseDepth)
	}
}

// ловим сигналы выключения
func watchSignals(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal, 2)

	signal.Notify(osSignalChan,
		syscall.SIGINT,
		syscall.SIGTERM)
	

	
	sig := <-osSignalChan
	log.Printf("got signal %q", sig.String())

	// если сигнал получен, отменяем контекст работы
	cancel()
}

func watchCrawler(ctx context.Context, results <-chan crawlResult, maxErrors, maxResults int) chan struct{} {
	readersDone := make(chan struct{})

	go func() {
		defer close(readersDone)
		for {
			select {
			case <-ctx.Done():
				return

			case result := <-results:
				if result.err != nil {
					maxErrors--
					if maxErrors <= 0 {
						log.Println("max errors exceeded")
						return
					}
					continue
				}

				log.Printf("crawling result: %v", result.msg)
				maxResults--
				if maxResults <= 0 {
					log.Println("got max results")
					return
				}
			}
		}
	}()

	return readersDone
}
