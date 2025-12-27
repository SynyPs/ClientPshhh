package repository

import (
	"github.com/fhg/ClientPshhh/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(dbPath string) (*MessageRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.Message{}, &models.KeyStore{}, &models.Contact{})
	return &MessageRepository{db: db}, nil
}

func (r *MessageRepository) SaveMessage(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *MessageRepository) GetLastMessage(limit int) ([]models.Message, error) {
	var msgs []models.Message
	retult := r.db.Order("created_at desc").Limit(limit).Find(&msgs)
	return msgs, retult.Error
}

func (r *MessageRepository) GetContact() ([]models.Contact, error) {
	var contacts []models.Contact
	result := r.db.Order("updated_at desc").Find(&contacts)
	return contacts, result.Error
}

func (r *MessageRepository) SeedContacts() {
	// Список фейковых друзей
	dummies := []models.Contact{
		{PublicKey: "key_alice_123", Name: "Alice (Friend)", LastMsg: "Привет, как дела?"},
		{PublicKey: "key_bob_456", Name: "Bob (Work)", LastMsg: "Скинь отчет"},
		{PublicKey: "key_eva_789", Name: "Eva (Spam)", LastMsg: "Купи гараж"},
	}

	for _, c := range dummies {
		// Ищем по PublicKey. Если нет — создаем.
		// &c в первом аргументе обновится данными из БД или созданными данными
		r.db.Where(models.Contact{PublicKey: c.PublicKey}).FirstOrCreate(&c)
	}
}

func (r *MessageRepository) FindOrCreateContact(pubKey string, defaultName string) error {
	var c models.Contact
	err := r.db.FirstOrCreate(&c, models.Contact{PublicKey: pubKey}).Error
	if err != nil {
		return err
	}
	if c.Name == "" {
		c.Name = defaultName
		r.db.Save(&c)
	}
	return nil
}

func (r *MessageRepository) UpdateLastMessage(pubKey string, msg string) {
	r.db.Model(&models.Contact{}).Where("public_key = ?", pubKey).Update("last_msg", msg)
}
