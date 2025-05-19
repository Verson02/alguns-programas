package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"strings"
)

// Conjuntos de caracteres para geração de senha
const (
	lowerLetters   = "abcdefghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers        = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
)

// Gera uma senha aleatória com os parâmetros especificados
func generatePassword(length int, includeNumbers, includeSpecial, includeUpper bool) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("o comprimento da senha deve ser maior que zero")
	}

	charPool := lowerLetters
	if includeUpper {
		charPool += upperLetters
	}
	if includeNumbers {
		charPool += numbers
	}
	if includeSpecial {
		charPool += specialChars
	}

	if len(charPool) == 0 {
		return "", fmt.Errorf("nenhum conjunto de caracteres foi selecionado")
	}

	var password strings.Builder
	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charPool))))
		if err != nil {
			return "", fmt.Errorf("erro ao gerar número aleatório: %v", err)
		}
		password.WriteByte(charPool[randIndex.Int64()])
	}

	return password.String(), nil
}

func main() {
	// Flags para os parâmetros do gerador de senha
	length := flag.Int("length", 12, "Comprimento da senha")
	includeNumbers := flag.Bool("numbers", true, "Incluir números na senha")
	includeSpecial := flag.Bool("special", true, "Incluir caracteres especiais na senha")
	includeUpper := flag.Bool("upper", true, "Incluir letras maiúsculas na senha")

	flag.Parse()

	// Gera a senha com os parâmetros fornecidos
	password, err := generatePassword(*length, *includeNumbers, *includeSpecial, *includeUpper)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println("Senha gerada:", password)
}
