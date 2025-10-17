package controls

import (
	"encoding/json"
	"fmt"

	"github.com/plusev-terminal/go-plugin-common/datapipe/types"
)

type NumberInputOption func(c *NumberInput)

type numberInputOptions struct {
	MinValue int `json:"min,omitempty"`
	MaxValue int `json:"max,omitempty"`
	Step     int `json:"step,omitempty"`
	Decimals int `json:"decimals,omitempty"`
}

func (o *numberInputOptions) ToMap() (map[string]any, error) {
	jsonData, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}

	var resultMap map[string]any
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json to map: %w", err)
	}

	return resultMap, nil
}

func WithMaxValue(maxValue int) NumberInputOption {
	return func(c *NumberInput) {
		c.numberInputOptions.MaxValue = maxValue
	}
}

func WithMinValue(minValue int) NumberInputOption {
	return func(c *NumberInput) {
		c.numberInputOptions.MinValue = minValue
	}
}

func WithStep(step int) NumberInputOption {
	return func(c *NumberInput) {
		c.numberInputOptions.Step = step
	}
}

func WithDecimals(decimals int) NumberInputOption {
	return func(c *NumberInput) {
		c.numberInputOptions.Decimals = decimals
	}
}

type NumberInput struct {
	Control
	numberInputOptions
}

func NewNumberInput(label, name string, options ...NumberInputOption) *types.GuiControl {
	c := &NumberInput{
		Control: Control{
			Label: label,
			Name:  name,
			Type:  types.NUMBER_INPUT,
		},
		numberInputOptions: numberInputOptions{
			MinValue: 1,
			MaxValue: 999999,
			Step:     1,
			Decimals: 0,
		},
	}

	for _, opt := range options {
		opt(c)
	}

	optionsMap, err := c.numberInputOptions.ToMap()
	if err != nil {
		panic(fmt.Sprintf("failed to convert options to map: %v", err))
	}

	return NewControl(label, name, types.NUMBER_INPUT, optionsMap)
}
