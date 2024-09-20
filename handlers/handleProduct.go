package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (api *Api) GetProduct(w http.ResponseWriter, r *http.Request) { //! у нас есть id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}

	product, err := api.productStorage.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, 500)
		return
	}

	resp, err := json.Marshal(product)
	if err != err {
		http.Error(w, `{"error":"json error"}`, 500)
		return
	}

	w.Write(resp)
}
