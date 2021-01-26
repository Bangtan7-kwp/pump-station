package app

import (
	"fmt"
	"log"
	"net/http"

	"api/app/handler"
	"api/app/model"
	"api/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/employees", a.GetAllEmployees)
	a.Post("/employees", a.CreateEmployee)
	a.Get("/employees/{title}", a.GetEmployee)
	a.Put("/employees/{title}", a.UpdateEmployee)
	a.Delete("/employees/{title}", a.DeleteEmployee)
	a.Put("/employees/{title}/disable", a.DisableEmployee)
	a.Put("/employees/{title}/enable", a.EnableEmployee)
	a.Post("/signupdata", a.CreateSignUPData)//check
	a.Post("/logindata",a.CreateLoginData)//check
	a.Post("/login",a.Login)//check
	a.Put("/deleteloginAccount/{username}", a.DeleteLoginData)//check
	a.Put("/ChangePwForAcc/{username}", a.ChangePwForAcc)//check
	a.Get("/petrolinfo", a.GetAllPetrolInfo)//check
	a.Post("/createPetrolInfo",a.CreatePetrolInfo)//check
	a.Put("/deletePetrolInfo/{petrol_name}", a.DeletePetrolData)//check
	a.Put("/changePetrolInfo/{petrol_name}", a.ChangePetrolData)//check
	a.Post("/createbook",a.CreateBook)//check
	a.Get("/bookinfo", a.GetAllBookInfo)//check
	a.Put("/changebook/{car_number}", a.ChangeBookInfo)//check
	a.Put("/deleteBookInfo/{car_number}", a.DeleteBookInfo)//check
	a.Post("/createtransaction",a.CreateTransaction)//check
	a.Get("/getTranperDay",a.GetAllTranPerDay)//check
	a.Get("/getallTran",a.GetAllTran)//check
	a.Get("/getallTranforDebt",a.GetAllTranforDebt)
	a.Put("/getallTranonDay/{income_date}",a.GetAllTranOnDay)//check
	a.Put("/getallTranbyCar/{car_number}",a.GetAllTranByCarNo)//check
	a.Put("/updateTransaction",a.UpdateTran)//check
	a.Put("/deleteTran/{id}", a.DeleteTran)//check
	a.Post("/amtforNormal",a.AmtforNormalCus)//check
	a.Post("/amtforVIP",a.AmtforVIPCus)//check
	a.Post("/creategift",a.CreateGiftLiterAndPrice)//check
	a.Get("/getgift", a.GetGiftLiterAndPrice)//check
	a.Get("/getallDebt", a.GetAllDebt)
	a.Get("/getallDebtforToday", a.GetAllDebtForToday)
	a.Put("/getallDebtonDay/{debt_date}",a.GetAllDebtOnDay)
	a.Put("/getallDebtbycar/{car_number}",a.GetAllDebtByCarNo)
	a.Put("/updateDebtTran",a.UpdateDebt)
	a.Post("/paidOnDebt",a.PaidOnDebt)
	a.Get("/getallBookTran", a.GetAllBookTran)
	a.Put("/getBookTranByCarNo/{car_number}",a.GetAllBookTranByCarNo)
	a.Get("/getallBookTranByGiftLiter", a.GetAllBookTranByGiftLiter)
	a.Post("/creategift_tran", a.CreateGiftTran)
	a.Get("/getlastBookTranByGiftLiter", a.GetLastBookTranByGiftLiter)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Handlers to manage Employee Data
func (a *App) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	handler.GetAllEmployees(a.DB, w, r)
}

func (a *App) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	handler.CreateEmployee(a.DB, w, r)
}

func (a *App) GetEmployee(w http.ResponseWriter, r *http.Request) {
	handler.GetEmployee(a.DB, w, r)
}

func (a *App) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	handler.UpdateEmployee(a.DB, w, r)
}

func (a *App) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	handler.DeleteEmployee(a.DB, w, r)
}

func (a *App) DisableEmployee(w http.ResponseWriter, r *http.Request) {
	handler.DisableEmployee(a.DB, w, r)
}

func (a *App) EnableEmployee(w http.ResponseWriter, r *http.Request) {
	handler.EnableEmployee(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *App) CreateSignUPData(w http.ResponseWriter, r *http.Request) {
	handler.CreateSignUPData(a.DB, w, r)
}

func (a *App) CreateLoginData(w http.ResponseWriter, r *http.Request) {
	handler.CheckLoginData(a.DB, w, r,"","")
}

func (a *App) DeleteLoginData(w http.ResponseWriter, r *http.Request) {
	handler.DeleteLoginAccount(a.DB, w, r)
}

func (a *App) ChangePwForAcc(w http.ResponseWriter, r *http.Request) {
	handler.ChangePasswordForLoginAcc(a.DB, w, r)
}

func (a *App) GetAllPetrolInfo(w http.ResponseWriter, r *http.Request) {
	handler.GetAllPetrolInfo(a.DB, w, r)
}

func (a *App) CreatePetrolInfo(w http.ResponseWriter, r *http.Request) {
	handler.CreatePetrolInfo(a.DB, w, r)
}

func (a *App) DeletePetrolData(w http.ResponseWriter, r *http.Request) {
	handler.DeletePetrolInfo(a.DB, w, r)
}

func (a *App) ChangePetrolData(w http.ResponseWriter, r *http.Request) {
	handler.ChangePetrolInfo(a.DB, w, r)
}

func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	handler.CreateBook(a.DB, w, r)
}

