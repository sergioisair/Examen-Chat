package main

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

var nick string
var listaMensajes list.List

func main() {
	fmt.Println("BIENVENIDO AL CHAT")
	fmt.Print("Ingresa tu nick: ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	nick = s.Text()
	cliente()
}

func cliente() {
	canal, error := net.Dial("tcp", ":5000")
	if error != nil {
		fmt.Println(error)
		return
	}
	error = gob.NewEncoder(canal).Encode(nick)
	if error != nil {
		fmt.Println(error)
	}
	go handleServidor(canal)
	var op string

	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n  1) Enviar MSJ    |   2) Enviar Archivo")
		fmt.Println("  3) Mostrar MSJs  |   4) Cerrar Sesión ")
		s.Scan()
		op = s.Text()
		if op == "1" {
			var mensaje string
			fmt.Print("ENVIAR: ")
			s.Scan()
			mensaje = s.Text()
			var opc uint64 = 1
			error := gob.NewEncoder(canal).Encode(opc)
			if error != nil {
				fmt.Println(error)
			} else {
				listaMensajes.PushBack("Tu: " + mensaje)
				mensaje := nick + ": " + mensaje
				gob.NewEncoder(canal).Encode(mensaje)
			}
		} else if op == "3" {
			fmt.Println("CHAT")
			band = true
			for i := listaMensajes.Front(); i != nil; i = i.Next() {
				fmt.Println(i.Value)
			}
			s.Scan()
			band = false
		} else if op == "4" {
			var opc uint64 = 3
			error := gob.NewEncoder(canal).Encode(opc)
			if error != nil {
				fmt.Println(error)
			}
			break
		} else {
			fmt.Println("Incorrecto")
		}
	}
	canal.Close()
}

var band bool = false

func handleServidor(canal net.Conn) {
	var opc uint64
	for {
		error := gob.NewDecoder(canal).Decode(&opc)
		if error != nil {
			continue
		}
		var mensaje string
		if opc == 1 {
			error = gob.NewDecoder(canal).Decode(&mensaje)
			if error != nil {
				continue
			}
			listaMensajes.PushBack(mensaje)
			if band {
				fmt.Println(mensaje)
			} else {
				fmt.Println()
				for i := listaMensajes.Front(); i != nil; i = i.Next() {
					fmt.Println(i.Value)
				}
			}
		}
	}
}
