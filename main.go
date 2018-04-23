//
// Реализовать чтение с нескольких каналов
// ps: программа корректно работает в режиме гонки -race
//

package main

import (
	"fmt"
	"sync"
)

// Канал с которого бум читать, принимает на вход числа
// Таких каналов можно создать множество
// Например они что-то вычисляют паралелньно и отправляют в канал
func SomeChannel(nums ...int) <-chan int {
	outChannel := make(chan int)
	go func() {
		for _, num := range nums {
			outChannel <- num * 2
		}
		close(outChannel)
	}()
	return outChannel
}

// Функция чтения с каналов
// На вход получаем каналы, и отдает в один канал подсчитанные данные
func OneWorkerReader(channels ...<-chan int) <-chan int {
	var vg sync.WaitGroup
	vg.Add(len(channels))
	outChannel := make(chan int)
	outPutWorkers := func(in <-chan int) {
		for ch := range in {
			outChannel <- ch
		}
		vg.Done()
	}

	for _, ch := range channels {
		go outPutWorkers(ch)
	}

	go func() {
		vg.Wait()
		close(outChannel)
	}()

	return outChannel
}

func main() {
	// Создаем воркеров, которое что-то вычисляют и возвращают каналы
	oneChannel := SomeChannel(1, 2, 3, 4, 5)
	twoChannel := SomeChannel(100, 200)
	threeChannel := SomeChannel(1000, 2000)

	// Передаем в функцию каналы наших воркеров и печатаем результат из одного канала
	for n := range OneWorkerReader(oneChannel, twoChannel, threeChannel) {
		fmt.Println(n)
	}
}
