package api

import (
	"embed"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRouter(handler *apiHandler, UIFiles embed.FS) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/api/team", handler.TeamList()).Methods(http.MethodGet)
	router.Handle("/api/team", handler.TeamAdd()).Methods(http.MethodPost)
	router.Handle("/api/team/{teamId}/settings", handler.TeamSettingsGet()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/settings", handler.TeamSettingsEdit()).Methods(http.MethodPost)

	router.Handle("/api/team/{teamId}/users", handler.UserList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/users", handler.UserAdd()).Methods(http.MethodPost)
	router.Handle("/api/team/{teamId}/users/{userId}", handler.UserGet()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/users/{userId}", handler.UserEdit()).Methods(http.MethodPatch)
	router.Handle("/api/team/{teamId}/users/{userId}", handler.UserDelete()).Methods(http.MethodDelete)

	router.Handle("/api/team/{teamId}/tasks", handler.TaskList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/tasks", handler.TaskAdd()).Methods(http.MethodPost)
	router.Handle("/api/team/{teamId}/tasks/{taskId}", handler.TaskGet()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/tasks/{taskId}", handler.TaskEdit()).Methods(http.MethodPatch)
	router.Handle("/api/team/{teamId}/tasks/{taskId}", handler.TaskDelete()).Methods(http.MethodDelete)
	router.Handle("/api/team/{teamId}/tasks/{taskId}/history", handler.HistoryList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/tasks/{taskId}/run", handler.TaskRun()).Methods(http.MethodPost)

	router.Handle("/api/team/{teamId}/sendMessage", handler.MessageSend()).Methods(http.MethodPost)

	router.Handle("/api/team/{teamId}/receivers", handler.ReceiverList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/receivers", handler.ReceiverAdd()).Methods(http.MethodPost)
	router.Handle("/api/team/{teamId}/receivers/{receiverId}", handler.ReceiverGet()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/receivers/{receiverId}", handler.ReceiverEdit()).Methods(http.MethodPatch)
	router.Handle("/api/team/{teamId}/receivers/{receiverId}", handler.ReceiverDelete()).Methods(http.MethodDelete)
	router.Handle("/api/team/{teamId}/receive", handler.Receive())
	router.Use(mux.CORSMethodMiddleware(router))

	spaHandler := SpaHandler{StaticPath: "ui/build", IndexPath: "index.html", Files: UIFiles}
	router.PathPrefix("/").Handler(spaHandler)

	return router
}
