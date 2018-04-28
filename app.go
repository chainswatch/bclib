package main

import (
  "app/chains"
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  email       string  `gorm:"type:varchar(100);unique_index"`
  password    string
  ip_address  string
}

type Role struct {
  gorm.Model
  name       string  `gorm:"type:varchar(100);not null;unique"`
}

func main() {
  chains.ChainsWatcher()
}
