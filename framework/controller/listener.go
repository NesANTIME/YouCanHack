package controller

import (
	"log"
	"youcanhack/framework/controller/connection"
)

type SubOption = connection.SubOption
type MenuItem = connection.MenuItem

func Get_BDLab() []connection.MenuItem {
	menu, err := connection.Connection_DB()
	if err != nil {
		log.Fatalf("no se pudo cargar el menú: %v", err)
	}
	return menu
}
