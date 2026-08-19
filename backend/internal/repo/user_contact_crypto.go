package repo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	userContactCipherPrefix   = "enc:"
	userContactKeyInfo        = "tier0-enterprise:sys_user_info:contact:aes-gcm"
	userContactEmbeddedKeyHex = "502f23e39140832460cd2d6a81d2d0b2403473399c4c5b78898abf3b50df10a0"
	userContactFieldEmail     = "email"
	userContactFieldPhone     = "phone"
)

var errUserContactCiphertextInvalid = errors.New("user contact ciphertext is invalid")

type userContactCipher struct {
	aead cipher.AEAD
}

var activeUserContactCipher *userContactCipher

func newUserContactCipher() (*userContactCipher, error) {
	// This embedded key intentionally protects only against plaintext exposure
	// from the database alone. It is part of the persistent ciphertext format
	// and must remain stable across releases or existing contacts become unreadable.
	key, err := hex.DecodeString(userContactEmbeddedKeyHex)
	if err != nil || len(key) != 32 {
		return nil, errors.New("embedded user contact key must encode exactly 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create user contact cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create user contact AEAD: %w", err)
	}
	return &userContactCipher{aead: aead}, nil
}

func encryptUserContact(field, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	if activeUserContactCipher == nil {
		return "", errors.New("user contact cipher is not initialized")
	}
	nonce := make([]byte, activeUserContactCipher.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate user contact nonce: %w", err)
	}
	sealed := activeUserContactCipher.aead.Seal(nonce, nonce, []byte(plaintext), userContactAAD(field))
	return userContactCipherPrefix + base64.RawStdEncoding.EncodeToString(sealed), nil
}

func decryptUserContact(field, stored string) (string, error) {
	if stored == "" {
		return "", nil
	}
	if activeUserContactCipher == nil {
		return "", errors.New("user contact cipher is not initialized")
	}
	if !strings.HasPrefix(stored, userContactCipherPrefix) {
		return "", errUserContactCiphertextInvalid
	}
	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(stored, userContactCipherPrefix))
	if err != nil || len(raw) < activeUserContactCipher.aead.NonceSize() {
		return "", errUserContactCiphertextInvalid
	}
	nonceSize := activeUserContactCipher.aead.NonceSize()
	plaintext, err := activeUserContactCipher.aead.Open(nil, raw[:nonceSize], raw[nonceSize:], userContactAAD(field))
	if err != nil {
		return "", errUserContactCiphertextInvalid
	}
	return string(plaintext), nil
}

func userContactAAD(field string) []byte {
	return []byte(userContactKeyInfo + ":" + field)
}

func userContactIsEncrypted(stored string) bool {
	return strings.HasPrefix(stored, userContactCipherPrefix)
}

func encryptUserContactPair(email, phone string) (string, string, error) {
	encryptedEmail, err := encryptUserContact(userContactFieldEmail, email)
	if err != nil {
		return "", "", err
	}
	encryptedPhone, err := encryptUserContact(userContactFieldPhone, phone)
	if err != nil {
		return "", "", err
	}
	return encryptedEmail, encryptedPhone, nil
}

func decryptUserContactPair(email, phone string) (string, string, error) {
	plaintextEmail, err := decryptUserContact(userContactFieldEmail, email)
	if err != nil {
		return "", "", err
	}
	plaintextPhone, err := decryptUserContact(userContactFieldPhone, phone)
	if err != nil {
		return "", "", err
	}
	return plaintextEmail, plaintextPhone, nil
}
