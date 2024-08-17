package main

import (
	"wb/internal/app"
)

func main() {
	//Создание приложения
	App, err := app.New()
	if err != nil {
		panic(err)
	}

	//Запуск приложения
	err = App.Run()
	if err != nil {
		panic(err)
	}
}

//TODO: Использование кэширования запросов
//TODO: Использование Goroutines
//TODO: Профилирование кода
