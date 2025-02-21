package date

import "time"

// Função para converter string para data
func ParseISODate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
