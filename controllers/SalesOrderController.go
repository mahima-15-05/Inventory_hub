package controllers

import (
	"encoding/json"
	"inventory-control-hub/database"
	"inventory-control-hub/models"
	"inventory-control-hub/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetSalesOrder(w http.ResponseWriter, r *http.Request) {
	var salesOrder []models.SalesOrder

	result := database.DB.Preload("Product").Find(&salesOrder)

	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Sales order not found")
		return
	}

	json.NewEncoder(w).Encode(salesOrder)

}

func CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	var salesOrder models.SalesOrder
	err := json.NewDecoder(r.Body).Decode(&salesOrder)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Input")
		return
	}
	if salesOrder.ProductID==nil{
		utils.RespondWithError(w, http.StatusBadRequest, "ProductID is required")
		return
	}

	var product models.Product

	result := database.DB.First(&product, salesOrder.ProductID)

	//checking if the product whose order has come exists or not
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	if salesOrder.Quantity <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	// if product is available-check the stock of the product
	if product.Quantity < salesOrder.Quantity {
		utils.RespondWithError(w, http.StatusBadRequest, "Product out of stock")
		return
	}

	// price and order-date
	salesOrder.TotalPrice = float64(salesOrder.Quantity) * product.Price
	salesOrder.OrderDate = time.Now()

	tx := database.DB.Begin()
	res := tx.Create(&salesOrder)

	if res.Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create sales order")
		return
	}

	// update product quantity
	product.Quantity = product.Quantity - salesOrder.Quantity

	if tx.Save(product).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product quantity")
		return

	}

	tx.Commit()
	if err := database.DB.Preload("Product").First(&salesOrder, salesOrder.ID).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Sales order not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, salesOrder)

}

func UpdateSalesOrder(w http.ResponseWriter, r *http.Request) {
	// fetch id from path
	params := mux.Vars(r)

	id := params["id"]
	var salesOrder models.SalesOrder

	err := json.NewDecoder(r.Body).Decode(&salesOrder)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// existing record

	var existingOrder models.SalesOrder

	if database.DB.First(&existingOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Sales order not found")
		return
	}

	var product models.Product
	if database.DB.First(&product, existingOrder.ProductID).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "The ordered product not found")
		return
	}

	if salesOrder.ProductID != nil {
		if salesOrder.ProductID != existingOrder.ProductID {
			utils.RespondWithError(w, http.StatusConflict, "The sales order cannot be updated for new product, please cancel the order and place another order again ")
		 	return
		}
	}

	// if product id is updated

	product.Quantity = product.Quantity + existingOrder.Quantity

	//insufficient
	if salesOrder.Quantity > product.Quantity {
		utils.RespondWithError(w, http.StatusBadRequest, "Insufficient stock for the updated quantity")
		return
	}

	product.Quantity = product.Quantity - salesOrder.Quantity

	//update order
	existingOrder.Quantity = salesOrder.Quantity
	existingOrder.TotalPrice = float64(salesOrder.Quantity) * product.Price
	existingOrder.OrderDate = time.Now()

	tx := database.DB.Begin()
	if err := tx.Save(&existingOrder).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product stock")
		return
	}
	tx.Commit()
	var updatedOrder models.SalesOrder
	if database.DB.Preload("Product").First(&updatedOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Updated sales order not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updatedOrder)

}

func DeleteSalesOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var salesOrder models.SalesOrder

	if database.DB.First(&salesOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Sales order not found")
		return
	}
	// finding the associates product
	var product models.Product
	if database.DB.First(&product, salesOrder.ProductID).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
	}

	product.Quantity = product.Quantity + salesOrder.Quantity

	tx := database.DB.Begin()
	if database.DB.Save(&product).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product stock")
		return
	}

	if tx.Delete(&salesOrder).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete sales order")
		return
	}

	tx.Commit()
	utils.RespondWithJSON(w, http.StatusOK, "Sales order deleted successfully")

}
