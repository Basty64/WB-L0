package handler

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
	"wb/internal/cache"
	"wb/internal/logs"
	"wb/internal/models"
)

type Server struct {
	server *http.Server
}

func NewServer(url string, inMemoryCache *cache.InMemoryCache) (*Server, error) {

	// Инициализация маршрутизатора
	mux := http.NewServeMux()

	// Инициализация сервера
	server := &http.Server{
		Addr:    url,
		Handler: mux,
	}

	// Настройка маршрутов
	mux.Handle("GET /order", logs.RequestLogger(handleOrder(inMemoryCache)))
	mux.HandleFunc("/", handleIndex("index"))

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Start(urlServer string) error {

	log.Printf("Сервер запущен по адресу %s\n", urlServer)
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

func handleIndex(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Инициализация шаблонов
		templates, err := template.ParseFiles("./frontend/" + filename + ".html")
		if err != nil {
			_ = fmt.Errorf("ошибка при загрузке шаблонов: %w", err)
		}

		err = templates.ExecuteTemplate(w, filename+".html", nil)
		if err != nil {
			_ = fmt.Errorf(err.Error())
			return
		}
	}
}

func handleOrder(InMemoryCache *cache.InMemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orderUIDstr := r.URL.Query().Get("orderUID")
		orderUID, err := strconv.Atoi(orderUIDstr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при загрузке страницы"), http.StatusBadRequest)
			//log.Errorf("%s", err)
			return
		}
		if orderUID == 0 {
			http.Error(w, fmt.Sprintf("Некорректный query параметр"), http.StatusBadRequest)
			log.Errorf("%s", err)
			return
		}

		// Поиск заказа в кэше
		order, ok := InMemoryCache.GetOrder(orderUID)

		if ok != true {
			str := models.OrderNotFoundError{Text: "Заказ отсутствует"}
			strjson, err := json.Marshal(str)
			if err != nil {
				log.Error(err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(strjson)
			if err != nil {
				log.Error(err)
			}
			return
		}

		// Отображение данных заказа
		data, err := json.Marshal(order)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при кодировании данных заказа: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			log.Error(err)
		}
	}
}
