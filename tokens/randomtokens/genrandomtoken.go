package Tokens

import "crypto/rand"

//Thanks to https://www.socketloop.com/tutorials/golang-how-to-generate-random-string
func GenerateRandomAlphaNumericString(strSize int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Token server internal error: " + err.Error())
	}
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
