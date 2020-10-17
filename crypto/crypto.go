package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/chacha20poly1305"
)

// GenerateRsaKeyPair generates a pub/priv rsa key pair
func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

// ExportRsaPrivateKeyAsPemStr exports the rsa priv key in pem format
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {

	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPEM)
}

// ExportRsaPublicKeyAsPemStr exports the rsa pub key in pem format
func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {

	pubkeyBytes := x509.MarshalPKCS1PublicKey(pubkey)

	pubkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPEM)
}

// ExportRsaPrivateKeyToFile exports the rsa priv key to a file
func ExportRsaPrivateKeyToFile(fileName string, key *rsa.PrivateKey) error {

	// Export the keys to pem string
	priv_pem := ExportRsaPrivateKeyAsPemStr(key)

	// Write to file
	return ioutil.WriteFile(fileName, []byte(priv_pem), 0644)
}

// ExportRsaPublicKeyToFile exports the rsa pub key to a file
func ExportRsaPublicKeyToFile(fileName string, key *rsa.PublicKey) error {

	// Export the keys to pem string
	pubKey := ExportRsaPublicKeyAsPemStr(key)

	// Write to file
	return ioutil.WriteFile(fileName, []byte(pubKey), 0644)
}

// ParseRsaPrivateKeyFromPemStr parses the private key from the given pem string
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// ParseRsaPublicKeyFromPemStr parses the pub key from the given pem string
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

// ParseRsaPrivateKeyFromFile parses the priv key from the given file
func ParseRsaPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {

	// read the whole file at once
	priv_pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Import the keys from pem string
	return ParseRsaPrivateKeyFromPemStr(string(priv_pem))
}

// GetHash returns the hash of a message
func GetHash(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)
	return h.Sum(nil)
}

// Sign returns a signature made by combining the message and the signers private key
func Sign(msg []byte, key *rsa.PrivateKey) (signature []byte, err error) {
	hs := GetHash(msg)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hs)
}

// Verify checks if a message is signed by a given Public Key
func Verify(msg []byte, sig []byte, pk *rsa.PublicKey) error {
	hs := GetHash(msg)
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hs, sig)
}

// EncryptUsingSymmKey Enrypts data using ChaCha2020
// TODO: Give option to change encryption algo
func EncryptUsingSymmKey(msg, privKey []byte) ([]byte, error) {
	hs := GetHash(privKey)
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(hs); err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	var encryptedMsg []byte
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg = aead.Seal(nonce, nonce, msg, nil)

	return encryptedMsg, nil
}

// DecryptUsingSymmKey Decrypts data using ChaCha2020
// TODO: Give option to change encryption algo
func DecryptUsingSymmKey(encryptedMsg, privKey []byte) ([]byte, error) {
	hs := GetHash(privKey)
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(hs); err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	if len(encryptedMsg) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
