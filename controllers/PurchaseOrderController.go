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

func GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	//fetch all the records from the db

	var purchaseOrder []models.PurchaseOrder

	if database.DB.Preload("Product").Find(&purchaseOrder).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Purchase order not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, purchaseOrder)

}

func GetPurchaseOrderById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	var purchaseOrder models.PurchaseOrder

	if database.DB.Preload("Product").First(&purchaseOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Purchase order not found")
		return
	}

	

	utils.RespondWithJSON(w, http.StatusOK, purchaseOrder)

}

func CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var purchaseOrder models.PurchaseOrder

	err := json.NewDecoder(r.Body).Decode(&purchaseOrder)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if purchaseOrder.ProductID==nil{
		utils.RespondWithError(w, http.StatusBadRequest, "ProductID is required")
		return
	}

	// check the quantity of purchase order

	if purchaseOrder.Quantity <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Order quantity must be greater than 0")
		return
	}

	//transaction begins

	tx := database.DB.Begin()

	//find the associated product

	var product models.Product
	if database.DB.First(&product, purchaseOrder.ProductID).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Associated product is not found")
		return
	}

	product.Quantity = product.Quantity + purchaseOrder.Quantity

	if tx.Save(&product).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product quantity")
		return
	}

	//set price
	purchaseOrder.OrderDate = time.Now()

	//save purchaseOrder

	if tx.Save(&purchaseOrder).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create purchase order")
		return
	}

	tx.Commit()

	var fullOrder models.PurchaseOrder
	if err := database.DB.Preload("Product").First(&fullOrder, purchaseOrder.ID).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch full purchase order")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, fullOrder)

}

func UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var oldPurchaseOrder models.PurchaseOrder
	if database.DB.First(&oldPurchaseOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "The purchase order not found")
		return
	}

	var newPurchaseOrder models.PurchaseOrder

	//decode req body
	err := json.NewDecoder(r.Body).Decode(&newPurchaseOrder)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	//start transaction

	tx := database.DB.Begin()

	if newPurchaseOrder.ProductID != nil {
		if oldPurchaseOrder.ProductID != newPurchaseOrder.ProductID {
			utils.RespondWithError(w, http.StatusConflict, "Ordered product cannot be changed, please cancel the purchase order and place an order again")
			return
		}
	}

	var product models.Product

	if tx.First(&product, oldPurchaseOrder.ProductID).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	product.Quantity -= oldPurchaseOrder.Quantity
	product.Quantity += newPurchaseOrder.Quantity

	//save updated product
	if tx.Save(&product).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	//updating other fields

	oldPurchaseOrder.OrderDate = time.Now()

	if newPurchaseOrder.Supplier != "" {
		oldPurchaseOrder.Supplier = newPurchaseOrder.Supplier
	}

	if newPurchaseOrder.Quantity != 0 {
		oldPurchaseOrder.Quantity = newPurchaseOrder.Quantity
	}

	if tx.Save(&oldPurchaseOrder).Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update purchase order")
		return
	}

	tx.Commit()

	var fullOrder models.PurchaseOrder

	if err := database.DB.Preload("Product").First(&fullOrder, id).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch full purchase order")
		return
	}

	// if database.DB.First()

	utils.RespondWithJSON(w, http.StatusOK, fullOrder)

}

func DeletePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var purchaseOrder models.PurchaseOrder
	if database.DB.First(&purchaseOrder, id).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Purchase order not found")
		return
	}

	// need to update the product quantity
	var product models.Product
	if database.DB.First(&product, purchaseOrder.ProductID).Error != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	product.Quantity = product.Quantity - purchaseOrder.Quantity

	tx := database.DB.Begin()

	if tx.Save(&product).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	if tx.Delete(&purchaseOrder).Error != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete purchase order")
		return
	}

	tx.Commit()
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Purchase order deleted successfully"})
}
