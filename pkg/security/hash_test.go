package security

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 1000),
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "test@#$%^&*()_+-=[]{}|;':\".,/<>?",
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "тест密码パスワード",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the hash format
				if !strings.HasPrefix(hashedPassword, "$argon2id$") {
					t.Errorf("HashPassword() hash format invalid, expected to start with $argon2id$")
				}

				// Verify the hash has correct number of parts
				parts := strings.Split(hashedPassword, "$")
				if len(parts) != 6 {
					t.Errorf("HashPassword() hash should have 6 parts, got %d", len(parts))
				}

				// Verify version
				expectedVersion := fmt.Sprintf("v=%d", argon2.Version)
				if parts[2] != expectedVersion {
					t.Errorf("HashPassword() version mismatch, expected %s, got %s", expectedVersion, parts[2])
				}

				// Verify parameters
				expectedParams := fmt.Sprintf("models=%d,t=%d,p=%d", A2IDmemory, A2IDtime, A2IDthreads)
				if parts[3] != expectedParams {
					t.Errorf("HashPassword() parameters mismatch, expected %s, got %s", expectedParams, parts[3])
				}

				// Verify salt and hash are valid base64
				_, err := base64.RawStdEncoding.DecodeString(parts[4])
				if err != nil {
					t.Errorf("HashPassword() salt is not valid base64: %v", err)
				}

				_, err = base64.RawStdEncoding.DecodeString(parts[5])
				if err != nil {
					t.Errorf("HashPassword() hash is not valid base64: %v", err)
				}
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "testpassword"

	// Generate multiple hashes of the same password
	hashes := make([]string, 10)
	for i := 0; i < 10; i++ {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() failed: %v", err)
		}
		hashes[i] = hash
	}

	// Verify all hashes are unique (due to random salt)
	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes); j++ {
			if hashes[i] == hashes[j] {
				t.Errorf("HashPassword() generated duplicate hashes: %s", hashes[i])
			}
		}
	}
}

func TestComparePassword(t *testing.T) {
	// Generate a valid hash for testing
	validPassword := "testpassword123"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantMatch      bool
		wantErr        bool
	}{
		{
			name:           "correct password",
			hashedPassword: validHash,
			password:       validPassword,
			wantMatch:      true,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: validHash,
			password:       "wrongpassword",
			wantMatch:      false,
			wantErr:        false,
		},
		{
			name:           "empty password against valid hash",
			hashedPassword: validHash,
			password:       "",
			wantMatch:      false,
			wantErr:        false,
		},
		{
			name:           "malformed hash - too few parts",
			hashedPassword: "$argon2id$v=19$models=65536",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
		{
			name:           "malformed hash - invalid parameters",
			hashedPassword: "$argon2id$v=19$invalid=params$c2FsdA$aGFzaA",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
		{
			name:           "malformed hash - invalid salt base64",
			hashedPassword: "$argon2id$v=19$models=65536,t=3,p=4$invalid-base64!@#$aGFzaA",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
		{
			name:           "malformed hash - invalid hash base64",
			hashedPassword: "$argon2id$v=19$models=65536,t=3,p=4$c2FsdA$invalid-base64!@#",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
		{
			name:           "empty hash string",
			hashedPassword: "",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
		{
			name:           "hash with wrong algorithm",
			hashedPassword: "$bcrypt$12$abcdefghijklmnopqrstuvwxyz",
			password:       "testpassword",
			wantMatch:      false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ComparePassword(tt.hashedPassword, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("ComparePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if match != tt.wantMatch {
				t.Errorf("ComparePassword() match = %v, wantMatch %v", match, tt.wantMatch)
			}
		})
	}
}

func TestComparePasswordWithDifferentParameters(t *testing.T) {
	// Create a hash with custom parameters to test parameter extraction
	password := "testpassword"
	salt := []byte("testsalt12345678") // 16 bytes
	customMemory := uint32(32 * 1024)
	customTime := uint32(2)
	customThreads := uint8(2)
	customKeyLen := uint32(24)

	hash := argon2.IDKey([]byte(password), salt, customTime, customMemory, customThreads, customKeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	customHashedPassword := fmt.Sprintf("$argon2id$v=%d$models=%d,t=%d,p=%d$%s$%s",
		argon2.Version, customMemory, customTime, customThreads, b64Salt, b64Hash)

	t.Logf("Custom hash: %s", customHashedPassword)

	// Test that comparison works with different parameters
	match, err := ComparePassword(customHashedPassword, password)
	if err != nil {
		t.Errorf("ComparePassword() with custom parameters failed: %v", err)
	}
	if !match {
		t.Errorf("ComparePassword() with custom parameters should match")
	}

	// Test with wrong password
	match, err = ComparePassword(customHashedPassword, "wrongpassword")
	if err != nil {
		t.Errorf("ComparePassword() with custom parameters failed: %v", err)
	}
	if match {
		t.Errorf("ComparePassword() with custom parameters should not match wrong password")
	}
}

func TestHashAndCompareIntegration(t *testing.T) {
	passwords := []string{
		"password123",
		"",
		"verylongpasswordwithlotsofcharacters",
		"helloworld",
		"pass@#$%^&*()",
	}

	for _, password := range passwords {
		t.Run(fmt.Sprintf("password_%s", password), func(t *testing.T) {
			// Hash the password
			hashedPassword, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() failed: %v", err)
			}

			// Compare with correct password
			match, err := ComparePassword(hashedPassword, password)
			if err != nil {
				t.Errorf("ComparePassword() failed: %v", err)
			}
			if !match {
				t.Errorf("ComparePassword() should match the original password")
			}

			// Compare with incorrect password
			wrongPassword := password + "wrong"
			match, err = ComparePassword(hashedPassword, wrongPassword)
			if err != nil {
				t.Errorf("ComparePassword() failed: %v", err)
			}
			if match {
				t.Errorf("ComparePassword() should not match incorrect password")
			}
		})
	}
}

func TestConstantTimeComparison(t *testing.T) {
	// This test ensures that comparison uses constant time
	// We can't easily test timing, but we can verify behavior
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	// Test multiple incorrect passwords of different lengths
	incorrectPasswords := []string{
		"a",
		"ab",
		"abc",
		"wrongpassword",
		"verylongwrongpasswordthatisdifferentlength",
	}

	for _, wrongPassword := range incorrectPasswords {
		match, err := ComparePassword(hashedPassword, wrongPassword)
		if err != nil {
			t.Errorf("ComparePassword() failed for password '%s': %v", wrongPassword, err)
		}
		if match {
			t.Errorf("ComparePassword() should not match for password '%s'", wrongPassword)
		}
	}
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword() failed: %v", err)
		}
	}
}

func BenchmarkComparePassword(b *testing.B) {
	password := "benchmarkpassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		b.Fatalf("HashPassword() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ComparePassword(hashedPassword, password)
		if err != nil {
			b.Fatalf("ComparePassword() failed: %v", err)
		}
	}
}
