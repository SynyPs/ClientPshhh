package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/fhg/ClientPshhh/internal/crypto"
	"github.com/fhg/ClientPshhh/internal/models"
	"github.com/fhg/ClientPshhh/internal/repository"
)

type Orchestrator struct {
	repo *repository.MessageRepository
	keys *crypto.CryptoManager
	conn *websocket.Conn

	writeMu sync.Mutex

	onMessage func(string)
}

func NewOrchestrator(repo *repository.MessageRepository, keys *crypto.CryptoManager) *Orchestrator {
	return &Orchestrator{
		repo: repo,
		keys: keys,
	}
}

func (o *Orchestrator) SetOnMessageReceived(f func(string)) {
	o.onMessage = f
}

func (o *Orchestrator) Connect(url string) error {
	// if o.conn != nil {
	// 	log.Println("Закрываю старое соединение...")
	// 	o.conn.Close(websocket.StatusNormalClosure, "Reconnecting")
	// }

	// ctx := context.Background()

	// // === ДОБАВЛЯЕМ ОПЦИИ ДЛЯ NGROK ===
	// opts := &websocket.DialOptions{
	// 	HTTPHeader: http.Header{},
	// }
	// // Этот заголовок заставляет Ngrok пропустить страницу с предупреждением
	// opts.HTTPHeader.Set("ngrok-skip-browser-warning", "true")

	// // Передаем opts вместо nil
	// c, _, err := websocket.Dial(ctx, url, opts)
	// if err != nil {
	// 	return err
	// }
	// o.conn = c

	// go o.readLoop()

	// return nil
	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}
	o.conn = c

	go o.readLoop()

	return nil
}

func (o *Orchestrator) readLoop() {
	defer o.conn.Close(websocket.StatusInternalError, "closing")

	myKey := o.keys.PublicBase64()

	for {
		var msgData models.Message
		err := wsjson.Read(context.Background(), o.conn, &msgData)
		if err != nil {
			log.Println("Разрыв соединения:", err)
			return
		}

		// decrypted := o.keys.Decrypt(msgData.Content, msgData.SenderKey)

		if msgData.SenderKey == myKey {
			// Это мое сообщение, игнорируем
			continue
		}

		if err := o.repo.SaveMessage(&msgData); err != nil {
			log.Println("Ошибка сохранения входящего сообщения в БД:", err)
			continue
		}

		text := string(msgData.Content)

		err = o.repo.FindOrCreateContact(msgData.SenderKey, "User "+shortKey(msgData.SenderKey))
		if err != nil {
			log.Println("Error create contact:", err)
		}

		o.repo.UpdateLastMessage(msgData.SenderKey, text)

		if o.onMessage != nil {
			o.onMessage(text)
		}
	}
}

func shortKey(k string) string {
	if len(k) > 5 {
		return k[:]
	}
	return k
}

func (o *Orchestrator) SendMessage(text string, receiver string) {
	if o.conn == nil {
		return
	}

	myKey := o.keys.PublicBase64()

	msg := models.Message{
		SenderKey:   myKey,
		ReceiverKey: receiver,
		Content:     []byte(text),
	}

	fmt.Println(msg)

	if err := o.repo.SaveMessage(&msg); err != nil {
		log.Println("Ошибка сохранения сообщения в БД:", err)
		return
	}

	o.writeMu.Lock()
	defer o.writeMu.Unlock()

	err := wsjson.Write(context.Background(), o.conn, msg)
	if err != nil {
		log.Println("Ошибка отправки:", err)
	}
}

func (o *Orchestrator) GetContacts() ([]models.Contact, error) {
	return o.repo.GetContact()
}

func (o *Orchestrator) AddContact(pubKey, name string) error {
	return o.repo.FindOrCreateContact(pubKey, name)
}

func (o *Orchestrator) LoadChatHistory(contactKey string) ([]models.UIMessage, error) {
	myKey := o.keys.PublicBase64()
	dbMsgs, err := o.repo.GetHistory(contactKey, myKey)
	if err != nil {
		return nil, err
	}

	var uiMsgs []models.UIMessage
	for _, m := range dbMsgs {
		uiMsgs = append(uiMsgs, models.UIMessage{
			Text: string(m.Content),
			IsMe: m.SenderKey != contactKey,
		})
	}
	return uiMsgs, nil
}
