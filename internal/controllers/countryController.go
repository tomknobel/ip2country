package controllers

import (
	"github.com/tomknobel/ip2country/internal/db"
	"github.com/tomknobel/ip2country/pkg/utils"
	"net/http"
)

type countryController struct {
	db db.Db
}
type ErrorResponse struct {
	error string
}

func NewCountryController(db db.Db) *countryController {
	return &countryController{
		db,
	}
}
func (c *countryController) GetCountryByIp(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	country, err := c.db.Find(ip)
	if err != nil {
		panic(err)
	}
	utils.JsonResponse(w, country)

}
