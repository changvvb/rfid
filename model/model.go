package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type User struct {
	gorm.Model
	// ID       uint `gorm:"primary_key" form:"id"`
	UserName string
	Age      uint
	Vip      bool
	Money    uint
}

type Admin struct {
	gorm.Model
	AdminName     string
	AdminPassword string
}

func init() {
	db,err := opendb()
	defer db.Close()
	if err != nil {
		log.Println("modle:", err)
	}
	db.CreateTable(&User{})
	db.CreateTable(&Admin{})
	// db.AutoMigrate(&User{})
}

//add administritor to database
func AddAdmin(name string, pass string) error {
	if name == "" || pass == "" {
		return fmt.Errorf("Name or password is empty")
	}
	db, err := opendb()
	if err != nil {
		return err
	}
	fmt.Println("insert info: ", name, pass)
	db.Save(&Admin{AdminName: name, AdminPassword: pass})
	return nil
}

//Exam administritor name and password
func Exam(name string, password string) bool {
	db, err := opendb()
	if err != nil {
		return false
	}
	r := !db.Where("admin_name=? AND admin_password=?", name, password).Find(&Admin{}).RecordNotFound()
	log.Println(r)
	return r
}

// func Admins() []Admin {
//     db, err := opendb()
//     if err != nil {
//         log.Println("model.Admins", err)
//         return nil
//     }
// }

func opendb() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:changvvb@/rfid?charset=utf8&parseTime=True&loc=Local")
	return db, err
}
