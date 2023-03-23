package main

import (
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/infrastructure/db"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/infrastructure/provider"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/infrastructure/webservice"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/controller"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/presenter"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/repository"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/usecase/interactor"
	httpPkg "github.com/CanobbioE/strict-clean-arch-go-webservice/pkg/http"
	"log"
	"net/http"
)

func main() {
	database := db.NewInMemoryDB()
	fileProvider := provider.NewDbFileProvider(database)
	fileRepo := repository.NewFilesRepository(fileProvider)

	wh := &webservice.Handler{File: &webservice.File{
		Controller: &controller.File{},
		Interactor: &interactor.File{Repo: fileRepo},
		Presenter:  &presenter.File{},
	}}

	r := httpPkg.NewRouter()
	wh.RegisterEndpoints(r)

	log.Fatal(http.ListenAndServe(":8080", r.Handler()))

}
