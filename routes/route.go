package routes

import (
	"inventory-control-hub/controllers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.HomeRoute).Methods("GET")

	r.HandleFunc("/products", controllers.GetProduct).Methods("GET")
	r.HandleFunc("/product/{id}", controllers.GetProductById).Methods("GET")
	r.HandleFunc("/add-product", controllers.AddProduct).Methods("POST")
	r.HandleFunc("/update-product/{id}", controllers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/delete-product/{id}", controllers.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/sales-order", controllers.GetSalesOrder).Methods("GET")
	r.HandleFunc("/add-sales-order", controllers.CreateSalesOrder).Methods("POST")
	r.HandleFunc("/update-sales-order/{id}", controllers.UpdateSalesOrder).Methods("PUT")
	r.HandleFunc("/delete-sales-order/{id}", controllers.DeleteSalesOrder).Methods("DELETE")

	r.HandleFunc("/purchase-orders", controllers.GetPurchaseOrder).Methods("GET")
	r.HandleFunc("/purchase-order/{id}", controllers.GetPurchaseOrderById).Methods("GET")
	r.HandleFunc("/add-purchase-order", controllers.CreatePurchaseOrder).Methods("POST")
	r.HandleFunc("/update-purchase-order/{id}", controllers.UpdatePurchaseOrder).Methods("PUT")
	r.HandleFunc("/delete-purchase-order/{id}", controllers.DeletePurchaseOrder).Methods("DELETE")

	return r
}
