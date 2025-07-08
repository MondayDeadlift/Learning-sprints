package calculationService

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationByID(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (c *calcService) calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression) // создание выражения (3+3)
	if err != nil {
		return "", err // если выражение не валидно
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err
}

// CreateCalculation implements CalculationService.
func (c *calcService) CreateCalculation(expression string) (Calculation, error) {
	result, err := c.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := c.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

// DeleteCalculation implements CalculationService.
func (c *calcService) DeleteCalculation(id string) error {
	return c.repo.DeleteCalculation(id)
}

// GetAllCalculations implements CalculationService.
func (c *calcService) GetAllCalculations() ([]Calculation, error) {
	return c.repo.GetAllCalculations()
}

// GetCalculationByID implements CalculationService.
func (c *calcService) GetCalculationByID(id string) (Calculation, error) {
	return c.repo.GetCalculationByID(id)
}

// UpdateCalculation implements CalculationService.
func (c *calcService) UpdateCalculation(id string, expression string) (Calculation, error) {
	calc, err := c.repo.GetCalculationByID(id)
	if err != nil {
		return Calculation{}, err
	}

	result, err := c.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := c.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}
