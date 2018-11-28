package routes

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"2018_2_YetAnotherGame/presentation/controllers"
	"2018_2_YetAnotherGame/presentation/middlewares"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

//var fooCount = prometheus.NewCounter(prometheus.CounterOpts{
//	Name: "foo_total",
//	Help: "Number of foo successfully processed.",
//})

//var (
//	Timings = prometheus.NewSummaryVec(prometheus.SummaryOpts{
//		Name: "method_tim",
//		Help:"fff",
//	},
//		[]string{"method"})
//	Counter = prometheus.NewCounterVec(prometheus.CounterOpts{
//		Name:"method_counter",
//		Help:"fffg",
//	},
//		[]string{"method"},
//	))



func Router(env *controllers.Environment) http.Handler {
	prometheus.MustRegister(env.Counter)
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/api/users/me", env.MeHandle).Methods("GET")
	routerAuth.HandleFunc("/api/users/me", env.UpdateHandle).Methods("POST")
	routerAuth.HandleFunc("/api/avatar", env.AvatarHandle).Methods("POST")
	routerAuth.HandleFunc("/api/session", env.LogOutHandle).Methods("DELETE")
	authHandler := middlewares.AuthMiddleware(routerAuth)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	//router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fooCount.Add(1)
	//	fmt.Fprintf(w, "foo_total increased")
	//})

	router.Handle("/api/users/me", authHandler)
	router.Handle("/api/session", authHandler).Methods("DELETE")
	router.Handle("/api/upload", authHandler)
	router.HandleFunc("/api/leaders", env.ScoreboardHandle).Methods("GET")

	router.HandleFunc("/api/session/new", env.RegistrationHandle).Methods("POST")
	router.HandleFunc("/api/session", env.LoginHandle).Methods("POST")
	router.HandleFunc("/api/vkauth", env.VKRegister)
	return router
}