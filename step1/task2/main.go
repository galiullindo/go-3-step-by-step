package main

import "fmt"

type Order struct {
	OrderNumber  int
	CustomerName string
	OrderAmount  float64
}

type OrderLogger struct{}

func NewOrderLogger() *OrderLogger {
	return &OrderLogger{}
}

func (l *OrderLogger) AddOrder(o Order) {
	fmt.Printf("Добавлен заказ #%d, Имя клиента: %s, Сумма заказа: $%.2f\n", o.OrderNumber, o.CustomerName, o.OrderAmount)
}
