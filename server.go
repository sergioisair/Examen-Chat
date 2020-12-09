package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
)

var listaMensajes list.List
var listaClientes list.List

func main() {
	go servidor()
	for {
		fmt.Println("SERVIDOR DE CHAT ACTIVO")
		for i := listaMensajes.Front(); i != nil; i = i.Next() {
			fmt.Println(i.Value)
		}
		fmt.Scanln()
	}
}

func servidor() {
	serv, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		canal, err := serv.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleCliente(canal)
	}
}

func handleCliente(canal net.Conn) {
	var opcion uint64
	var err error
	var mensaje string
	err = gob.NewDecoder(canal).Decode(&mensaje)
	if err != nil {
		fmt.Println(err)
	}
	listaClientes.PushBack(canal)
	mensaje = "* " + mensaje + " se ha conectado *"
	fmt.Println(mensaje)
	listaMensajes.PushBack(mensaje)
	for {
		err = gob.NewDecoder(canal).Decode(&opcion)
		if err != nil {

			continue
		}
		if opcion == 1 { // ENVIAR MENSAJE
			err = gob.NewDecoder(canal).Decode(&mensaje)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(mensaje)
			listaMensajes.PushBack(mensaje)
			for i := listaClientes.Front(); i != nil; i = i.Next() {
				if i.Value.(net.Conn) != canal {
					err = gob.NewEncoder(i.Value.(net.Conn)).Encode(opcion)
					if err != nil {
						continue
					}
					err = gob.NewEncoder(i.Value.(net.Conn)).Encode(mensaje)
					if err != nil {
						continue
					}
				}
			}
		}
	}
}
