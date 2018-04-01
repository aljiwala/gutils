package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const minL2N = 14

var (
	// DefaultTimeout ...
	DefaultTimeout = 5 * time.Second
	dursMu         sync.Mutex
	durs           = make([]time.Duration, 0, 8)
)

type (
	// WriteCloser ...
	WriteCloser struct {
		io.Writer
	}

	// SecretWriter ...
	SecretWriter struct {
		key, nonce []byte
		w          io.WriteCloser
		buf        bytes.Buffer
	}

	// Key ...
	Key struct {
		Bytes []byte `json:"-"`
		Salt  []byte
		L2N   uint
		R, P  int
	}
)

// Close ...
func (wc WriteCloser) Close() error {
	if c, ok := wc.Writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Write ...
func (sw *SecretWriter) Write(p []byte) (int, error) {
	return sw.buf.Write(p)
}

// Close ...
func (sw *SecretWriter) Close() error {
	var (
		key   [32]byte
		nonce [24]byte
	)

	copy(key[:], sw.key)
	copy(nonce[:], sw.nonce)
	out := make([]byte, 0, sw.buf.Len()+secretbox.Overhead)
	_, err := sw.w.Write(secretbox.Seal(out, sw.buf.Bytes(), &nonce, &key))
	if err != nil {
		return err
	}

	return sw.w.Close()
}

func (key Key) String() string {
	k := struct {
		Bytes, Salt []byte
		L2N         uint
		R, P        int
	}(key)
	b, err := json.Marshal(k)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// Populate ...
func (key *Key) Populate(password []byte, keyLen int) error {
	var err error
	key.Bytes, err = scrypt.Key(password, key.Salt, 1<<key.L2N, key.R, key.P, keyLen)
	return err
}

// Salt creates new random salt with the given length.
func Salt(saltLen int) (salt []byte, err error) {
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt)
	return
}

// Open ...
func Open(filename string, passphrase []byte) (Key, io.Reader, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return Key{}, nil, err
	}
	return OpenReader(fh, passphrase)
}

// OpenReader ...
func OpenReader(r io.Reader, passphrase []byte) (Key, io.Reader, error) {
	var (
		k     [32]byte
		nonce [24]byte
		key   Key
	)

	dec := json.NewDecoder(r)
	if err := dec.Decode(&key); err != nil {
		return key, nil, err
	}
	if err := key.Populate(passphrase, 32); err != nil {
		return key, nil, err
	}

	box, err := ioutil.ReadAll(io.MultiReader(dec.Buffered(), r))
	if err != nil {
		return key, nil, err
	}

	box = box[1:] // Trim \n of json.Encoder.Encode.
	out := make([]byte, 0, len(box)-secretbox.Overhead)

	copy(nonce[:], key.Salt)
	copy(k[:], key.Bytes)
	data, ok := secretbox.Open(out, box, &nonce, &k)
	if !ok {
		return key, nil, errors.New("failed open box")
	}

	return key, bytes.NewReader(data), nil
}

// GenKey derives a key from the password, using scrypt.
// It tries to create the strongest key within the given time window.
func GenKey(password []byte, saltLen, keyLen int, timeout time.Duration,
) (Key, error) {
	salt, err := Salt(saltLen)
	if err != nil {
		return Key{}, err
	}
	key := Key{Salt: salt, R: 8, P: 1, L2N: minL2N}
	dursMu.Lock()
	defer dursMu.Unlock()
	for i, d := range durs {
		if d < timeout {
			key.L2N = minL2N + uint(i)
			continue
		}
		if d > timeout {
			key.L2N--
		}
		break
	}
	deadline := time.Now().Add(timeout)
	for now := time.Now(); now.Before(deadline); {
		if key.Bytes, err = scrypt.Key(password, salt, 1<<key.L2N, key.R, key.P, keyLen); err != nil {
			return key, err
		}
		now2 := time.Now()
		dur := now2.Sub(now)
		i := int(key.L2N - minL2N)
		if len(durs) <= i {
			if cap(durs) > i {
				durs = durs[:i+1]
			} else {
				durs = append(durs, make([]time.Duration, len(durs))...)
			}
		}
		durs[key.L2N-minL2N] = dur
		now = now2

		if now.Add(2 * dur).After(deadline) {
			break
		}
		key.L2N++
	}
	return key, nil
}

// Encrypt binary data to a base64 string with AES using the key provided.
func Encrypt(keyString string, data []byte) (string, error) {
	key := []byte(keyString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt a base64 string to binary data with AES using the key provided.
func Decrypt(keyString string, base64Data string) ([]byte, error) {
	key := []byte(keyString)
	ciphertext, err := base64.URLEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Ciphertext provided is smaller than AES block size")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
