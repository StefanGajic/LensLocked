package main

import (
	"flag"
	"fmt"
	"net/http"

	//"github.com/lenslocked/email"
	"github.com/lenslocked/middleware"
	"github.com/lenslocked/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/lenslocked/controllers"
	"github.com/lenslocked/models"
)

func main() {

	boolPtr := flag.Bool("prod", false, "Provide this flag "+
		"in production. This ensures that a .config file is "+
		"provided before the application starts.")
	flag.Parse()
	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	// mgCfg := cfg.Mailgun
	// emailer := email.NewClient(
	// 	email.WithSender("Lenslocked.com Support", "support@"+mgCfg.Domain),
	// 	email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	// )

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User) //, emailer
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}

	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").
		Name(controllers.IndexGalleries)
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").
		Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").
		Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))

}
