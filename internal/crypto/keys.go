package crypto

import "encoding/base64"

func (km *CryptoManager) PublicBase64() string {
	return base64.StdEncoding.EncodeToString(km.SingPublic)
	// Или BoxPublic, смотря что ты используешь как ID. Обычно SignPublic (Ed25519) это ID.
}
