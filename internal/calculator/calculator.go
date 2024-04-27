package calculator

import (
	"strconv"
	"strings"
	"time"
)

type Operation struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

var (
	Operations = []Operation{
		{"Сложение", 2 * time.Second},
		{"Вычитание", 3 * time.Second},
		{"Умножение", 2 * time.Second},
		{"Деление", 5 * time.Second},
	}
)

func EvaluateExpression(expression string) float64 {
	// Удаляем все пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	// Создаем стеки для операндов и операторов
	operandStack := make([]float64, 0)
	operatorStack := make([]rune, 0)

	// Функция для выполнения операции
	performOperation := func() {
		if len(operandStack) < 2 || len(operatorStack) == 0 {
			return
		}

		b := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]

		a := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]

		op := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]

		var result float64
		switch op {
		case '+':
			time.Sleep(Operations[0].Duration)
			result = a + b
		case '-':
			time.Sleep(Operations[1].Duration)
			result = a - b
		case '*':
			time.Sleep(Operations[2].Duration)
			result = a * b
		case '/':
			time.Sleep(Operations[3].Duration)
			result = a / b
		}
		operandStack = append(operandStack, result)
	}

	// Обходим каждый символ в выражении
	for _, char := range expression {
		switch char {
		case '(':
			operatorStack = append(operatorStack, char)
		case ')':
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != '(' {
				performOperation()
			}
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] == '(' {
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
		case '+', '-':
			for len(operatorStack) > 0 && (operatorStack[len(operatorStack)-1] == '+' ||
				operatorStack[len(operatorStack)-1] == '-' || operatorStack[len(operatorStack)-1] == '*' || operatorStack[len(operatorStack)-1] == '/') {
				performOperation()
			}
			operatorStack = append(operatorStack, char)
		case '*', '/':
			for len(operatorStack) > 0 && (operatorStack[len(operatorStack)-1] == '*' ||
				operatorStack[len(operatorStack)-1] == '/') {
				performOperation()
			}
			operatorStack = append(operatorStack, char)
		default:
			// Если символ - цифра или точка, добавляем ее в стек операндов
			operand, _ := strconv.ParseFloat(string(char), 64)
			operandStack = append(operandStack, operand)
		}
	}

	for len(operatorStack) > 0 {
		performOperation()
	}

	return operandStack[0]
}
