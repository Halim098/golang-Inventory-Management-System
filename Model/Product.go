package Model

import (
	"errors"
	"fmt"
	"ims/Database"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description" binding:"required" gorm:"type:text"`
	Price       int        `json:"price" binding:"required"`
	Quantity    int        `json:"quantity" binding:"required"`
	Category    string     `json:"category" binding:"required"`
	Code        string     `json:"code" binding:"required" gorm:"unique"`
	Activities  []Activity `gorm:"foreignKey:ProductID"`
}

type ProductInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
	Quantity    *int    `json:"quantity"`
	Changetype  *string `json:"changetype"`
	Category    *string `json:"category"`
	Code        *string `json:"code"`
}

func (p *Product) Save(tx *gorm.DB) (*Product, error) {
	err := tx.Exec("INSERT INTO products (name, description, price, quantity, category, code, created_at) VALUES (?, ?, ?, ?, ?, ?,?)", p.Name, p.Description, p.Price, p.Quantity, p.Category, p.Code, time.Now()).Error
	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

func (p *Product) Update(id uint, input ProductInput, tx *gorm.DB) (*Product, error) {
	checkUpdate(p, input)
	err := tx.Exec("UPDATE products SET name = ?, description = ?, price = ?, quantity = ?, category = ?, updated_at=? WHERE id = ?", p.Name, p.Description, p.Price, p.Quantity, p.Category, time.Now(), id).Error
	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

func DeleteProduct(tx *gorm.DB, id uint) error {
	err := tx.Exec("UPDATE products SET deleted_at = ? WHERE id = ?", time.Now(), id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFilterProduct(query string) ([]Product, error) {
	var product []Product
	err := Database.Database.Raw(query).Scan(&product)
	fmt.Println(query)
	if err.Error != nil {
		return []Product{}, err.Error
	}

	if err.RowsAffected == 0 {
		return []Product{}, errors.New("product not found")
	}

	return product, nil
}

func GetProductsByID(id uint) (Product, error) {
	var product Product
	err := Database.Database.Raw("SELECT * FROM products WHERE id = ? AND deleted_at is NULL", id).Scan(&product)
	if err.Error != nil {
		return Product{}, err.Error
	}

	if err.RowsAffected == 0 {
		return Product{}, errors.New("product not found")
	}

	return product, nil
}

func GetProductsByCode(code string, tx *gorm.DB) (Product, error) {
	var product Product
	err := tx.Raw("SELECT * FROM products WHERE code = ? AND deleted_at is NULL", code).Scan(&product)
	if err.Error != nil {
		return Product{}, err.Error
	}

	if err.RowsAffected == 0 {
		return Product{}, errors.New("product not found")
	}

	return product, nil
}

func ProductQuery(name string, category string) string {
	if name != "" && category != "" {
		return fmt.Sprintf("SELECT * FROM products WHERE name ILIKE '%%%s%%' AND LOWER(category) = LOWER('%s') AND deleted_at is NULL", name, category)
	} else if name != "" {
		return fmt.Sprintf("SELECT * FROM products WHERE name ILIKE '%%%s%%' AND deleted_at is NULL", name)
	} else if category != "" {
		return fmt.Sprintf("SELECT * FROM products WHERE LOWER(category) = '%s' AND deleted_at is NULL", category)
	}
	return "SELECT * FROM products WHERE deleted_at is NULL"
}

func checkUpdate(product *Product, input ProductInput) {
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Quantity != nil {
		product.Quantity = *input.Quantity
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.Code != nil {
		product.Code = *input.Code
	}
}
