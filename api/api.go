package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"personae-fasti/data"
	"personae-fasti/opt"

	"github.com/rs/cors"
)

type APIFunc func(http.ResponseWriter, *http.Request) *APIError
type APIFuncAuth func(http.ResponseWriter, *http.Request, *data.Player) *APIError

type APIServer struct {
	server  *http.Server
	storage *data.Storage
}

type APIError struct {
	Error   error  `json:"-"`
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (api *APIServer) HandleError(e error) *APIError {
	return &APIError{
		Error: e,
		Code:  http.StatusInternalServerError,
	}
}

func (api *APIServer) HandleErrorString(e string) *APIError {
	return api.HandleError(fmt.Errorf(e))
}

func (a *APIError) WithCode(c int) *APIError {
	a.Code = c
	return a
}

func (a *APIError) WithMessage(m string) *APIError {
	a.Message = m
	return a
}

func ReadBody(r *http.Request) []byte {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func ReadJsonBody(r *http.Request, v any) error {
	bodyBytes := ReadBody(r)
	return json.Unmarshal(bodyBytes, v)
}

func (api *APIServer) Respond(r *http.Request, w http.ResponseWriter, status int, v any) *APIError {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	respErr := ""

	if APIErr, ok := v.(*APIError); ok {
		if APIErr.Message == "" {
			APIErr.Message = APIErr.Error.Error()
		}
		respErr = APIErr.Message
	}

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return api.HandleError(err)
	}

	jsonData, _ := json.Marshal(v)

	log := &data.Log{
		Time:     time.Now(),
		User:     0, //r.Context().Value("user").(string),
		URI:      r.RequestURI,
		Method:   r.Method,
		Request:  string(ReadBody(r)),
		Response: string(jsonData),
		Error:    respErr,
		HTTPCode: status,
	}

	api.storage.Log(log, r.Context())

	return nil
}

func InitServer(c *opt.Conf, s *data.Storage) *APIServer {

	router := http.NewServeMux()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "AccessKey"},
		AllowCredentials: true,
	})

	api := &APIServer{
		server: &http.Server{
			Addr:    c.App.Port,
			Handler: crs.Handler(router),
		},
		storage: s,
	}

	api.SetHandlers(router)

	log.Println("API server running on ", api.server.Addr)

	if err := api.server.ListenAndServe(); err != nil {
		panic(err)
	}

	return api

}
