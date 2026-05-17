package core

import (
	"log"
	"youcanhack/core/database"
)

type SubOption = database.SubOption
type MenuItem = database.MenuItem

func Get_BDLab() []database.MenuItem {
	menu, err := database.Connection_DB()
	if err != nil {
		log.Fatalf("no se pudo cargar el menú: %v", err)
	}
	return menu
}
