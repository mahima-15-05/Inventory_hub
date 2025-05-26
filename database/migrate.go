// this file is used for updating the changes in model file to the database

package database

import "inventory-control-hub/models"

func Migrate() {
	DB.AutoMigrate(&models.Product{}, &models.PurchaseOrder{}, &models.SalesOrder{})
}
