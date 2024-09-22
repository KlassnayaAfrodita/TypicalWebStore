package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/KlassnayaAfrodita/mylib/storage"

	"github.com/gorilla/mux"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

type Api struct {
	session        *storage.Session
	users          *storage.UserStorage
	productStorage *storage.ProductStorage
}

func NewApi() *Api {
	productStorage := storage.NewProductStorage()
	_, err := productStorage.AddProduct(storage.Product{
		ID:       1,
		Name:     "laptop",
		Price:    1000.,
		Quantity: 5,
		About:    "laptop",
	})
	if err != nil {
		panic(err)
	}
	return &Api{
		session:        storage.NewSession(),
		users:          storage.NewUserStorage(),
		productStorage: &productStorage,
	}
}

func (api *Api) GetAllCart(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("hello"))
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userID, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	products, err := user.Cart.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), 500)
		logger.Error("error", err)
		return
	}

	resp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), 500)
		logger.Error("error", err)
		return
	}
	w.Write(resp)
}

func (api *Api) GetProductCart(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, 400)
		logger.Error("error", err)
		return
	}

	product, err := user.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error": "server error"`, 500)
		logger.Error("error", err)
		return
	}

	w.Write(resp)
}

func (api *Api) AddProductCart(w http.ResponseWriter, r *http.Request) { //! ожидается json с новым пользователем

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Error("error", err)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
		logger.Error("error", err)
		return
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("error", err)
		return
	}

	_, err = user.Cart.AddProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		logger.Error("error", err)
		return
	}

	//! меняем базу всех продуктов
	productOrigin, err := api.productStorage.GetProduct(product.ID)
	if err != nil {
		http.Error(w, `"error":"product not found"`, 404)
		logger.Error("error", err)
		return
	}

	if productOrigin.Quantity-product.Quantity < 0 {
		http.Error(w, `{"error":"product is missing"}`, 400)
		logger.Error("error", err)
		return
	}

	productOrigin.Quantity = productOrigin.Quantity - product.Quantity
	api.productStorage.ChangeProduct(productOrigin)

	http.Redirect(w, r, "/cart", 200)
}

// TODO если изменено количество, то уменьшить количество.
func (api *Api) ChangeProductCart(w http.ResponseWriter, r *http.Request) { //! получаем json полного пользователя

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Error("error", err)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
		logger.Error("error", err)
		return
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("error", err)
		return
	}

	_, err = user.Cart.ChangeProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		logger.Error("error", err)
		return
	}
	w.Write(body)
}

func (api *Api) DeleteProductCart(w http.ResponseWriter, r *http.Request) { //! получаем только id

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		logger.Error("error", err)
		return
	}

	product, err := user.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		logger.Error("error", err)
		return
	}

	product, err = user.Cart.DeleteProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		logger.Error("error", err)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error":"server error"`, 500)
		logger.Error("error", err)
		return
	}

	w.Write(resp)
}

func (api *Api) CommentProduct(w http.ResponseWriter, r *http.Request) { //! принимаем post json продукта и конкретный продукт
	//* проверяем метод
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"bad method"}`, 500)
		return
	}

	//* получаем конкретного юзера
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		logger.Error("error", err)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		logger.Error("error", err)
		return
	}

	//* получаем продукт
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"]) //TODO переделать другую переменную принимать
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		logger.Error("error", err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
		logger.Error("error", err)
		return
	}
	defer r.Body.Close()

	var comment storage.Comment

	err = json.Unmarshal(body, &comment)
	if err != nil {
		http.Error(w, `{"error":"json error"}`, 500)
		logger.Error("error", err)
		return
	}

	//TODO бизнес логика
	product, err := user.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, 500)
		logger.Error("error", err)
		return
	}

	newProduct := storage.Product{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
		About:    product.About,
		Comments: append(product.Comments, &comment),
	}

	api.productStorage.ChangeProduct(newProduct)

	http.Redirect(w, r, "/cart/comments", 400)
}
