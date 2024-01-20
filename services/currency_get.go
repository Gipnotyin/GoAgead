package services

import (
	"bytes"
	"encoding/xml"
	"example.com/m/config"
	"example.com/m/entity"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

type services struct {
	app    fiber.Router
	config *config.Config
	client *resty.Client
}

type Services interface {
	Ping(ctx *fiber.Ctx) error
	GetCurrency(ctx *fiber.Ctx) error
}

func NewServices(app fiber.Router, configLoad *config.Config) *services {
	return &services{
		app:    app,
		config: configLoad,
		client: resty.New(),
	}
}

func (s *services) RegisterRouters() {
	s.app.Get("/ping", s.Ping)
	s.app.Get("/currency/:date<regex(\\d{2}-\\d{2}-\\d{4})>/:val<len(6)>", s.GetCurrency)
}

func (s *services) Ping(ctx *fiber.Ctx) error {
	return ctx.SendString("ok")
}

func (s *services) GetCurrency(ctx *fiber.Ctx) error {
	bytesXML, err := getValue(ctx.Params("date"), ctx.Params("val"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(err.Error())
	}

	result, err := decoder(bytesXML)
	if err != nil {
		return ctx.JSON(err.Error())
	}

	response := new(entity.ResponseCurrency)
	response.Value = result.Record[0].Value
	return ctx.JSON(response)
}

func getValue(date, val string) ([]byte, error) {
	client := &http.Client{}
	date = strings.Join(strings.Split(date, "-"), "/")
	URL := fmt.Sprintf("https://cbr.ru/scripts/XML_dynamic.asp?date_req1=%v&date_req2=%v&VAL_NM_RQ=%v", date, date, val)
	fmt.Println(URL)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, fmt.Errorf("Read body: %v", err)
		}
		return data, nil
	default:
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
}

func decoder(value []byte) (*entity.XMLAnswer, error) {
	var result entity.XMLAnswer

	decoder := xml.NewDecoder(bytes.NewReader(value))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
