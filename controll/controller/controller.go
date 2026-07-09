package controller

import (
	"context"
	"fmt"

	"github.com/coigo/micro-cloud/infra"
)

type NewContainer struct {
	ContainerSize ContainerSize
}

type ContainerSize int

const (
	SMALL ContainerSize = iota
	MEDIUM
	LARGE
)

func OnUpdate() {

}

func Decide (ctx context.Context) {
	

	iter := infra.Redis.Scan(ctx, 0,"machine-status:*", 0).Iterator()
		
	for iter.Next(ctx) {
	    key := iter.Val()
	
	    value, err := infra.Redis.Get(ctx, key).Result()
	    if err != nil {
	        continue
	    }

		// ver o que resta
		// pegar o tamanho do container
		// ver quantos containeres eu consigo rodar
		// pegar o que tem menos
		// se estiver parecido pegar o menor			
	
	    fmt.Println(key, value)

	}

	// eu tenho que buscar fechar uma maquina por inteiro? na minha visao isso parece o correto
	// pesos
	// pegar o menor peso?
	// desconsidera o que nao consegue rodar
	// se for um peso muito diferente pega esse, senao, pega o menor
}
