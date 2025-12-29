package models

import "gorm.io/gorm"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	PublicKey string `gorm:"uniqueIndex"`
	Name      string
	IsMe      bool
}

type KeyStore struct {
	ID              uint
	PrivateKeyEd    []byte // Для подписки
	PrivateKeyCurve []byte // Для расшифровки
}

type Contact struct {
	gorm.Model
	PublicKey string `gorm:"uniqueIndex"`
	Name      string
	LastMsg   string
}

type Message struct {
	gorm.Model
	SenderKey   string // Public key
	ReceiverKey string // Public key
	Content     []byte // Зашифрованный контент
	Decrypted   string `gorm:"-"` // Поле для UI
	Nonce       []byte // для box
	Signature   []byte // Подпись
}

type UIMessage struct {
	Text string
	IsMe bool
}
