package handlers

import (
	"GOLANG/net/projext-net/storage"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// type Result struct {
// 	Body interface{} `json:"body,omitempty"` //! это для того, чтобы можно было юзать разные типы
// 	Err  string      `json:"err,omitempty"`
// }

// func (r *Result) formResult(field string) ([]byte, error) {
// 	body := map[string]interface{}{
// 		field: r.Body,
// 	}
// 	ret
// }

// type User struct {
// 	cart storage.ProductStore
// 	user *storage.User
// }

type Api struct {
	session *storage.Session
	users   *storage.UserStorage
}

func (api *Api) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("hello"))
	products, err := api.users.GetUser(id).GetProducts()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	resp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(resp)
}

func (c *storage.User) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, 400)
		return
	}

	product, err := c.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error": "server error"`, 500)
		return
	}

	w.Write(resp)
}

func (c *storage.User) AddProduct(w http.ResponseWriter, r *http.Request) { //! ожидается json с новым пользователем
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.Cart.AddProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
	}

	w.Write([]byte(body))
}

func (c *storage.User) ChangeProduct(w http.ResponseWriter, r *http.Request) { //! получаем json полного пользователя
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.Cart.ChangeProducts(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}
	w.Write(body)
}

func (c *storage.User) DeleteProduct(w http.ResponseWriter, r *http.Request) { //! получаем только id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}

	product, err := c.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}

	product, err = c.Cart.DeleteProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error":"server error"`, 500)
	}

	w.Write(resp)
}
