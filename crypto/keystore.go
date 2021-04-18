package crypto

import (
	"crypto/rsa"
)

// KeyStore is a as the name says a keystore which stores multiple privatekeys
type KeyStore struct {
	keys map[string]*rsa.PrivateKey
}

// AddKeyFromFile Adds a private rsa key to the keystore from the file
func (store *KeyStore) AddKeyFromFile(keyName, fileName string) error {
	key, err := ParseRsaPrivateKeyFromFile(fileName)
	if err != nil {
		return err
	}
	store.keys[keyName] = key
	return nil
}

// AddKey Adds a private rsa key to the keystore
func (store *KeyStore) AddKey(keyName, pem string) error {
	key, err := ParseRsaPrivateKeyFromPemStr(pem)
	if err != nil {
		return err
	}
	store.keys[keyName] = key
	return nil
}

// Key Gets the private key for the given identifier
func (store *KeyStore) Key(keyName string) (*rsa.PrivateKey, bool) {
	val, isOk := store.keys[keyName]
	return val, isOk
}

// NewKeyStore Creates a new keystore
func NewKeyStore() *KeyStore {
	return &KeyStore{
		keys: make(map[string]*rsa.PrivateKey),
	}
}
