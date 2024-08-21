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
	mux.HandleFunc("/", handleIndex("show"))

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
			http.Error(w, fmt.Sprintf("Ошибка при загрузке страницы: %v", err), http.StatusInternalServerError)
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

//кусок кода с валидацией

//		if err != nil {
//			// Если заказ не найден в кэше, получить его из базы данных
//			if errors.Is(err, cache.OrderNotFoundErr) {
//				order, err := db.GetOrder(ctx, orderUID)
//				if err != nil {
//					log.Printf("Ошибка при поиске заказа №%d: %s", orderUID, err)
//					http.Error(w, fmt.Sprint("Ошибка при получении данных заказа"), http.StatusInternalServerError)
//					return
//				}
//				if order.ID == 0 {
//					http.Error(w, fmt.Sprint("Заказ не найден"), http.StatusNotFound)
//					return
//				}
//				// Добавить заказ в кэш
//				err = InMemoryCache.InsertOrder(orderUID, order)
//				if err != nil {
//					log.Printf("ошибка при добавлении заказа в кэш: %v", err)
//				}
//			} else {
//				log.Printf("Ошибка при поиске заказа №%d: %s", orderUID, err)
//				http.Error(w, fmt.Sprint("Ошибка при получении данных заказа"), http.StatusInternalServerError)
//			}
//		}
