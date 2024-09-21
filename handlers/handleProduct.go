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
		logger.Error("error", err)
		return
	}

	resp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, `{"error":"json error"}`, 500)
		logger.Error("error", err)
		return
	}

	w.Write(resp)
}

func (api *Api) GetProduct(w http.ResponseWriter, r *http.Request) { //! у нас есть id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		logger.Error("error", err)
		return
	}

	product, err := api.productStorage.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, 500)
		logger.Error("error", err)
		return
	}

	resp, err := json.Marshal(product)
	if err != err {
		http.Error(w, `{"error":"json error"}`, 500)
		logger.Error("error", err)
		return
	}

	w.Write(resp)
}

//TODO реализовать доступ только продовцу
// func (api *Api) AddProduct(w http.ResponseWriter, r *http.Request) { //! принимаем post json
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, `{"error": "server error"}`, 500)
// 		logger.Error("error", err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var product storage.Product

// 	err = json.Unmarshal(body, product)
// 	if
// }
