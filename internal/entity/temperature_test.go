package entity

import (
	"testing"
)

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name     string
		kelvin   float64
		expected Temperature
	}{
		{
			name:   "Zero Kelvin",
			kelvin: 0,
			expected: Temperature{
				Farenheit: -459.67,
				Celsius:   -273,
				Kelvin:    0,
			},
		},
		{
			name:   "Freezing Point of Water",
			kelvin: 273.15,
			expected: Temperature{
				Farenheit: 32,
				Celsius:   0,
				Kelvin:    273.15,
			},
		},
		{
			name:   "Boiling Point of Water",
			kelvin: 373.15,
			expected: Temperature{
				Farenheit: 212,
				Celsius:   100,
				Kelvin:    373.15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, err := NewTemperature("São Paulo", tt.kelvin)
			if err != nil {
				t.Fatalf("NewTemperature() error = %v", err)
			}
			if temp.Farenheit != tt.expected.Farenheit {
				t.Errorf("NewTemperature() Farenheit = %v, expected %v", temp.Farenheit, tt.expected.Farenheit)
			}
			if temp.Celsius != tt.expected.Celsius {
				t.Errorf("NewTemperature() Celsius = %v, expected %v", temp.Celsius, tt.expected.Celsius)
			}
			if temp.Kelvin != tt.expected.Kelvin {
				t.Errorf("NewTemperature() Kelvin = %v, expected %v", temp.Kelvin, tt.expected.Kelvin)
			}
			if temp.City != "São Paulo" {
				t.Errorf("NewTemperature() City = %v, expected %v", temp.City, "São Paulo")
			}
		})
	}
}
