package auth

var JwtMaker Maker

func GetJWTMaker() (Maker, error) {
	secretKey := "6892e55f3d2d1fae3665a0fa6586987cc4c4f3c3a5ff2adb1bd36a6351a24e2d"
	JwtMaker, err := NewJWTMaker(secretKey)
	if err != nil {
		return nil, err
	}
	return JwtMaker, nil
}
