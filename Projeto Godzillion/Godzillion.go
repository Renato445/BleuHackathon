package main


import (
	"fmt"
	"encoding/json"
	"crypto/sha512"
	"encoding/hex"
    "crypto/hmac"
	"io/ioutil"
	"net/http"
	"strings"
    "log"
    //"unicode/utf8"
    "os/exec"

)

//FUNCAO OK
func verSaldos(apikey string, apisecret string) bool{
	
	url 	  := "https://bleutrade.com/api/v2/account/getbalances?apikey="

	values := []string{}
    values = append(values, url)
    values = append(values, apikey)
    sign  := strings.Join(values, "")
	  
    //gera o hash e inclui na url
    sign = hashHmac512(sign, apisecret)
    
    values = append(values, "&apisign=")
    values = append(values, sign)

    result := strings.Join(values, "")
	
	//////////////////////////////// GJSON
	resp, err := http.Get(result)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	//le o html da pagina	
	html, err := ioutil.ReadAll(resp.Body)	
	////////////////////////////////
	var b bleuStruct

	jsonStr := json.Unmarshal(html, &b)
	if jsonStr != nil {
		fmt.Println("Erro:", jsonStr)
		return false
	}

	fmt.Printf("Seus saldos são os seguintes:\n")
	for i:= 0 ; i < 45; i++{
		fmt.Printf("[%s\t]:%s\t   | [Valor Pendente]:%s\n", b.Result[i].Currency, b.Result[i].Balance, b.Result[i].Pending)
		
	}
	
	return true

}

/// FUNCAO OK: NAO DA PRA SACAR POR CAUSA DO VALOR DO FEE = VALOR DO BALANCE ;-;
func sacar(apikey string, apisecret string, moeda string, valor string, carteira string) bool{

	url 	  := "https://bleutrade.com/api/v2/account/withdraw?apikey="

	values := []string{}
    values = append(values, url)
    values = append(values, apikey)
    values = append(values, "&currency=")
    values = append(values, moeda)
    values = append(values, "&quantity=")
    values = append(values, valor)
    values = append(values, "&address=")
    values = append(values, carteira)

    sign  := strings.Join(values, "")
	  
    //gera o hash e inclui na url
    sign = hashHmac512(sign, apisecret)
    
    values = append(values, "&apisign=")
    values = append(values, sign)

    result := strings.Join(values, "")
	
	//////////////////////////////// GJSON
	resp, err := http.Get(result)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	//le o html da pagina	
	html, err := ioutil.ReadAll(resp.Body)	
	////////////////////////////////
	var b bleuStruct

	jsonStr := json.Unmarshal(html, &b)
	if jsonStr != nil {
		fmt.Println("Erro:", jsonStr)
		return false
	}

	fmt.Printf("Parabens seu saque foi realizado com sucesso!\n%s\n", b)
	
	return true
	
}

///FUNCAO OK
func depositosRealizados(apikey string, apisecret string) bool{

	url 	  := "https://bleutrade.com/api/v2/account/getdeposithistory?apikey="

	values := []string{}
    values = append(values, url)
    values = append(values, apikey)
    sign  := strings.Join(values, "")
	  
    //gera o hash e inclui na url
    sign = hashHmac512(sign, apisecret)
    
    values = append(values, "&apisign=")
    values = append(values, sign)

    result := strings.Join(values, "")
	
	//////////////////////////////// GJSON
	resp, err := http.Get(result)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	//le o html da pagina	
	html, err := ioutil.ReadAll(resp.Body)	
	////////////////////////////////
	var b bleuStruct

	jsonStr := json.Unmarshal(html, &b)
	if jsonStr != nil {
		fmt.Println("Erro:", jsonStr)
		return false
	}
	if b.Result != nil {
		fmt.Printf("Seus depositos realizados foram os seguintes:\n")
		i := 0
		for ;i < 1; i++{
			fmt.Printf("%s", b.Result[i])
		}
	}else{
		fmt.Printf("Voce nao realizou nenhum deposito ainda.\n")
	}
		
	
	
	return true
	
}

