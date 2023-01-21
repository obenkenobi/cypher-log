package app

import "fmt"

type App struct {
}

func (a App) Start() {
	fmt.Println("Hello world!!!")
}

func NewApp() *App {
	return &App{}
}
