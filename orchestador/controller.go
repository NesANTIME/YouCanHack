package core

import (
	"log"
	"youcanhack/orchestador/qemu"
)

func reader(content string) {

}

// funcion principal

func Running(path_recipe string) {
	path_isos := "/home/nesantime/Descargas/"
	// receta := "RECETA"

	log.Println("Buscando ISO file...")

	iso := path_isos + "alpine-standard-3.23.4-x86_64.iso"

	err := qemu.Execute_Lab("ISo", iso, "512", "display")
	if err != nil {
		log.Fatal(err)
	}
}
