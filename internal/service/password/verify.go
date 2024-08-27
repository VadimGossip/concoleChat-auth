package password

import "golang.org/x/crypto/bcrypt"

func (s *service) Verify(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
