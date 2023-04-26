/*
Напишите код, реализующий пайплайн, работающий с целыми числами и состоящий из следующих стадий:
    //Стадия фильтрации отрицательных чисел (не пропускать отрицательные числа).
    //Стадия фильтрации чисел, не кратных 3 (не пропускать такие числа), исключая также и 0.
    //Стадия буферизации данных в кольцевом буфере с интерфейсом, соответствующим тому,
		который был дан в качестве задания в 19 модуле. В этой стадии предусмотреть опустошение
		буфера (и соответственно, передачу этих данных, если они есть, дальше) с определённым
		интервалом во времени. Значения размера буфера и этого интервала времени сделать
		настраиваемыми (как мы делали: через константы или глобальные переменные).

Написать источник данных для конвейера. Непосредственным источником данных должна быть консоль.

Также написать код потребителя данных конвейера.
Данные от конвейера можно направить снова в консоль построчно, сопроводив их каким-нибудь
поясняющим текстом, например: «Получены данные …».

При написании источника данных подумайте о фильтрации нечисловых данных,
которые можно ввести через консоль. Как и где их фильтровать, решайте сами.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
var errLog = log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime)

// объявление структуры кольц буфера
type RingIntBuffer struct {
	array []int
	pos   int
	size  int
	m     sync.Mutex //блокировка сущности во время записи
}

// сканирование ввода через консоль, также фильтрация не числовых данных
func readInput(input chan<- int) {
	infoLog.Println("Starting reading users input")
	for {
		var u int
		_, err := fmt.Scanf("%d \n", &u)
		if err != nil {
			fmt.Println("This isn't a number")
			errLog.Println("Wrong user input!")
		} else {
			input <- u
		}
	}
}

// создание динамического буффера
func NewRingIntBuffer(size int) *RingIntBuffer {
	infoLog.Println("Creating dynamic buffer")
	return &RingIntBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

// функция свдига элемента
func (r *RingIntBuffer) Push(el int) {
	infoLog.Println("Staring pushing function")
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else { //если позиция не совпала
		r.pos++
		r.array[r.pos] = el
	}
}

// функция получения элемента
func (r *RingIntBuffer) Get() []int {
	infoLog.Println("Starting Get function")
	if r.pos <= 0 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.array[:r.pos]

	r.pos = 0
	return output
}

// функция фильтрации отрицательных чисел
func removeNegatives(curntChn <-chan int, nxtChn chan<- int) {
	infoLog.Println("Starting negative numbers filtering function")
	for n := range curntChn {
		if n >= 0 {
			nxtChn <- n
		}
	}
}

// функция фильтрации чисел не кратных 3, исключая 0 также
func notDivToThree(curntChn <-chan int, nxtChn chan<- int) {
	infoLog.Println("Staring function remove 0 numbers and not multiples of 3")
	for n := range curntChn {
		if n%3 != 0 {
			nxtChn <- n
		}
	}
}

// функция записи в буффер значений
func writeToBuffer(curntChn <-chan int, r *RingIntBuffer) {
	infoLog.Println("Start function writting data in buffer")
	for n := range curntChn {
		r.Push(n)
	}
}

// функция показа данных из буфера  в консоль
func writeToConsole(r *RingIntBuffer, t *time.Ticker) {
	infoLog.Println("Starting function to show data in console")
	for range t.C {
		buffer := r.Get()
		if len(buffer) > 0 {
			fmt.Println("The buffer is", buffer)
		}
	}
}

func main() {
	//var l = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	//Create channel for input numbers from console
	infoLog.Print("Create channel for input numbers from console")
	input := make(chan int)
	go readInput(input)

	// declaring next chanel for filter and remove all negative numbers
	infoLog.Print("Create chanel for filter and remove all negative numbers")
	rmvNegat := make(chan int)
	go removeNegatives(input, rmvNegat)

	//declaring next chanel for filter and remove all numbers that multiply of 3, including 0
	infoLog.Print("Create  chanel for filter numbers")
	notDivTo3 := make(chan int)
	go notDivToThree(rmvNegat, notDivTo3)

	size := 4
	rng := NewRingIntBuffer(size)
	infoLog.Print("Create chanel showing data in console")
	go writeToBuffer(notDivTo3, rng)

	delay := 5
	ticker := time.NewTicker(time.Second * time.Duration(delay))
	go writeToConsole(rng, ticker)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case sig := <-c:
		infoLog.Println("Exit programm")
		fmt.Printf("Got %s signal. Aborting .... \n", sig)
		os.Exit(0)
	}
}
