package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/tomknobel/ip2country/internal/controllers"
	"github.com/tomknobel/ip2country/internal/db"
)

func InitIp2CountryRouter(r chi.Router, ipDb db.Db) {
	err := ipDb.Connect()
	if err != nil {
		return
	}
	defer ipDb.Close()
	ic := controllers.NewIpController(ipDb)

	r.Get("/find-country", ic.GetCountryByIp)
}
