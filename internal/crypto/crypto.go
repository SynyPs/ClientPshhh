package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/nacl/box"
)

// KeyManager хранит ключи в памяти
type CryptoManager struct {
	SingPublic  ed25519.PublicKey
	SingPrivate ed25519.PrivateKey

	// Ключи для шифрования (Encryption)
	// nacl/box требует указатель на массив [32]byte
	pubKey  *[32]byte
	privKey *[32]byte
}

// storageStruct используется для сериализации в JSON
type storageStruct struct {
	SignPublic  []byte `json:"sign_public"`
	SignPrivate []byte `json:"sign_private"`
	BoxPublic   []byte `json:"box_public"`
	BoxPrivate  []byte `json:"box_private"`
}

func (cm *CryptoManager) Encrypt(msg []byte, peerPubKey *[32]byte) ([]byte, []byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, nil, err
	}
	encrypted := box.Seal(nil, msg, &nonce, peerPubKey, cm.privKey)
	return encrypted, nonce[:], nil
}

func (cm *CryptoManager) Decrypt(encrypted []byte, nonce []byte, senderBoxPubKey *[32]byte) ([]byte, error) {
	var n [24]byte
	if len(nonce) != 24 {
		return nil, fmt.Errorf("invatelid nonce len")
	}
	copy(n[:], nonce)

	dectypted, ok := box.Open(nil, encrypted, &n, senderBoxPubKey, cm.privKey)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}
	return dectypted, nil
}

func SignMessage(privKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privKey, message)
}

// Загружает клю из файла или новый делает
func LoadOrGenerateKeys(path string) (*CryptoManager, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// new file
		return generateAndSave(path)
	}
	// load file
	return loadFromFile(path)
}

func generateAndSave(path string) (*CryptoManager, error) {
	pubSign, privSign, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Fatal to gen ed25:", err)
	}
	pubBox, privBox, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Fatal to gen box ed25:", err)
	}

	cm := &CryptoManager{
		SingPublic:  pubSign,
		SingPrivate: privSign,
		pubKey:      pubBox,
		privKey:     privBox,
	}
	if err := cm.save(path); err != nil {
		return nil, err
	}
	return cm, nil
}

func loadFromFile(path string) (*CryptoManager, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s storageStruct
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid key file format: %w", err)
	}

	// Восстанавливаем типы
	// Box ключи нужно превратить из slice []byte в array *[32]byte
	if len(s.BoxPublic) != 32 || len(s.BoxPrivate) != 32 {
		return nil, fmt.Errorf("corrupted box key")
	}

	var boxPub, boxPriv [32]byte
	copy(boxPub[:], s.BoxPublic)
	copy(boxPriv[:], s.BoxPrivate)

	return &CryptoManager{
		SingPublic:  ed25519.PublicKey(s.SignPublic),
		SingPrivate: ed25519.PrivateKey(s.SignPrivate),
		pubKey:      &boxPub,
		privKey:     &boxPriv,
	}, nil
}

func (cm *CryptoManager) save(path string) error {
	s := storageStruct{
		SignPublic:  []byte(cm.SingPublic),
		SignPrivate: []byte(cm.SingPrivate),
		BoxPublic:   cm.pubKey[:],
		BoxPrivate:  cm.privKey[:],
	}

	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	// 0600 - права на чтение/запись только для владельца
	return os.WriteFile(path, data, 0o600)
}
