package providers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/tashima42/whatsapp-rastreio/data"
	"github.com/tashima42/whatsapp-rastreio/helpers"
	"io"
	"net/http"
	"strings"
)

type CorreiosProvider struct {
	BaseUrl string
	DB      *sql.DB
}

func (cp *CorreiosProvider) RegisterPackage(user data.User, object data.Object) {
	correiosObjects, err := cp.GetCorreiosObjects([]string{object.Code})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	err = object.CreateObject(cp.DB)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	err = user.AddObject(cp.DB, object)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	correiosObject := correiosObjects[0]
	event := correiosObject.Event[0].ToEvent()
	event.ObjectId = object.ID
	err = event.CreateEvent(cp.DB)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

func (cp *CorreiosProvider) GetCorreiosObjects(objects []string) ([]CorreiosObject, error) {
	body := fmt.Sprintf("%s%s%s",
		"<rastroObjeto><usuario></usuario><senha></senha><tipo>L</tipo><resultado>T</resultado><objetos>",
		strings.Join(objects, ""),
		"</objetos><lingua>101</lingua><token></token></rastroObjeto>",
	)

	res, err := cp.correiosRequest(http.MethodPost, "service/rest/rastro/rastroMobile", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	var rastro rastroResponse
	err = decoder.Decode(&rastro)
	if err != nil {
		return nil, err
	}
	if strings.Contains(rastro.Objeto[0].Category, "ERRO") {
		return nil, fmt.Errorf("Invalid Code")
	}
	return rastro.Objeto, nil
}

func (cp *CorreiosProvider) correiosRequest(method string, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", cp.BaseUrl, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/xml")
	client := &http.Client{}
	return client.Do(req)
}

type CorreiosObject struct {
	Number   string                `json:"numero"`
	Initials string                `json:"sigla"`
	Name     string                `json:"nome"`
	Category string                `json:"categoria"`
	Event    []CorreiosObjectEvent `json:"evento"`
}

type CorreiosObjectEvent struct {
	Type        string   `json:"tipo"`
	Status      string   `json:"status"`
	Date        string   `json:"data"`
	Hour        string   `json:"hora"`
	Creation    string   `json:"criacao"`
	Description string   `json:"descricao"`
	Details     string   `json:"detalhe,omitempty"`
	Receiver    struct{} `json:"recebedor,omitempty"`
	Unity       struct {
		Local     string          `json:"local"`
		Code      string          `json:"codigo"`
		City      string          `json:"cidade"`
		State     string          `json:"uf"`
		Sto       string          `json:"sto"`
		UnityType string          `json:"tipounidade"`
		Address   correiosAddress `json:"endereco"`
	} `json:"unidade"`
	DestinyZipCode string `json:"cepDestino"`
	KeepDeadline   string `json:"prazoGuarda"`
	WorkDays       string `json:"diasUteis"`
	SendedDate     string `json:"dataPostagem"`
	DetailsOEC     struct {
		DeliveryMan string `json:"carteiro"`
		District    string `json:"distrito"`
		List        string `json:"lista"`
		Unity       string `json:"unidade"`
	} `json:"detalheOEC,omitempty"`
	Destiny []struct {
		Place        string          `json:"local"`
		Code         string          `json:"codigo"`
		City         string          `json:"cidade"`
		Neighborhood string          `json:"bairro"`
		State        string          `json:"uf"`
		Address      correiosAddress `json:"endereco"`
	} `json:"destino,omitempty"`
	Postage struct {
		DestinyZipCode    string `json:"cepdestino"`
		Ar                string `json:"ar"`
		Mp                string `json:"mp"`
		Dh                string `json:"dh"`
		Weight            string `json:"peso"`
		Volume            string `json:"volume"`
		ProgrammedDate    string `json:"dataprogramada"`
		PostageDate       string `json:"datapostagem"`
		TreatmentDeadline string `json:"prazotratamento"`
	} `json:"postagem,omitempty"`
}

func (coe *CorreiosObjectEvent) ToEvent() data.Event {
	bodyByte, _ := json.Marshal(coe)
	hash := helpers.GetMD5Hash(string(bodyByte))
	event := data.Event{
		Hash: hash,
		Body: string(bodyByte),
	}
	return event
}

type correiosAddress struct {
	Code         string `json:"codigo"`
	ZipCode      string `json:"cep"`
	Street       string `json:"logradouro"`
	Number       string `json:"numero"`
	Locale       string `json:"localidade"`
	State        string `json:"uf"`
	Neighborhood string `json:"bairro"`
	Latitude     string `json:"latitude,omitempty"`
	Longitude    string `json:"longitude,omitempty"`
}

type rastroResponse struct {
	Versao     string           `json:"versao"`
	Quantidade string           `json:"quantidade"`
	Pesquisa   string           `json:"pesquisa"`
	Resultado  string           `json:"resultado"`
	Objeto     []CorreiosObject `json:"objeto"`
}
