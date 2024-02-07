package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type CustomerList struct {
	queue []string
	mux   sync.Mutex
}

var customerList = CustomerList{}

func CustomerService(customerName string) (string, error) {
	customerList.mux.Lock()
	defer customerList.mux.Unlock()
	if customerName == "" {
		c := customerList.queue[0]
		customerList.queue = append(customerList.queue[:0], customerList.queue[1:]...)
		return c, nil
	} else {
		if len(customerList.queue) >= 2 {
			return "", errors.New("OUT OF SERVICE DUE TO FULL OF CAPACITY")
		}
		customerList.queue = append(customerList.queue, customerName)
		return customerName, nil
	}
}
func DoBarber(ch chan string) {
	for {
		if len(customerList.queue) == 0 {
			fmt.Println("Barber is sleeping...")
			select {
			case msg := <-ch:
				fmt.Println("Barber walking..." + msg)
			}
		}
		customerName, _ := CustomerService("")
		fmt.Println("Doing Barber with the " + customerName)
		time.Sleep(5 * time.Second)
	}
}
func main() {
	channel := make(chan string)
	go DoBarber(channel)
	var customerName string
	fmt.Println("Welcome to My Hairdresser...")
	for {
		_, err := fmt.Scanln(&customerName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(customerList.queue) == 0 {
			select {
			case channel <- "Wake up":
			case <-time.After(1 * time.Second):
			}
		}
		_, err = CustomerService(customerName)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
