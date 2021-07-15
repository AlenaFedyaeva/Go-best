package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//sigint -2  /usr1 -10 / sigterm -15 / sighup -1  /sigquit -3
//ps -aux | grep mainMe
//kill -1 18164
func main() {
	log.Print("start")

	fmt.Println("my pid ", os.Getpid())
	signalChan := make(chan os.Signal, 4)

	signal.Notify(signalChan, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT) // регистрируем каналы для получения нотификации указанных сигналов

	sigUser1 := make(chan os.Signal)
	signal.Notify(sigUser1, syscall.SIGUSR1)

	select {
	case sig := <-signalChan:
		fmt.Println(sig.String())
	case sig := <-sigUser1:
		fmt.Println("user",sig.String())
	// default:
		// fmt.Println("nothing available")
	}


	// sig := <-signalChan // ожидаем получения сигнала из канала
	// log.Printf("got %s signal", sig.String())
	log.Print("end")
}
