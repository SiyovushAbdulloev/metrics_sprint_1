package metric

import "strconv"

func (h *Handler) validValue(metricType string, value string) (any, bool) {
	switch metricType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, false
		}
		return v, true
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, false
		}
		return v, true
	default:
		return nil, false
	}
}
