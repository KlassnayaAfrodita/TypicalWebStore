package handlers

import (
	"encoding/json"
	"net/http"
)

func (api *Api) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := api.productStorage.GetProducts()
	if err != nil {
		http.Error(w, `{"error":"db error"}`, 500)
		return
	}

	resp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, `{"error":"json error"}`, 500)
	}

	w.Write(resp)
}

func (api *Api) GetProduct(w http.ResponseWriter, r *http.Request) {

}
