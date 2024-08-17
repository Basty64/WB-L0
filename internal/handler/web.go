package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"wb/internal/cache"
	"wb/internal/db"
	"wb/internal/db/postgres"
)

type Server struct {
	server *http.Server
}

func NewServer(ctx context.Context, url string, database *postgres.RepoPostgres, inMemoryCache *cache.InMemoryCache) (*Server, error) {

	// Инициализация маршрутизатора
	mux := http.NewServeMux()

	// Инициализация шаблонов
	templates, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке шаблонов: %w", err)
	}

	// Инициализация сервера
	server := &http.Server{
		Addr:    url,
		Handler: mux,
	}

	// Настройка маршрутов
	mux.HandleFunc("/", handleIndex(templates))
	mux.HandleFunc("/order/{orderUID}", handleOrder(ctx, inMemoryCache, database))

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Start(urlServer string) error {

	log.Printf("Сервер запущен по адресу %d\n", urlServer)
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

func handleIndex(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {

			return
		}
	}
}

func handleOrder(ctx context.Context, InMemoryCache *cache.InMemoryCache, db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orderUIDstr := r.URL.Path[len("/order/"):]
		orderUID, err := strconv.Atoi(orderUIDstr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при загрузке страницы: %v", err), http.StatusInternalServerError)
		}

		// Поиск заказа в кэше
		order, err := InMemoryCache.GetOrder(orderUID)
		if err != nil {
			// Если заказ не найден в кэше, получить его из базы данных
			var err error
			order, err = db.GetOrder(ctx, orderUID)
			if err != nil {
				if errors.Is(err, cache.OrderNotFoundErr) {
					http.Error(w, fmt.Sprint("заказ отсутствует"), http.StatusNotFound)
					return
				} else {
					http.Error(w, fmt.Sprintf("Ошибка при получении данных заказа: %v", err), http.StatusInternalServerError)
					return
				}
			}
			// Добавить заказ в кэш
			err = InMemoryCache.InsertOrder(orderUID, order)
			if err != nil {
				fmt.Errorf("ошибка при добавлении заказа в кэш: %v", err)
			}
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
			return
		}
	}
}
