package models

import "gorm.io/gorm"

// Product model info
// @Description Product Information
type Transaction struct {
	gorm.Model
	Id       uint     `json:"id" validate:"required"`
	CartId     int  `form:"productid" json:"productid" validate:"required"`
	UserId     int  `form:"userid" json:"userid" validate:"required"`
	Status     string  `form:"status" json:"status" gorm:"default:process"`
	Total    float64 `form:"total" json:"total" validate:"required"`
	Cart Cart `gorm:"foreignkey:CartId;references:Id"`
	User User `gorm:"foreignkey:UserId;references:Id"`
}

func ViewTransaction(db *gorm.DB, trans *[]Transaction, id int) (err error) {
	err = db.Where("user_id = ?", id).Preload("Cart").Preload("User").Find(trans).Error
	if err != nil {
		return err
	}
	return nil
}

func AddtoTransaction(db *gorm.DB, data *Transaction) (err error) {
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func FinishTransaction(db *gorm.DB, trans *Transaction, id int) (err error) {
	err = db.Model(trans).Where("id = ?", id).Update("status", "finish").Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteTransaction(db *gorm.DB, data *Transaction, id int) (err error) {
	err = db.Where("id=?", id).Delete(data).Error
	if err != nil {
		return err
	}
	return nil
}