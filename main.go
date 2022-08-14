package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/order", createOrder).Methods("POST")
	router.HandleFunc("/orders", getAllOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", getOneOrder).Methods("GET")
	router.HandleFunc("/orders/{id}", updateOrder).Methods("PATCH")
	router.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type order struct {
	OrderID string `json:"OrderID"`
	Price   string `json:"Price"`
	Title   string `json:"Title"`
}

type allOrders []order

var orders = allOrders{
	{
		OrderID: "1",
		Price:   "1000",
		Title:   "Burger",
	}, {
		OrderID: "2",
		Price:   "1500",
		Title:   "pizza",
	}, {
		OrderID: "3",
		Price:   "500",
		Title:   "pepsi soda",
	}, {
		OrderID: "4",
		Price:   "700",
		Title:   "french fries",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder order
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter correct data")
	}

	json.Unmarshal(reqBody, &newOrder)
	orders = append(orders, newOrder)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newOrder)
}

func getOneOrder(w http.ResponseWriter, r *http.Request) {
	OrderID := mux.Vars(r)["id"]

	for _, singleOrder := range orders {
		if singleOrder.OrderID == OrderID {
			json.NewEncoder(w).Encode(singleOrder)
		}
	}
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(orders)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	OrderID := mux.Vars(r)["id"]
	var updatedOrder order

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter correct data")
	}
	json.Unmarshal(reqBody, &updatedOrder)

	for i, singleOrder := range orders {
		if singleOrder.OrderID == OrderID {
			singleOrder.Price = updatedOrder.Price
			singleOrder.Title = updatedOrder.Title
			orders[i] = singleOrder
			json.NewEncoder(w).Encode(singleOrder)
		}
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	OrderID := mux.Vars(r)["id"]

	for i, singleOrder := range orders {
		if singleOrder.OrderID == OrderID {
			orders = append(orders[:i], orders[i+1:]...)
			fmt.Fprintf(w, "The order with ID: %v has been deleted", OrderID)
		}
	}
}
