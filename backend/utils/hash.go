package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ArgonParams defines the computational cost for password hashing
type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// Current recommended OWASP parameters for Argon2id
var params = &ArgonParams{
	Memory:      64 * 1024, // 64 MB of RAM required to compute hash
	Iterations:  3,         // 3 linear passes
	Parallelism: 2,         // 2 threads
	SaltLength:  16,        // 16 bytes of random salt
	KeyLength:   32,        // 32 bytes output key
}

// HashPassword creates a military-grade secure hash using Argon2id
func HashPassword(password string) (string, error) {
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Compute the memory-hard hash
	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
	
	// Encode in standard PHC string format (e.g., $argon2id$v=19$m=65536,t=3,p=2$salt$hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

// VerifyPassword checks a plaintext password against a PHC formatted hash in constant time
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Recompute hash with the extracted salt and parameters
	hashToCompare := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(decodedHash)))

	// Use ConstantTimeCompare to prevent timing attacks
	if subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1 {
		return true, nil
	}
	return false, nil
}
