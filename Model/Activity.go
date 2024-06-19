package Model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"ims/Database"
	"time"

	"gorm.io/gorm"
)

type Changetype string

const (
	Addition    Changetype = "addition"
	Subtraction Changetype = "subtraction"
	Insert      Changetype = "insert"
	Delete      Changetype = "delete"
	Update      Changetype = "update"
)

func (ct *Changetype) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ct = Changetype(v)
	case string:
		*ct = Changetype(v)
	default:
		return fmt.Errorf("cannot scan type %T into Changetype", v)
	}
	return nil
}

func (ct Changetype) Value() (driver.Value, error) {
	return string(ct), nil
}

type Activity struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Changetype   Changetype `json:"changetype" binding:"required" gorm:"type:Changetype"`
	ChangeAmount int        `json:"change_amount" binding:"required"`
	Timestamp    time.Time  `json:"timestamp" binding:"required" gorm:"default:CURRENT_TIMESTAMP"`
	ProductID    uint       `json:"product_id"`
	UserID       uint       `json:"user_id" `
	User         User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Product      Product    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ActivityResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	ProductID    uint      `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Changetype   string    `json:"changetype"`
	ChangeAmount int       `json:"change_amount"`
	Timestamp    time.Time `json:"timestamp"`
}

func (a *Activity) Save(tx *gorm.DB) (*Activity, error) {
	err := tx.Exec("INSERT INTO activities (changetype, change_amount, product_id, user_id) VALUES (?, ?, ?, ?)", a.Changetype, a.ChangeAmount, a.ProductID, a.UserID).Error
	if err != nil {
		return &Activity{}, err
	}
	return a, nil
}

func FilterActivity(query string) ([]ActivityResponse, error) {
	var activity []ActivityResponse
	err := Database.Database.Raw(query).Scan(&activity)

	fmt.Println(query)

	if err.Error != nil {
		return []ActivityResponse{}, err.Error
	}

	if err.RowsAffected == 0 {
		return []ActivityResponse{}, errors.New("activity not found")
	}

	return activity, nil
}

func GetActivityByID(id string) (ActivityResponse, error) {
	var activity ActivityResponse
	err := Database.Database.Raw(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE a.id = ?`, id).Scan(&activity)
	if err.Error != nil {
		return ActivityResponse{}, err.Error
	}

	if err.RowsAffected == 0 {
		return ActivityResponse{}, errors.New("activity not found")
	}

	return activity, nil
}

func ActivityQuery(changetype string, userName string, timestamp string, productName string) string {
	if changetype != "" && userName != "" && timestamp != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND DATE(a.timestamp) = '%s' AND p.name ILIKE '%%%s%%' AND a.changetype = '%s'`, userName, timestamp, productName, changetype)
	} else if changetype != "" && userName != "" && timestamp != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND DATE(a.timestamp) = '%s' AND a.changetype = '%s'`, userName, timestamp, changetype)
	} else if changetype != "" && userName != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND p.name ILIKE '%%%s%%' AND a.changetype = '%s'`, userName, productName, changetype)
	} else if changetype != "" && timestamp != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE DATE(a.timestamp) = '%s' AND p.name ILIKE '%%%s%%' AND a.changetype = '%s'`, timestamp, productName, changetype)
	} else if userName != "" && timestamp != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND DATE(a.timestamp) = '%s' AND p.name ILIKE '%%%s%%'`, userName, timestamp, productName)
	} else if changetype != "" && userName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND a.changetype = '%s'`, userName, changetype)
	} else if changetype != "" && timestamp != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE DATE(a.timestamp) = '%s' AND a.changetype = '%s'`, timestamp, changetype)
	} else if changetype != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE p.name ILIKE '%%%s%%' AND a.changetype = '%s'`, productName, changetype)
	} else if userName != "" && timestamp != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND DATE(a.timestamp) = '%s'`, userName, timestamp)
	} else if userName != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s' AND p.name ILIKE '%%%s%%'`, userName, productName)
	} else if timestamp != "" && productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE DATE(a.timestamp) = '%s' AND p.name ILIKE '%%%s%%'`, timestamp, productName)
	} else if changetype != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE a.changetype = '%s'`, changetype)
	} else if userName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE u.username = '%s'`, userName)
	} else if timestamp != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE DATE(a.timestamp) = '%s'`, timestamp)
	} else if productName != "" {
		return fmt.Sprintf(`SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id
		WHERE p.name ILIKE '%%%s%%'`, productName)
	}
	return `SELECT
		a.id, 
		a.user_id AS User_ID,
		u.username As User_Name, 
		a.product_id AS Product_ID, 
		p.name AS Product_Name,
		a.changetype,
		a.change_amount,
		a.timestamp
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN products p ON a.product_id = p.id`
}
