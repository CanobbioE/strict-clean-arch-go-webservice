package test

//go:generate mockgen -package mocks -source ../domain/book.go -destination mocks/book_repository.go
//go:generate mockgen -package mocks -source ../infrastructure/webservice/handler.go -destination mocks/book_controller.go
//go:generate mockgen -package mocks -source ../interfaces/controller/book_http_controller.go -destination mocks/book_interactor.go BookInteractor