func (a *App) ChangeBookInfo(w http.ResponseWriter, r *http.Request) {
	handler.ChangeBookInfo(a.DB, w, r)
}

func (a *App) GetAllBookInfo(w http.ResponseWriter, r *http.Request) {
	handler.GetAllBook(a.DB, w, r)
}

func (a *App) DeleteBookInfo(w http.ResponseWriter, r *http.Request) {
	handler.DeleteBookInfo(a.DB, w, r)
}

func (a *App) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	handler.CreateTransaction(a.DB, w, r)
}

func (a *App) GetAllTranPerDay(w http.ResponseWriter, r *http.Request) {
	handler.GetAllIncomeTransactionPerDay(a.DB, w, r)
}

func (a *App) GetAllTran(w http.ResponseWriter, r *http.Request) {
	handler.GetAllIncomeTransaction(a.DB, w, r)
}

func (a *App) GetAllTranforDebt(w http.ResponseWriter, r *http.Request) {
	handler.GetAllIncomeTransactionForDebt(a.DB, w, r)
}

func (a *App) GetAllTranOnDay(w http.ResponseWriter, r *http.Request) {
	handler.GetAllIncomeTransactionOnDay(a.DB, w, r)
}

func (a *App) GetAllTranByCarNo(w http.ResponseWriter, r *http.Request) {
	handler.GetAllIncomeTransactionbyCarNo(a.DB, w, r)
}

func (a *App) UpdateTran(w http.ResponseWriter, r *http.Request) {
	handler.UpdateTransaction(a.DB, w, r)
}

func (a *App) DeleteTran(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTransaction(a.DB, w, r)
}

func (a *App) AmtforNormalCus(w http.ResponseWriter, r *http.Request) {
	handler.AmtForNormalCustomer(a.DB, w, r)
}

func (a *App) AmtforVIPCus(w http.ResponseWriter, r *http.Request) {
	handler.AmtForVIPCustomer(a.DB, w, r)
}

func (a *App) CreateGiftLiterAndPrice(w http.ResponseWriter, r *http.Request) {
	handler.CreateGiftLiterAndPrice(a.DB, w, r)
}

func (a *App) GetGiftLiterAndPrice(w http.ResponseWriter, r *http.Request) {
	handler.GetGiftLiterAndPrice(a.DB, w, r)
}

func (a *App) GetAllDebt(w http.ResponseWriter, r *http.Request) {
	handler.GetAllDebt(a.DB, w, r)
}

func (a *App) GetAllDebtForToday(w http.ResponseWriter, r *http.Request) {
	handler.GetAllDebtForToday(a.DB, w, r)
}

func (a *App) GetAllDebtOnDay(w http.ResponseWriter, r *http.Request) {
	handler.GetAllDebtOnDay(a.DB, w, r)
}

func (a *App) GetAllDebtByCarNo(w http.ResponseWriter, r *http.Request) {
	handler.GetAllDebtbyCarNo(a.DB, w, r)
}

func (a *App) UpdateDebt(w http.ResponseWriter, r *http.Request) {
	handler.UpdateDebt(a.DB, w, r)
}

func (a *App) PaidOnDebt(w http.ResponseWriter, r *http.Request) {
	handler.PaidOnDebt(a.DB, w, r)
}

func (a *App) GetAllBookTran(w http.ResponseWriter, r *http.Request) {
	handler.GetBookTransaction(a.DB, w, r)
}

func (a *App) GetAllBookTranByCarNo(w http.ResponseWriter, r *http.Request) {
	handler.GetBookTransactionByCarNo(a.DB, w, r)
}

func (a *App) GetAllBookTranByGiftLiter(w http.ResponseWriter, r *http.Request) {
	handler.GetBookTransactionByGiftLiter(a.DB, w, r)
}

func (a *App) CreateGiftTran(w http.ResponseWriter, r *http.Request) {
	handler.CreateGiftTrans(a.DB, w, r)
}

func (a *App) GetLastBookTranByGiftLiter(w http.ResponseWriter, r *http.Request) {
	handler.GetLastBookTransactionByGiftLiter(a.DB, w, r)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	handler.GetLastBookTransactionByGiftLiter(a.DB, w, r)
}