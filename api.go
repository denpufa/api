package api

import (
	"math"
	"encoding/json"
	"log"
	"net/http"
)

// Corpos de requisições
type PoupancaRequestBody struct {
	Value float64 `json:"value"`
	Years int `json:"years"`
}

type TradeRequestBody struct {
	ValueInitial float64 `json:"value_initial"`
	ValueFinal   float64 `json:"value_final"`
	Days         int     `json:"days"`
}

type TesouroRequestBody struct {
	ValueInitial float64 `json:"value_initial"`
	Years        int     `json:"years"`
}

type JurosRequestBody struct {
	Juros float64 `json:"juros"`

}

// Respostas das requisições
type ImpostoResponse struct {
	Descricao string `json:"descricao"`
}

type PoupancaResponse struct {
	Descricao string `json:"descricao"`
	ValorFinal float64 `json:"valor_final"`
	ValorInicial float64 `json:"valor_inicial"`

}

type TradeResponse struct {
	TipoNegociacao string `json:"tipo_negociacao"`
	ImpostoParaPagar float64 `json:"imposto_para_pagar"`
}

type TesouroResponse struct {
	Imposto      float64 `json:"imposto"`
	ValorFinal   float64 `json:"valor_final"`
}

var jurosAnual = 10.0

// Handler's
func PoupancaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var reqBody PoupancaRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Corpo de requisição inválido", http.StatusBadRequest)
		return
	}
	// juros compostos
	valorFinal := reqBody.Value*(math.Pow((1 + jurosAnual*0.7/100),float64(reqBody.Years)))

	res := PoupancaResponse{
		Descricao:       "Rendimento da poupança é isento de imposto e representa 70% do juros base",
		ValorFinal:    valorFinal,
		ValorInicial:  reqBody.Value,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func ImpostoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	tipo := r.URL.Path[len("/imposto/"):]
	res := ImpostoResponse{
		Descricao: "",
	}

	switch tipo {
	case "poupanca":
		res.Descricao = "Poupança é isenta de imposto"
	case "tesouro":
		res.Descricao = "Impostos sobre os lucros, entre 0 a 6 meses o imposto é de 22,5%, a partir de 6 meses até 1 ano, fica 20%, de 1 ano até 2 meses fica 17,5%, a partir de 2 anos fica 15%"
	case "daytrade":
		res.Descricao = "Imposto de 20% sobre lucros (caso tenha), para negociações diárias"
	case "swingtrade":
		res.Descricao = "Imposto de 15% sobre lucros (caso tenha), para negociações semanais"
	case "long":
		res.Descricao = "Imposto de 15% sobre lucros caso vendas acima de 20.000 mensais, para negociações com durações maiores de 1 mês"
	default:
		http.Error(w,"Tipo inválido", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func TradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var reqBody TradeRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Corpo de requisição inválido", http.StatusBadRequest)
		return
	}

	tipoNegociacao := ""
	imposto := 0.0

	if reqBody.Days == 1 {
		tipoNegociacao = "daytrading"

		if reqBody.ValueFinal - reqBody.ValueInitial > 0 {
			imposto = (reqBody.ValueFinal - reqBody.ValueInitial)*0.2
		}
	} else if reqBody.Days > 1 && reqBody.Days <= 30 {
		tipoNegociacao = "swingtrade"
		if reqBody.ValueFinal - reqBody.ValueInitial > 0 {
			imposto = (reqBody.ValueFinal - reqBody.ValueInitial)*0.15
		}
	} else {
		tipoNegociacao = "long"
		if reqBody.ValueFinal - reqBody.ValueInitial > 0 && reqBody.ValueFinal > 20.000 {
			imposto = (reqBody.ValueFinal - reqBody.ValueInitial)*0.15
		}
	}

	res := TradeResponse{
		TipoNegociacao: tipoNegociacao,
		ImpostoParaPagar: imposto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func TesouroHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var reqBody TesouroRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Corpo de requisição inválido", http.StatusBadRequest)
		return
	}

	imposto := 0.0
	rendimento := 0.0
	if reqBody.Years == 1 {
		imposto = float64(reqBody.ValueInitial) * 0.175
		rendimento = float64(reqBody.ValueInitial) * (math.Pow((1 + jurosAnual/100),float64(reqBody.Years)))

	} else {
		imposto = float64(reqBody.ValueInitial) * 0.15
		rendimento = float64(reqBody.ValueInitial) * (math.Pow((1 + jurosAnual/100),float64(reqBody.Years)))
	}

	res := TesouroResponse{
		Imposto:      imposto,
		ValorFinal:   rendimento,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func JurosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var reqBody JurosRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Corpo de requisição inválido", http.StatusBadRequest)
		return
	}

	jurosAnual = reqBody.Juros

	w.WriteHeader(http.StatusOK)
}

func StartApi() {
	http.HandleFunc("/poupanca", PoupancaHandler)
	http.HandleFunc("/imposto/", ImpostoHandler)
	http.HandleFunc("/trade", TradeHandler)
	http.HandleFunc("/tesouro", TesouroHandler)
	http.HandleFunc("/juros", JurosHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
