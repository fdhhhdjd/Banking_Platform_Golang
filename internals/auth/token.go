package auth

import "os"

var JwtMaker Maker

func GetJWTMaker() (Maker, error) {

	secretKey := os.Getenv("SECRET_KEY_TOKEN")
	JwtMaker, err := NewJWTMaker(secretKey)
	if err != nil {
		return nil, err
	}
	return JwtMaker, nil
}
