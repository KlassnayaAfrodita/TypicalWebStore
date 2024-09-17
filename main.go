package main

import (
	"./handlers/handlers"
	"net/http"
)

func main() {
	mux := mux.NewRouter()

	// cart := handlers.Cart{}
	user := handlers.NewUser{}
	cart := handlers.NewCart()
	products := handlers.NewProducts()
	http.HandleFunc("/cart", cart.GetAllProduct).Methods("GET")
	http.HandleFunc("/cart/{id=[0-9]+"}, cart.GetProduct).Methods("GET")
 	http.HandleFunc("/cart", cart.AddProduct).Methods("POST") //! передаем json с новым Product
	http.HandleFunc("/cart/{id=[0-9]+}", cart.ChangeProduct).Methods("PUT") //! передаем json с новым Product
	http.HandleFunc("/cart/{id=[0-9]+}", cart.DeleteProduct).Methods("DELETE") //! передаем только id в url4

	http.HandleFunc("/product", products.GetAllProduct).Methods("GET")
	
	
	http.HandleFunc("/auth", user.)

	http.ListenAndServe(":8080", nil)
}
