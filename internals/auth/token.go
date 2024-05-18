package auth

import "log"

var jwtMaker Maker

func GetJWTMaker() Maker {
	// TOKEN
	secretKey := "123123123213213131231232132132131231231232"
	var err error
	jwtMaker, err = NewJWTMaker(secretKey)
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}
	return jwtMaker
}
