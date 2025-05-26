package controllers

import (
	"encoding/json"
	"inventory-control-hub/database"
	"inventory-control-hub/models"
	"inventory-control-hub/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func HomeRoute(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Welcome to inventory hub's home route...")
}
func GetProduct(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter is used to send back the response to the browser

	// r *http.Request represents the incoming request from the user contains method(get, post ..), header, params

	var products []models.Product // here product is a variable which is slice type [], which contains Product like struct which is located in models package.
	result := database.DB.Find(&products)

	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}
	// here database.DB, DB is a global variable .Find() tells the database to find all the records present in the products table and store them in products slice

	//1. Taking the struct name (e.g., Product)
	//2. Converting it to snake_case (if needed)
	//3. Pluralizing it → products

	utils.RespondWithJSON(w, http.StatusOK, products)
}
func GetProductById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // uses Gorilla mux router to get path variables, mux.Vars(r) returns a map[string] string, means key is string and value is also string,
	id := params["id"]    //  and then we can fetch value of key "id"

	var product models.Product

	result := database.DB.First(&product, id) // first record with matching id
	if result.Error != nil {
		// w.WriteHeader(http.StatusNotFound)
		// json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})

		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	//if found, return as a response

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(product)

	utils.RespondWithJSON(w, http.StatusOK, product)
}
func AddProduct(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	var product models.Product

	//decoding json into product structure
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	if product.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Product's name is required")
		return
	}

	//checking if the product with the same name already lies in db
	var existingProduct models.Product
	duplicate_rec := database.DB.Where("name=?", product.Name).First(&existingProduct)
	if duplicate_rec.Error == nil {
		utils.RespondWithError(w, http.StatusConflict, "This product already exists")
		return
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		// http.Error(w, "Failed to add product", http.StatusInternalServerError)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add product")
		return
	}

	json.NewEncoder(w).Encode(product)
}
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var product models.Product
	result := database.DB.First(&product, id)

	if result.Error != nil {

		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	//if product exists
	var updatedData models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		// http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json data")
		return
	}
	var existingProduct models.Product
	duplicate_rec := database.DB.Where("name=?", updatedData.Name).First(&existingProduct)
	if duplicate_rec.Error == nil {
		utils.RespondWithError(w, http.StatusConflict, "This product already exists")
		return
	}

	if updatedData.Name != "" {
		product.Name = updatedData.Name
	}

	if updatedData.Price != 0 {
		product.Price = updatedData.Price
	}

	if updatedData.Quantity != 0 {
		product.Quantity = updatedData.Quantity
	}

	if updatedData.Description != "" {
		product.Description = updatedData.Description
	}

	res := database.DB.Save(&product)

	if res.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, product)
}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	param := mux.Vars(r)
	id := param["id"]
	var count int64

	data := database.DB.First(&product, id)
	if data.Error != nil {
		// w.WriteHeader(http.StatusNotFound)
		// json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	//if any sales order is associated with the product
	database.DB.Model(&models.SalesOrder{}).Where("product_id = ?", product.ID).Count(&count)
	if count > 0 {
		utils.RespondWithError(w, http.StatusConflict, "Cannot delete product because sales orders exist")
		return
	}
	deleteResult := database.DB.Delete(&product)
	if deleteResult.Error != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete product"})
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	// json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})

}
