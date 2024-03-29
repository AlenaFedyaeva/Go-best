package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type crawlResult struct {
	err error
	msg string
}

type crawler struct {
	muxVisit   sync.Mutex
	rwmuxDepth sync.RWMutex
	visited    map[string]string
	maxDepth   int
}

func (c *crawler) GetMaxDepth() int {
	c.rwmuxDepth.RLock()
	defer c.rwmuxDepth.RUnlock()

	return c.maxDepth
}
//task2: увеличить глубину поиска при приеме USR1
func (c *crawler) IncreaseMaxDepth(val int) int {
	c.rwmuxDepth.Lock()
	defer c.rwmuxDepth.Unlock()

	c.maxDepth=c.maxDepth+val
	return c.maxDepth
}


func newCrawler(maxDepth int) *crawler {
	return &crawler{
		visited:  make(map[string]string),
		maxDepth: maxDepth,
	}
}

// рекурсивно сканируем страницы
func (c *crawler) run(ctx context.Context, url string, results chan<- crawlResult, depth int) {
	// просто для того, чтобы успевать следить за выводом программы, можно убрать :)
	time.Sleep(2 * time.Second)

	// проверяем что контекст исполнения актуален
	select {
	case <-ctx.Done():
		//fmt.Println("ctx.done return")
		return

	default:
		// проверка глубины
		if depth >= c.GetMaxDepth() {
			//fmt.Println("depth return")
			return
		}

		page, err := parse(ctx,url)
		if err != nil {
			// ошибку отправляем в канал, а не обрабатываем на месте
			results <- crawlResult{
				err: errors.Wrapf(err, "parse page %s", url),
			}
			return
		}

		title := pageTitle(page)
		links := pageLinks(nil, page)

		fmt.Println("links",links)

		// блокировка требуется, т.к. мы модифицируем мапку в несколько горутин
		c.muxVisit.Lock()
		c.visited[url] = title
		c.muxVisit.Unlock()

		// отправляем результат в канал, не обрабатывая на месте
		results <- crawlResult{
			err: nil,
			msg: fmt.Sprintf("%s -> %s\n", url, title),
		}

		// рекурсивно ищем ссылки
		for link := range links {
			// если ссылка не найдена, то запускаем анализ по новой ссылке
			// fmt.Println("link",link)
			if c.checkVisited(link) {
				continue
			}

			go c.run(ctx, link, results, depth+1)
		}
	}
}

func (c *crawler) checkVisited(url string) bool {
	c.muxVisit.Lock()
	defer c.muxVisit.Unlock()

	_, ok := c.visited[url]
	return ok
}
