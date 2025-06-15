package services

type Stock struct{
	RemainingStock float64
	SuppliedStock float64
	LitersSold float64
}
func (s *Stock) CalculateStockLevel()float64 {
	s.RemainingStock = (s.RemainingStock + s.SuppliedStock) - s.LitersSold
	return s.RemainingStock
}