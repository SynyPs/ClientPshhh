package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/fhg/ClientPshhh/internal/crypto"
	"github.com/fhg/ClientPshhh/internal/gui"
	"github.com/fhg/ClientPshhh/internal/repository"
	"github.com/fhg/ClientPshhh/internal/service"
)

func main() {
	a := app.NewWithID("com.pshhh.fuc")
	a.Settings().SetTheme(theme.DarkTheme())

	repo, err := repository.NewMessageRepository("pshh.bin")
	if err != nil {
		log.Printf("Критическиая ошибка инициализации:", err)
		os.Exit(1)
	}

	repo.SeedContacts()

	keyManager, err := crypto.LoadOrGenerateKeys("keys.dat")
	if err != nil {
		log.Printf("Ошибка работы с ключами: %v", err)
		os.Exit(1)
	}

	myPubKey := keyManager.PublicBase64() // Нужно добавить этот метод в crypto (см. ниже)
	fmt.Println("\n========================================")
	fmt.Println("МОЙ ПУБЛИЧНЫЙ КЛЮЧ (Скопируй его):")
	fmt.Println(myPubKey)
	fmt.Println("========================================\n")

	svc := service.NewOrchestrator(repo, keyManager)

	mainWindow := gui.NewMainWindow(a, svc)

	mainWindow.ShowAndRun()
}
