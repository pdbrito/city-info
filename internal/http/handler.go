package http

import (
	"encoding/json"
	"github.com/pdbrito/city-info/internal/cityinfoservice"
	"net/http"
)

type handler struct {
	cityInfoService cityinfoservice.CityInfoService
}

// NewHandler returns a new http.Handler that can serve our routes.
func NewHandler(cis cityinfoservice.CityInfoService) http.Handler {
	h := handler{cityInfoService: cis}

	mux := http.NewServeMux()
	mux.HandleFunc("/city-info", h.getCityInfo)

	return mux
}

// ErrorResponse defines our error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

func (h handler) getCityInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	city := r.URL.Query().Get("name")

	if city == "" {
		resp := ErrorResponse{
			Error: "missing required name query parameter",
		}
		responseData, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(responseData)
		return
	}

	cityInfo, err := h.cityInfoService.CurrentStatus(city)

	if err != nil {
		resp := ErrorResponse{
			Error: err.Error(),
		}
		responseData, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(responseData)
		return
	}

	responseData, _ := json.Marshal(cityInfo)
	w.Write(responseData)
}
