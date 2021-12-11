package cpfCNPJ

import (
	"strconv"
)

func ValidateCpfCNPJ(cpfCNPJ string) bool {

	if len(cpfCNPJ) == 11 {
		if validateFirstCpfDigit(cpfCNPJ) {

			if validateSecondCpfDigit(cpfCNPJ) {
				return true
			}

			return false
		}
	} else if len(cpfCNPJ) == 14 {
		if validateFirstCnpjDigit(cpfCNPJ) {

			if validateSecondCnpjDigit(cpfCNPJ) {
				return true
			}

			return false
		}
	}

	return false

}

func validateFirstCpfDigit(cpf string) bool {
	newCpf := cpf[:len(cpf)-2]
	numberToMul := 10
	result := 0

	for _, digit := range newCpf {
		numberToSum, _ := strconv.Atoi(string(digit))
		result = result + (numberToSum * numberToMul)
		numberToMul--
	}

	result = result * 10 % 11
	if result == 10 {
		result = 0
	}

	firstDigit, _ := strconv.Atoi(string(cpf[len(cpf)-2]))

	if result == firstDigit {
		return true
	}

	return false
}

func validateSecondCpfDigit(cpf string) bool {
	newCpf := cpf[:len(cpf)-1]
	numberToMul := 11
	result := 0

	for _, digit := range newCpf {
		numberToSum, _ := strconv.Atoi(string(digit))
		result = result + (numberToSum * numberToMul)
		numberToMul--
	}

	result = result * 10 % 11
	if result == 10 {
		result = 0
	}

	secondDigit, _ := strconv.Atoi(string(cpf[len(cpf)-1]))

	if result == secondDigit {
		return true
	}

	return false
}

func validateFirstCnpjDigit(cnpj string) bool {
	newCnpj := cnpj[:len(cnpj)-2]
	numberToMul := 5
	result := 0

	for _, digit := range newCnpj {
		numberToSum, _ := strconv.Atoi(string(digit))
		result = result + (numberToSum * numberToMul)
		numberToMul--
		if numberToMul == 1 {
			numberToMul = 9
		}
	}

	result = result % 11
	if result < 2 {
		result = 0
	} else {
		result = 11 - result
	}

	firstDigit, _ := strconv.Atoi(string(cnpj[len(cnpj)-2]))

	if result == firstDigit {
		return true
	}

	return false
}

func validateSecondCnpjDigit(cpf string) bool {
	newCpf := cpf[:len(cpf)-1]
	numberToMul := 6
	result := 0

	for _, digit := range newCpf {
		numberToSum, _ := strconv.Atoi(string(digit))
		result = result + (numberToSum * numberToMul)
		numberToMul--
		if numberToMul == 1 {
			numberToMul = 9
		}
	}

	result = result % 11
	if result < 2 {
		result = 0
	} else {
		result = 11 - result
	}

	secondDigit, _ := strconv.Atoi(string(cpf[len(cpf)-1]))

	if result == secondDigit {
		return true
	}

	return false
}
