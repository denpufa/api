package api

import (
	"math"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPoupancaHandler(t *testing.T) {
	reqBody := PoupancaRequestBody{
		Value: 1000,
		Years: 2,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/poupanca", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PoupancaHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status esperado 200, mas obteve %d", rr.Code)
	}

	var res PoupancaResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	expectedValorInicial := 1000.0
	expectedValorFinal := expectedValorInicial * (math.Pow((1+0.7*10.0/100), float64(2)))

	if res.ValorFinal != expectedValorFinal {
		t.Errorf("Valor total esperado %.2f, mas obteve %.2f", expectedValorFinal, res.ValorFinal)
	}

	if res.ValorInicial != expectedValorInicial {
		t.Errorf("Valor inicial esperado %.2f, mas obteve %.2f", expectedValorInicial, res.ValorInicial)
	}
}

func TestImpostoHandler(t *testing.T) {
	tt := []struct {
		tipo             string
		expectedDescricao string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{tipo: "poupanca", expectedDescricao: "Poupança é isenta de imposto", expectedStatus: http.StatusOK},
		{tipo: "tesouro", expectedDescricao: "Impostos sobre os lucros, entre 0 a 6 meses o imposto é de 22,5%, a partir de 6 meses até 1 ano, fica 20%, de 1 ano até 2 meses fica 17,5%, a partir de 2 anos fica 15%", expectedStatus: http.StatusOK},
		{tipo: "daytrade", expectedDescricao: "Imposto de 20% sobre lucros (caso tenha), para negociações diárias", expectedStatus: http.StatusOK},
		{tipo: "swingtrade", expectedDescricao: "Imposto de 15% sobre lucros (caso tenha), para negociações semanais", expectedStatus: http.StatusOK},
		{tipo: "long", expectedDescricao: "Imposto de 15% sobre lucros caso vendas acima de 20.000 mensais, para negociações com durações maiores de 1 mês", expectedStatus: http.StatusOK},
		{tipo: "invalido", expectedDescricao: "", expectedStatus: http.StatusBadRequest},
	}

	for _, tc := range tt {
		req, err := http.NewRequest("GET", "/imposto/"+tc.tipo, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ImpostoHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != tc.expectedStatus {
			t.Errorf("Status esperado %d, mas obteve %d", tc.expectedStatus, rr.Code)
		}

		if tc.expectedStatus == http.StatusOK {
			var res ImpostoResponse
			err = json.NewDecoder(rr.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if res.Descricao != tc.expectedDescricao {
				t.Errorf("Descrição esperada %v, recebida %v ", tc.expectedDescricao, res.Descricao)
			}
		}
	}
}

func TestTradeHandler(t *testing.T) {
	reqBody := TradeRequestBody{
		ValueInitial: 1000,
		ValueFinal:   1200,
		Days:         30,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/trade", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TradeHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status esperado 200, mas obteve %d", rr.Code)
	}

	var res TradeResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	expectedTipoNegociacao := "swingtrade"

	if res.TipoNegociacao != expectedTipoNegociacao {
		t.Errorf("Tipo de negociação esperado %s, mas obteve %s", expectedTipoNegociacao, res.TipoNegociacao)
	}
}

func TestTesouroHandler(t *testing.T) {
	reqBody := TesouroRequestBody{
		ValueInitial: 1000,
		Years:        1,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/tesouro", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TesouroHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status esperado 200, mas obteve %d", rr.Code)
	}

	var res TesouroResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	expectedImposto := 175.0
	expectedValorFinal := 1100.0

	if res.Imposto != expectedImposto {
		t.Errorf("Imposto esperado %.2f, mas obteve %.2f", expectedImposto, res.Imposto)
	}

	if res.ValorFinal != expectedValorFinal {
		t.Errorf("Rendimento esperado %.2f, mas obteve %.2f", expectedValorFinal, res.ValorFinal)
	}
}

func TestJurosHandler(t *testing.T) {
	reqBody := JurosRequestBody{
		Juros: 12,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/juros", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(JurosHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status esperado 200, mas obteve %d", rr.Code)
	}

	if jurosAnual != reqBody.Juros {
		t.Errorf("Juros esperados %.2f, mas obteve %.2f", reqBody.Juros, jurosAnual)
	}
}
