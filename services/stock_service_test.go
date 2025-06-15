package services

import (
	"testing"
)

func TestCalculateStockLevel(t *testing.T) {
	tests := []struct {
		name           string
		initialStock   float64
		suppliedStock  float64
		litersSold     float64
		expectedResult float64
	}{
		{
			name:           "Basic calculation",
			initialStock:   100,
			suppliedStock:  50,
			litersSold:     30,
			expectedResult: 120,
		},
		{
			name:           "No supplied stock",
			initialStock:   80,
			suppliedStock:  0,
			litersSold:     20,
			expectedResult: 60,
		},
		{
			name:           "No liters sold",
			initialStock:   60,
			suppliedStock:  40,
			litersSold:     0,
			expectedResult: 100,
		},
		{
			name:           "All stock sold",
			initialStock:   50,
			suppliedStock:  50,
			litersSold:     100,
			expectedResult: 0,
		},
		{
			name:           "Negative result",
			initialStock:   10,
			suppliedStock:  5,
			litersSold:     20,
			expectedResult: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := &Stock{
				RemainingStock: tt.initialStock,
				SuppliedStock:  tt.suppliedStock,
				LitersSold:     tt.litersSold,
			}
			result := stock.CalculateStockLevel()
			if result != tt.expectedResult {
				t.Errorf("got %v, want %v", result, tt.expectedResult)
			}
			if stock.RemainingStock != tt.expectedResult {
				t.Errorf("RemainingStock field = %v, want %v", stock.RemainingStock, tt.expectedResult)
			}
		})
	}
}