func funcaoTopper(traduzir string) bool{

	values := []string{}
    values = append(values, "echo ")
    values = append(values, "\"")
    values = append(values, traduzir)
    values = append(values, "\"")
    values = append(values, " >> saida.txt")
    criaArq  := strings.Join(values, "")

    fmt.Printf("-------------------------------------------------------------------\n\n%s\n\n-------------------------------------------------------------------", criaArq)


	cmd := exec.Command("bash", "-c", criaArq)
    _, err := cmd.CombinedOutput()

    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
   

    ////////////// AQUI CHAMO O TRAD EM PY
    //// ARQ DIVIDIR EM 4, POR UMA THREAD PARA CADA ARQUIVO  E DPS EXIBIR UM CAT NA ORDEM
    cmd = exec.Command("bash", "-c", "python3 tradutor.py")
    _, err = cmd.CombinedOutput()

    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
   
    //////////////


	return true
}

func hashHmac512(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}


type bleuStruct struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Result  []struct {
		Currency      string      `json:"Currency"`
		Balance       string      `json:"Balance"`
		Available     string      `json:"Available"`
		Pending       string      `json:"Pending"`
		CryptoAddress interface{} `json:"CryptoAddress"`
		IsActive      string      `json:"IsActive"`
		AllowDeposit  string      `json:"AllowDeposit"`
		AllowWithdraw string      `json:"AllowWithdraw"`
	} `json:"result"`
}

func verNoticia(url string) bool{
	
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	//le o html da pagina
	html, err := ioutil.ReadAll(resp.Body)
	//transforma de bytearray para string, e da um split na mensagem
	str := fmt.Sprintf("%s", html)
	str2 := strings.Split(str, "'inputdata'>0x")
	finalStr := strings.Split(str2[1], "<")

	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	//fmt.Printf("%s\n", str)  ok pego o html como string
	
	//fmt.Printf("%s\n", finalStr[0]) // aqui é o texto em hexadecimal
	
	src := []byte(finalStr[0])

	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(dst, src)
	
	if err != nil {
		log.Fatal(err)
	}
	///ao inves de criar o arquivo de texto com a saida, eu posso simplesmente separar em 4 arrays traduzir os 4 gerar 4 textos e dar cat neles

	myString := string(dst[:])
	funcaoTopper(myString)
	return true
}


func main() {
	
	///---------------Menu---------------///
	flag := true

	for(flag == true){
		fmt.Printf("\n///-----------------BLEU-----------------///\n///1-Ver saldos\t\t\n///2-Saque\t\t\n///3-Depositos\t\t\n///4-Comprar\t\t\n///5-Vender\t\t\n///6-Ver mensagem em uma transacao\t\t\n")
		opcode   := 0
		moeda    := ""
		valor    := ""
		carteira := ""
		url      := ""

		fmt.Scanf("%d", &opcode)

		switch opcode {
			case 1: 
				verSaldos("a127ab557b222a6586b4985d5ffb5735", "df30f5b970207aa3c0007a7d27be031ef86c8afa")

			case 2: 
				fmt.Printf("Qual moeda deseja sacar? Exemplo: bitcoin = BTC, doge = DOGE\n")
				fmt.Scanf("%s", &moeda)
				fmt.Printf("Qual valor deseja sacar? Exemplo: 0.00001\n")
				fmt.Scanf("%s", &valor)
				fmt.Printf("Para qual carteira de %s deseja enviar seu saque?\n", moeda)
				fmt.Scanf("%s", &carteira)
				sacar("a127ab557b222a6586b4985d5ffb5735", "df30f5b970207aa3c0007a7d27be031ef86c8afa", moeda, valor, carteira)

			case 3:
				depositosRealizados("a127ab557b222a6586b4985d5ffb5735", "df30f5b970207aa3c0007a7d27be031ef86c8afa")

			case 4:
				///Comprar

			case 5:
				///Vender

			case 6:	
				fmt.Printf("Opa, eh so nos passar o endereco da pagina que quer verificar se existe a mensagem\nPor exemplo: https://etherscan.io/tx/0xb1ed364e4333aae1da4a901d5231244ba6a35f9421d4607f7cb90d60bf45578a\n\n")
				fmt.Scanf("%s", &url)
				verNoticia(url)
			default: fmt.Println("oops... voce nao digitou uma opcao correta\n")
		}

	}




}