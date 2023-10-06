package main

import (
	_ "embed"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	//go:embed sql/busca_pessoas.sql
	buscaPessoasSql string

	//go:embed sql/cria_pessoa.sql
	criaPessoaSql string
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgrespw@localhost:5432/banco_dos_amigos?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.GET("/ping", handlePing)

	r.GET("/pessoas", getPessoasHandler(db))
	r.POST("/pessoas", createPessoaHandler(db))

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type PessoaBuilder interface {
	WithNome(nome string) PessoaBuilder
	WithIdade(idade int) PessoaBuilder
	WithAltura(altura float64) PessoaBuilder
	WithPeso(peso float64) PessoaBuilder
	Build() Pessoa
}

type Pessoa struct {
	ID     string  `json:"id" db:"id"`
	Nome   string  `json:"nome" db:"nome"`
	Idade  int     `json:"idade" db:"idade"`
	Altura float64 `json:"altura" db:"altura"`
	Peso   float64 `json:"peso" db:"peso"`
}

func (p Pessoa) WithNome(nome string) PessoaBuilder {
	p.Nome = nome
	return p
}

func (p Pessoa) WithIdade(idade int) PessoaBuilder {
	p.Idade = idade
	return p
}

func (p Pessoa) WithAltura(altura float64) PessoaBuilder {
	p.Altura = altura
	return p
}

func (p Pessoa) WithPeso(peso float64) PessoaBuilder {
	p.Peso = peso
	return p
}

func (p Pessoa) Build() Pessoa {
	return p
}

func createPessoaHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// p := Pessoa{}.WithNome("Mateus")
		// fmt.Printf("%#v", p)
		var payload, pessoa Pessoa

		err := c.BindJSON(&payload)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = db.Get(&pessoa, criaPessoaSql, payload.Nome, payload.Idade, payload.Altura, payload.Peso)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": pessoa,
		})
	}
}

func getPessoasHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pessoa Pessoa

		err := db.Get(&pessoa, buscaPessoasSql)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": pessoa,
		})
	}
}
