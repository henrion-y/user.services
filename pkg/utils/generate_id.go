package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	Str_Base_Alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumBer_Base_Alphabet = "1234567890"
	All_Base_Alphabet    = Str_Base_Alphabet + NumBer_Base_Alphabet
)

func GenerateId() (string, error) {
	return gonanoid.Generate(All_Base_Alphabet, 16)
}

func GenerateUsername() (string, error) {
	randStr, err := gonanoid.Generate(All_Base_Alphabet, 6)
	if err != nil {
		return "", err
	}
	return "用户" + randStr, nil
}

func GenerateSmsCode() (string, error) {
	return gonanoid.Generate(NumBer_Base_Alphabet, 6)
}
