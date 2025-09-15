//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/fingo-martpedia/fingo-transaction/internal/controller"
	"github.com/fingo-martpedia/fingo-transaction/internal/interfaces"
	"github.com/fingo-martpedia/fingo-transaction/internal/repository"
	"github.com/fingo-martpedia/fingo-transaction/internal/services"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Dependency struct {
	TransactionController interfaces.ITransactionController
}

func provideDB() *gorm.DB {
	return helpers.DB
}

func InitDependency() Dependency {
	wire.Build(
		provideDB,

		repository.NewTransactionRepository,
		wire.Bind(new(interfaces.ITransactionRepository), new(*repository.TransactionRepository)),

		services.NewTransactionService,
		wire.Bind(new(interfaces.ITransactionService), new(*services.TransactionService)),

		controller.NewTransactionController,
		wire.Bind(new(interfaces.ITransactionController), new(*controller.TransactionController)),

		wire.Struct(new(Dependency), "*"),
	)
	return Dependency{}
}
