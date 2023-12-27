package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

type Persona struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Id    int    `json:"id"`
}

var Contador = 1

func pegaInformacao(w http.ResponseWriter, r *http.Request) {
	username := ""
	email := ""

	username = r.FormValue("nome") // dado vindo do formulario arquivo index.html
	email = r.FormValue("email")   // dado vindo do formulario arquivo index.html

	abrirArquivo, err := os.ReadFile("./data/bd.json")

	if err != nil {
		fmt.Println("erro na abertura do arquivo")
	}
	Contador++
	var ArquivoStringFormatado []Persona
	err = json.Unmarshal(abrirArquivo, &ArquivoStringFormatado)
	Contador = len(ArquivoStringFormatado)

	if err != nil {
		fmt.Println("erro na conversão de json para string")
	}
	if username != "" {
		ArquivoStringFormatado = append(ArquivoStringFormatado, Persona{username, email, Contador})
		username = ""
	}

	//joga o nome e o email pegos no index para o array
	arquivoAtualizado, err := json.Marshal(ArquivoStringFormatado) // converte o arquivo de string para json
	if err != nil {
		panic(err)
	}
	// Sobrescreve o bd.json com o conteúdo atualizado
	err = os.WriteFile("./data/bd.json", arquivoAtualizado, 0755)

	if err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("./template/list.html")

	tmpl.ExecuteTemplate(w, "list.html", ArquivoStringFormatado)

}

func pegaHTTML(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./template/index.html")

	tmpl.ExecuteTemplate(w, "index.html", nil)

}

var controle = 0

func deletaRegistro(w http.ResponseWriter, r *http.Request) {

	var idRegistro = r.URL.Query().Get("nome")

	abrirArquivo, err := os.ReadFile("./data/bd.json")

	if err != nil {
		fmt.Println("erro na abertura do arquivo")
	}

	var arquivoJson []Persona
	err = json.Unmarshal(abrirArquivo, &arquivoJson)
	if err != nil {
		fmt.Println("erro na conversão de json para string")
	}

	var valorExcluido = idRegistro

	for index, str := range arquivoJson {

		if str.Nome == valorExcluido {

			arquivoJson = append(arquivoJson[:index], arquivoJson[index+1:]...)
			arquivoAtualizado, err := json.Marshal(arquivoJson)

			if err != nil {
				panic(err)
			}
			err = os.WriteFile("./data/bd.json", arquivoAtualizado, 0755)
			http.Redirect(w, r, "/mensagem", http.StatusSeeOther)
			break
		}

	}

}

func confirmaOperacao(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./template/mensagem.html")

	tmpl.ExecuteTemplate(w, "mensagem.html", nil)

}

func main() {
	fmt.Print("Servidor Rodando...")

	//Rotas
	http.HandleFunc("/login", pegaHTTML) // Carrega a list de login
	http.HandleFunc("/", pegaInformacao) // Pega Nome / Email
	http.HandleFunc("/deletar", deletaRegistro)
	http.HandleFunc("/mensagem", confirmaOperacao)

	http.ListenAndServe(":3000", nil)
}
