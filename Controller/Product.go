package Controller

import (
	"fmt"
	"ims/Database"
	"ims/Model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
	var product Model.Product
	userid := c.MustGet("user_id")

	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := Database.Database.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	defer tx.Rollback()

	_, err = product.Save(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product, err = Model.GetProductsByCode(product.Code, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activities := Model.Activity{
		ProductID:    product.ID,
		Changetype:   "insert",
		ChangeAmount: product.Quantity,
		Timestamp:    time.Now(),
		UserID:       userid.(uint),
	}

	_, err = activities.Save(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
}

func UpdateProduct(c *gin.Context) {
	var productInput Model.ProductInput
	userid := c.MustGet("user_id")
	prodid := c.Param("id")
	productid, err := strconv.Atoi(prodid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.BindJSON(&productInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wau"})
		return
	}

	product, err := Model.GetProductsByID(uint(productid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := Database.Database.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	defer tx.Rollback()

	activities := Model.Activity{
		ProductID: uint(productid),
		Timestamp: time.Now(),
		UserID:    userid.(uint),
	}

	if productInput.Changetype != nil {
		activities.Changetype = Model.Changetype(*productInput.Changetype)
		if productInput.Quantity == nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be filled"})
			return
		}
		activities.ChangeAmount = *productInput.Quantity

		if *productInput.Changetype == "addition" {
			*productInput.Quantity = product.Quantity + *productInput.Quantity
			fmt.Println(*productInput.Quantity)
		} else {
			*productInput.Quantity = product.Quantity - *productInput.Quantity
			fmt.Println(*productInput.Quantity)
		}

		if *productInput.Quantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock Not Enough"})
		}
	} else {
		activities.Changetype = "update"
		activities.ChangeAmount = 0
	}

	_, err = product.Update(uint(productid), productInput, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = activities.Save(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProduct(c *gin.Context) {
	userid := c.MustGet("user_id")
	productid := c.Param("id")
	id, err := strconv.Atoi(productid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	tx := Database.Database.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	defer tx.Rollback()

	err = Model.DeleteProduct(tx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activities := Model.Activity{
		ProductID:    uint(id),
		Changetype:   "delete",
		ChangeAmount: 0,
		Timestamp:    time.Now(),
		UserID:       userid.(uint),
	}

	_, err = activities.Save(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func GetProductsByID(c *gin.Context) {
	code := c.Param("id")
	id, err := strconv.Atoi(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := Model.GetProductsByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func GetProducts(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")

	query := Model.ProductQuery(name, category)
	products, err := Model.GetFilterProduct(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
