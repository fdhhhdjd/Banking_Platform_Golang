package auth

import "os"

var JwtMaker Maker

func init() {
	secretKey := os.Getenv("SECRET_KEY_TOKEN")
	var err error
	JwtMaker, err = NewJWTMaker(secretKey)
	if err != nil {
		panic("failed to create JWT maker: " + err.Error())
	}
}
func GetJWTMaker() (Maker, error) {
	return JwtMaker, nil
}
