package main

import (
  "app/chains"
)

type User struct {
  email       string  `gorm:"type:varchar(100);unique_index"`
  password    string
  ip_address  string
}

type Role struct {
  name       string  `gorm:"type:varchar(100);not null;unique"`
}

func main() {
  chains.ChainsWatcher()
}
