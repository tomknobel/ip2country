package controllers

import (
	"github.com/tomknobel/ip2country/internal/db"
	"github.com/tomknobel/ip2country/pkg/utils"
	"net/http"
)

type ipController struct {
	db db.Db
}
type ErrorResponse struct {
	error string
}

func NewIpController(db db.Db) *ipController {
	return &ipController{
		db,
	}
}
func (ic *ipController) GetCountryByIp(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	country, err := ic.db.Find(ip)
	if err != nil {
		utils.JsonResponse(w, ErrorResponse{
			error: err.Error(),
		})

	}
	utils.JsonResponse(w, country)

}
