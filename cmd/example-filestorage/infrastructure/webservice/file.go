package webservice

import (
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/controller"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/presenter"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/usecase/interactor"
	httpPkg "github.com/CanobbioE/strict-clean-arch-go-webservice/pkg/http"
	"net/http"
)

type Handler struct {
	File *File
}

type File struct {
	Controller controller.FileController
	Interactor interactor.FileInteractor
	Presenter  presenter.FilePresenter
}

func (h *Handler) RegisterEndpoints(router httpPkg.RouterInterface) {
	router.HandleFunc("/file", http.MethodPut, h.Upload)
	router.HandleFunc("/file/{id}", http.MethodGet, h.FindByID)
	router.HandleFunc("/file", http.MethodDelete, h.Delete)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := h.File.Controller.ParseRequestBody(ctx, r.Body)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	err = h.File.Controller.ValidateRequest(ctx, f)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	f, err = h.File.Interactor.Upload(ctx, f)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	httpPkg.EncodeResponse(w, f)
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := httpPkg.GetRouteVariable(r, "id")

	f, err := h.File.Interactor.Retrieve(ctx, id)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	f, err = h.File.Presenter.AnonymizeOwner(ctx, f)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	httpPkg.EncodeResponse(w, f)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, err := h.File.Controller.ParseRequestBody(ctx, r.Body)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	err = h.File.Controller.ValidateRequest(ctx, f)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	f, err = h.File.Interactor.Delete(ctx, *f.Id)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	f, err = h.File.Presenter.RemoveID(ctx, f)
	if err != nil {
		httpPkg.EncodeResponse(w, err)
		return
	}

	httpPkg.EncodeResponse(w, f)
}
