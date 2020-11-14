package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"
)

//Result type struct parece com uma interface, só uma estrutura de dado para retorna uma info desejada p/ ficar mais organizado
type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	//chamado no form action do checkout
	http.HandleFunc("/process", process)
	http.ListenAndServe(":8090", nil)
}

func home(writer http.ResponseWriter, request *http.Request) {
	temp := template.Must(template.ParseFiles("templates/home.html"))
	temp.Execute(writer, "")
}

func process(writer http.ResponseWriter, request *http.Request) {
	//demonstração de que o process está pegando a info do campo cupon e cartão de crédito
	//log.Println(request.FormValue("coupon"))
	//log.Println(request.FormValue("cc-number"))
	temp := template.Must(template.ParseFiles("templates/home.html"))
	temp.Execute(writer, "")
}

// retorna o status do cupon
func makeHTTPCall(urlMicroService string, coupon string, ccNumber string) Result {
	// empacotar os calores de cupon e ccNumber na mesma requisição, transformar em um obj
	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	// http.PostForm vai retornar o result, se for tudo ok e o err, caso algo dê errado
	response, err := http.PostForm(urlMicroService, values)
	// tratamento de erro em GO. nil == null. vai entender
	if err != nil {
		log.Fatal("Microservice b out")
	}

	// fecha a conexão depois de todo o cód da função for executado (não importa a posição do código na função, o defer sempre é executado por último)
	defer response.Body.Close()

	// data vai ser um json, com algum status - status: {"alguma coisa"}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error processing result", err)
	}

	result := Result{} // aquele type struct

	// Unmarshal transforma o JSON em formato struct e depois passa ele para esse result. Tá com o &, pq tá apontando, deve ser estilo C
	json.Unmarshal(data, &result)

	return result

}
