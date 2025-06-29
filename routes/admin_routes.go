package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/dancankarani/safa/middleware"
	"github.com/dancankarani/safa/models"
	"github.com/gofiber/fiber/v2"
)

func SetAdminRoutes(app *fiber.App) {
	g := app.Group("/api/v1", middleware.JWTMiddleware)
	//stations
	stations := g.Group("/admin/stations")
	stations.Get("/", controllers.ReadAllStationsController)
	stations.Get("/sales", controllers.GetAllSalesInStationHandler)
	stations.Get("/expenses", controllers.GetStationExpensesHandler)
	stations.Get("/:id", controllers.ReadStationByIDController)
	stations.Patch("/:id", controllers.ReadStationByIDController)
	stations.Post("/", controllers.NewStationHandler)

	// suppliers
	suppliers := g.Group("/admin/suppliers")
	suppliers.Get("/", controllers.GetSuppliersHandler)
	suppliers.Get("/balance/:id", controllers.GetSupplierBalanceHandler)
	suppliers.Get("/balances", controllers.GetAllSupplierBalancesHandler)
	suppliers.Get("/:id", controllers.GetSupplierHandler)
	suppliers.Post("/", controllers.AddSupplierHandler)
	suppliers.Patch("/:id", controllers.UpdateSupplierHandler)
	suppliers.Delete("/:id", controllers.DeleteSupplierHandler)


	// fuel products
	fuelProducts := g.Group("/admin/fuel-products")
	fuelProducts.Get("/", controllers.GetAllFuelProductsHandler)
	fuelProducts.Get("/:id", controllers.GetFuelProductByIDHandler)
	fuelProducts.Post("/", controllers.CreateFuelProductHandler)
	fuelProducts.Patch("/:id", controllers.UpdateFuelProductHandler)
	fuelProducts.Delete("/:id", controllers.DeleteFuelProductHandler)

	// supplies
	supplies := g.Group("/admin/supplies")
	supplies.Get("/", controllers.GetSuppliesHandler)
	supplies.Get("/:id", controllers.GetSupplyByIDHandler)
	supplies.Post("/", controllers.AddSupplyHandler)
	supplies.Patch("/:id", controllers.UpdateSupplyHandler)
	supplies.Delete("/:id", controllers.DeleteSupplyHandler)
	
	//debts
	debts := g.Group("/admin/supplier/debts")
	debts.Get("/", controllers.GetSupplierDebtsHandler)
	

	// dippings
	dippings := g.Group("/admin/dippings")
	dippings.Get("/", controllers.GetAllDippingsHandler)
	dippings.Get("/:id", controllers.GetDippingByIDHandler)
	dippings.Get("/station/:id", controllers.GetDippingsByStationHandler)
	dippings.Get("/product/:id", controllers.GetDippingByFuelProductHandler)
	dippings.Get("/date/:date", controllers.GetDippingByDippingDateHandler)
	dippings.Post("/", controllers.CreateDippingHandler)
	dippings.Patch("/:id", controllers.UpdateDippingHandler)

	//
	d := g.Group("/admin/dippings-sales")
	d.Get("/", controllers.CompareDippingsAndSales)

	

	//sales
	sales := g.Group("/admin/sales")
	sales.Post("/", controllers.AddNewSalesHandler)
	sales.Patch("/:id", controllers.UpdateSalesHandler)
	sales.Delete("/:id", controllers.DeleteSalesHandler)
	sales.Get("/", controllers.GetSalesHandler)
	sales.Get("/date/:start_date/:end_date", controllers.GetSalesByDateHandler)

	//tanks
	tanks := g.Group("/admin/tanks")
	tanks.Get("/:station/:id", controllers.GetAllTanksHandler)
	tanks.Get("/:id", controllers.GetTankByIDHandler)
	tanks.Post("/", controllers.AddNewTankHandler)
	tanks.Patch("/:id", controllers.UpdateTankHandler)
	tanks.Delete("/:id", controllers.DeleteTankHandler)

	//pump
	pumps := g.Group("/admin/pumps")
	pumps.Get("/station/:id", controllers.GetPumpsByStationHandler)
	pumps.Get("/:id", controllers.GetPumpByIDHandler)
	pumps.Post("/", controllers.AddNewPumpHandler)
	pumps.Patch("/:id", controllers.UpdatePumpHandler)
	pumps.Delete("/:id", controllers.DeletePumpHandler)
	pumps.Post("/:tank_id/:pump_id", controllers.AssignPumpToTankHandler)
	pumps.Delete("/:tank_id/:pump_id", controllers.ReassignPumpToTankHandler)

	//create tank with pumps
	tanksWithPumps := g.Group("/admin/tanks-with-pumps")
	tanksWithPumps.Post("/", models.CreateTankWithPumps)

	//nozzles
	nozzles := g.Group("/admin/nozzles")
	nozzles.Get("/:id", controllers.GetNozzleByIDHandler)
	nozzles.Get("/", controllers.GetAllNozzlesHandler)
	nozzles.Post("/", controllers.CreateNozzleHandler)
	nozzles.Patch("/:id", controllers.UpdateNozzleHandler)
	nozzles.Delete("/:id", controllers.DeleteNozzleHandler)

	//payments
	payments := g.Group("/admin/payments")
	
	payments.Post("/", controllers.AddNewEmployeePayment)
	payments.Get("/report", controllers.GetPayrollReportHandler)


	//fuel stock
	stock := g.Group("/admin/stock")
	stock.Get("/", controllers.GetFuelStockHandler)

	//supplier payments
	payment := g.Group("/admin/supplier-payments")
	payment.Post("/", controllers.AddSupplierPayments)
	
}

