package services 

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

/*compares the parsed password with the password from the database
@params hashedPassword
@params ParsedPassword
*/
func CompareHashAndPassword(hashedPassword ,parsedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(parsedPassword))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}