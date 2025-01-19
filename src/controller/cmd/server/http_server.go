package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"com.kvs/deviceorbit/controller/cmd/k8s"
	"github.com/gorilla/mux"
)

type httpServer struct {
	podController *k8s.PodController

	server *http.Server
}

func NewHttpServer(podController *k8s.PodController) *httpServer {
	return &httpServer{podController: podController}
}

func (httpServer *httpServer) StartHttpServer(port string) error {
	muxer := mux.NewRouter()
	muxer.HandleFunc("/pods", httpServer.handlePod).Methods("POST", "DELETE")
	muxer.HandleFunc("/health", httpServer.handleHealth).Methods("GET")

	httpServer.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: muxer,
	}

	return httpServer.server.ListenAndServe()
}

func (httpServer *httpServer) ShutdownHttpServer(ctx context.Context) error {
	if httpServer.server != nil {
		err := httpServer.server.Shutdown(ctx)
		if err != nil {
			return err
		}
		httpServer.server = nil
	}

	return nil
}

func (httpServer *httpServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := httpServer.podController.ValidateRole()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (httpServer *httpServer) handlePod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body := struct {
			DeviceSerial string `json:"deviceSerial"`
			Platform     int    `json:"platform"`
		}{}

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		httpServer.podController.CreatePod(body.DeviceSerial, body.Platform)
		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		body := struct {
			DeviceSerial string `json:"deviceSerial"`
		}{}

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		httpServer.podController.DeletePod(body.DeviceSerial)
		w.WriteHeader(http.StatusOK)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
