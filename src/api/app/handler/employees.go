package handler

import (
	"api/app/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetAllEmployees(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	employees := []model.Employee{}
	db.Find(&employees)
	respondJSON(w, http.StatusOK, employees)
}

func GetAllPetrolInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	pinfo := []model.Petrol_Info{}
	db.Where("access = ?", true).Find(&pinfo)
	respondJSON(w, http.StatusOK, pinfo)
}

func CreateSignUPData(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	signup := model.SignUp{}
	login :=model.LogIn{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signup); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	sforName:=model.SignUp{}
	sforPassword:=model.SignUp{}
	err1:=db.First(&sforName,model.SignUp{SName:signup.SName}).Error
	err2:=db.First(&sforPassword, model.SignUp{SPassword: signup.SPassword}).Error
	fmt.Println(err1,"Error of Name")
	fmt.Println(sforName,"Name")
	fmt.Println(sforName,"Password")
	fmt.Println(err2,"Error of Password")
	if err:=db.First(&signup,model.SignUp{SName:signup.SName,SPassword:signup.SPassword}).Error; err==nil{
		respondJSON(w,http.StatusCreated,"Username and password exit")
	}
	if sforName.SName=="" && sforPassword.SPassword==""{
		fmt.Println("Hello 1st if statement")
				if signup.SPassword == signup.ReEnter_Password {
					id, err := uuid.NewV4()
					if err != nil {
						fmt.Println("There is error in generating uuid")
					}
					signup.SignUp_ID = id
					signup.CreatedAt = time.Now().Local()
					fmt.Println(signup.CreatedAt, "Created time")
					/*a,err:=time.Parse(createdtime.String(),"07:05:45PM")
					fmt.Println(err,"error of date format")
					signup.CreatedAt=a*/
					signup.UpdatedAt = time.Now().Local()
					fmt.Println(signup.UpdatedAt, "Updated time")
					if err := db.Save(&signup).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
					fmt.Println(signup.ID, "Id of SignUP")
					fmt.Println(signup.CreatedAt,"Created time of signup")

					logid, err := uuid.NewV4()
					if err != nil {
						fmt.Println("There is error in generating uuid")
					}
					login.Login_ID = logid
					login.SignUp_ID_Fk = signup.SignUp_ID
					login.UserName = signup.SName
					login.Password = signup.SPassword
					login.Access = true
					if err := db.Save(&login).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
					//SendEmail(signup.Email,"This is your signup Password for Pump Station"+signup.SPassword)
					respondJSON(w, http.StatusOK, "SignUp is successful")
				}else{
					respondJSON(w,http.StatusCreated,"Passsword and ReEnter Password should be same")
				}
		}else if err1==nil && err2!=nil{
			respondJSON(w,http.StatusCreated,"Please write UserName again")
	}else if err1!=nil && err2==nil{
		respondJSON(w,http.StatusCreated,"Please write Password again")
	}
}

func CheckLoginData(db *gorm.DB, w http.ResponseWriter, r *http.Request,username string,password string) {
	log:=model.LogIn{}
	/*decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&log); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Println(log.UserName,"UserName")
	fmt.Println(log.Password,"Password")*/
	sign:=model.SignUp{}
	if err:= db.First(&sign, model.SignUp{SName: username,SPassword: password}).Error; err!=nil {
			respondJSON(w,http.StatusCreated,"Username and Password don't exist!!!")
	} else {
		fmt.Println(sign.SignUp_ID, "Sign Up UUID")
		log1:=model.LogIn{}
		if err := db.First(&log1, model.LogIn{UserName: username,Password: password}).Error; err != nil {
			respondJSON(w, http.StatusCreated, "Username and Password are wrong!!!")
		} else {
			id_log, err := uuid.NewV4()
			if err != nil {
				fmt.Println("There is error in generating uuid")
			}
			log.Login_ID = id_log
			log.SignUp_ID_Fk = sign.SignUp_ID
			log.Access = true
			if err := db.Save(&log).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondJSON(w, http.StatusOK, "Login is Successful")
		}
	}
}

func CreateEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	employee := model.Employee{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, employee)
}

func CreatePetrolInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	petrolinfo := model.Petrol_Info{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&petrolinfo); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("There is error in generating uuid")
	}
	if err := db.First(&petrolinfo,model.Petrol_Info{Petrol_Name:petrolinfo.Petrol_Name}).Error; err == nil {
		respondError(w, http.StatusCreated, "This Petrol is Existing!")
	}else {

		petrolinfo.Petrol_ID = id
		petrolinfo.Access = true
		if err := db.Save(&petrolinfo).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("Create new Petrol Info")
		respondJSON(w, http.StatusCreated, petrolinfo)
	}
}

func GetEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func UpdateEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func DeleteEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	if err := db.Delete(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func DisableEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	employee.Disable()
	if err := db.Save(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func EnableEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	employee.Enable()
	if err := db.Save(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getEmployeeOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Employee {
	employee := model.Employee{}
	if err := db.First(&employee, model.Employee{Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &employee
}

func getLoginOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.LogIn {
	login := model.LogIn{}
	if err := db.First(&login, model.LogIn{UserName: name}).Error; err != nil {
		//respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &login
}

func getSignUpOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.SignUp {
	singup := model.SignUp{}
	if err := db.First(&singup, model.SignUp{SName: name}).Error; err != nil {
		//respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &singup
}

func getPetrolInfoOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Petrol_Info {
	pinfo := model.Petrol_Info{}
	fmt.Println(name,"Name of input")
	if err := db.Find(&pinfo, model.Petrol_Info{Petrol_Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	fmt.Println(&pinfo.Petrol_Name,"Name from database")
	return &pinfo
}

func DeleteLoginAccount (db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["username"]
	fmt.Println(name,"Username from UI")
	loginacc := getLoginOr404(db, name, w, r)
	signupacc := getSignUpOr404(db, name, w, r)
	if loginacc == nil {
		respondJSON(w, http.StatusCreated, "Deleting Login Account is Unsuccessful")
		return
	}
	loginacc.Access=false
	if err := db.Save(&loginacc).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if signupacc == nil{
		respondJSON(w, http.StatusCreated, "Deleting Login Account is Unsuccessful")
		return
	}
	db.Delete(signupacc)
	fmt.Println("Set Deletedat time in mysql for signup")
	db.Delete(loginacc)
	fmt.Println("Set Deletedat time in mysql for login")
	respondJSON(w, http.StatusCreated, "Delete Login Account Successful")
}

func DeletePetrolInfo (db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["petrol_name"]
	fmt.Print(name,"Name of Delete Petrol Info")
	pinfo1 := getPetrolInfoOr404(db, name, w, r)
	if pinfo1==nil {
		return
		respondJSON(w, http.StatusCreated, "Delete Petrol Info Unsuccessful")
	}
		pinfo1.Access = false
		if err := db.Save(&pinfo1).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		db.Delete(pinfo1)
		fmt.Println("Set Deletedat time in mysql for petrol info")
		respondJSON(w, http.StatusCreated, "Delete Petrol Info Successful")

}

func ChangePasswordForLoginAcc(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["username"]
	fmt.Println(name,"Name from UI")
	logforUpdate := getLoginOr404(db, name, w, r)
	signUpdate := getSignUpOr404(db, name, w, r)
	if logforUpdate == nil {
		return
	}
	if signUpdate==nil{
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&logforUpdate); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	signUpdate.SPassword=logforUpdate.Password
	signUpdate.ReEnter_Password=logforUpdate.Password
	if err := db.Save(&signUpdate).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Save(&logforUpdate).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated,"Changing your password is successful")
}

func ChangePetrolInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["petrol_name"]
	petrol := getPetrolInfoOr404(db, name, w, r)
	if petrol == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&petrol); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	petrol.UpdatedAt=time.Now().Local()
	fmt.Println(petrol.UpdatedAt,"Updated time from golang")
	if err := db.Save(&petrol).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated,"Changing petrol information is successful")
}

func CreateBook(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	book := model.Customer_Info{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("There is error in generating uuid")
	}
	if err := db.Find(&book,model.Customer_Info{Car_Number:book.Car_Number}).Error; err == nil {
		respondError(w, http.StatusCreated, "This Book is Existing!")
	}else {

		book.Customer_ID = id
		book.Access = true
		if err := db.Save(&book).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("Create new Book Info")
		respondJSON(w, http.StatusCreated, book)
	}
}

func GetAllBook(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	b := []model.Customer_Info{}
	db.Where("access = ?", true).Find(&b)
	respondJSON(w, http.StatusOK, b)
}

func getBookOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Customer_Info {
	b1 := model.Customer_Info{}
	fmt.Println(name,"Name of input")
	db.Find(&b1, model.Customer_Info{Car_Number: name,Access:true})

	/*if err := db.Find(&b1, model.Customer_Info{Car_Number: name}).Error; err != nil {
		//respondError(w, http.StatusNotFound, err.Error())
		return nil
	}*/
	fmt.Println(&b1.Customer_ID,"Car Number from database")
	return &b1
}

func getDebtInfo(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) []model.Debt {
	b1 := []model.Debt{}
	fmt.Println(name,"Name of input")
	db.Find(&b1, model.Debt{Car_Number: name,Status:false})

	/*if err := db.Find(&b1, model.Customer_Info{Car_Number: name}).Error; err != nil {
		//respondError(w, http.StatusNotFound, err.Error())
		return nil
	}*/
	fmt.Println("Hello GetDebtInfo Function")
	return b1
}

func ChangeBookInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["car_number"]
	url.QueryEscape(name)
	b := getBookOr404(db, name, w, r)
	if b == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&b); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	/*b.UpdatedAt=time.Now().Local()
	fmt.Println(b.UpdatedAt,"Updated time from golang")*/
	if err := db.Save(&b).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated,"Changing Book Information is successful")
}

func DeleteBookInfo (db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["car_number"]
	url.QueryEscape(name)
	book := getBookOr404(db, name, w, r)
	if book == nil {
		return
	}
	book.Access=false
	if err := db.Save(&book).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	db.Delete(book)
	fmt.Println("Set Deletedat time in mysql for book info")
	respondJSON(w, http.StatusCreated, "Delete Book Information Successful")
}

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	t := model.Transaction{}
	income:=model.OneDay_Transaction{}
	fmt.Println(income.Car_Number,"Car Number before calling api")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("There is error in generating uuid")
	}
	p := model.Petrol_Info{}
	if err := db.Find(&p, model.Petrol_Info{Petrol_Name: t.Petrol_Name,Access: true}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	fmt.Println(t.Car_Number,"Car Number from Json")
	customer:=getBookOr404(db,t.Car_Number,w,r)
	fmt.Println(customer, "Customer Information")
	if customer.Customer_ID==uuid.Nil{
		fmt.Println("Hello,customer book is nil")
		income.Petrol_ID_Fk = p.Petrol_ID
		income.Income_ID = id
		income.Income_Date = t.Income_Date
		income.Liter_Per_Once = t.Liter
		income.Car_Number = t.Car_Number
		income.Info = t.Info
		income.Status = t.Paid
		income.Permission = t.Permission
		income.Amount = t.Amount
		income.Low_Price = t.LowPrice
		income.Customer_ID_Fk =uuid.Nil

		if income.Status == true { //paid for amount
			income.Paid_Amt = income.Amount
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		d := model.Debt{}
		debt_id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("There is error in generating uuid")
		}
		d.Debt_ID = debt_id
		d.Tran_ID=id
		d.Customer_ID_Fk = uuid.Nil
		d.Debt_Date = income.Income_Date
		d.Car_Number=income.Car_Number

		if t.Debt == true {
			income.Status = false
			income.Paid_Amt = 0.0
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			d.Debt_Amount = income.Amount
			d.Residual_Amount = income.Amount
			d.Status = false
			if err := db.Save(&d).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			debtarr:=getDebtInfo(db,t.Car_Number,w,r)
			fmt.Println("Length of Debt Array", len(debtarr))
			if len(debtarr)==1{
				d.Total_Debt = d.Residual_Amount
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}else if len(debtarr) > 1 {
				fmt.Println("Hello when length of debtarr is greater than 1")
				debt := model.Debt{}
				if err := db.Where(" car_number= ? and status= ?",t.Car_Number,false).Last(&debt).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
				}
				/*debttotal := model.Debt{}
				debtid := debt.ID - 1
				db.Where(" id= ? and status= ?", debtid,false).Find(&debttotal)
				total := debttotal.Total_Debt + debt.Residual_Amount
				d.Total_Debt = total*/
				/*if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}*/
				debttotal:=[]model.Debt{}
				db.Where("status=? and car_number= ?",false,t.Car_Number).Find(&debttotal)
				for i:=0;i<len(debttotal);i++{
					d.Total_Debt+=debttotal[i].Residual_Amount
				}
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		if t.Both == true {
			income.Paid_Amt = t.Paid_Amt
			income.Status = false
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			//d.Debt_Amount=t.Debt_Amt
			d.Debt_Amount = t.Amount - t.Paid_Amt
			d.Residual_Amount = t.Amount - t.Paid_Amt
			d.Paid_Amount = 0
			d.Paid_Date = t.Income_Date
			d.Status = false
			if err := db.Save(&d).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			debtarr:=getDebtInfo(db,t.Car_Number,w,r)
			fmt.Println("Length of Debt Array", len(debtarr))
			if len(debtarr)==1{
				d.Total_Debt = d.Residual_Amount
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}else if len(debtarr) > 1 {
				fmt.Println("Hello when length of debtarr is greater than 1")
				debt := model.Debt{}
				if err := db.Where(" car_number= ? and status= ?",t.Car_Number,false).Last(&debt).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
				}
				/*debttotal := model.Debt{}
				debtid := debt.ID - 1
				db.Where(" id= ? and status= ?", debtid,false).Find(&debttotal)
				total := debttotal.Total_Debt + debt.Residual_Amount
				d.Total_Debt = total*/
				/*if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}*/
				debttotal:=[]model.Debt{}
				db.Where("status=? and car_number= ?",false,t.Car_Number).Find(&debttotal)
				for i:=0;i<len(debttotal);i++{
					d.Total_Debt+=debttotal[i].Residual_Amount
				}
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}

	}else if customer.Customer_ID!=uuid.Nil{
		fmt.Println(customer.ID, "Customer ID")
		income.Petrol_ID_Fk = p.Petrol_ID
		income.Income_ID = id
		income.Income_Date = t.Income_Date
		income.Liter_Per_Once = t.Liter
		income.Car_Number = t.Car_Number
		income.Info = t.Info
		income.Status = t.Paid
		income.Permission = t.Permission
		income.Amount = t.Amount
		income.Low_Price = t.LowPrice
		income.Customer_ID_Fk = customer.Customer_ID

		if income.Status == true { //paid for amount
			income.Paid_Amt = income.Amount
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		d := model.Debt{}
		debt_id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("There is error in generating uuid")
		}
		d.Debt_ID = debt_id
		d.Tran_ID=id
		d.Customer_ID_Fk = customer.Customer_ID
		d.Debt_Date = income.Income_Date
		d.Car_Number=income.Car_Number
		//Debt is true
		if t.Debt == true {
			income.Status=false
			income.Paid_Amt = 0.0
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			d.Debt_Amount = income.Amount
			d.Residual_Amount = income.Amount
			d.Status = false
			if err := db.Save(&d).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			debtarr:=getDebtInfo(db,t.Car_Number,w,r)
			fmt.Println("Length of Debt Array", len(debtarr))
			if len(debtarr)==1{
				d.Total_Debt = d.Residual_Amount
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}else if len(debtarr) > 1 {
				fmt.Println("Hello when length of debtarr is greater than 1")
				debt := model.Debt{}
				if err := db.Where(" car_number= ? and status= ?",t.Car_Number,false).Last(&debt).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
				}
				/*debttotal := model.Debt{}
				debtid := debt.ID - 1
				db.Where(" id= ? and status= ?", debtid,false).Find(&debttotal)
				total := debttotal.Total_Debt + debt.Residual_Amount
				d.Total_Debt = total*/
				/*if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}*/
				debttotal:=[]model.Debt{}
				db.Where("status=? and car_number= ?",false,t.Car_Number).Find(&debttotal)
				for i:=0;i<len(debttotal);i++{
					d.Total_Debt+=debttotal[i].Residual_Amount
				}
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		//Paid for some amt and Debt for rest amt
		if t.Both == true {
			income.Paid_Amt = t.Paid_Amt
			income.Status = false
			if err := db.Save(&income).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			//d.Debt_Amount=t.Debt_Amt
			d.Debt_Amount = t.Amount - t.Paid_Amt
			d.Residual_Amount = t.Amount - t.Paid_Amt
			d.Paid_Amount = 0
			d.Paid_Date = t.Income_Date
			d.Status = false
			if err := db.Save(&d).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			debtarr:=getDebtInfo(db,t.Car_Number,w,r)
			fmt.Println("Length of Debt Array", len(debtarr))
			if len(debtarr)==1{
				d.Total_Debt = d.Residual_Amount
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}else if len(debtarr) > 1 {
				fmt.Println("Hello when length of debtarr is greater than 1")
				debt := model.Debt{}
				if err := db.Where(" car_number= ? and status= ?",t.Car_Number,false).Last(&debt).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
				}
				/*debttotal := model.Debt{}
				debtid := debt.ID - 1
				db.Where(" id= ? and status= ?", debtid,false).Find(&debttotal)
				total := debttotal.Total_Debt + debt.Residual_Amount
				d.Total_Debt = total*/
				/*if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}*/
				debttotal:=[]model.Debt{}
				db.Where("status=? and car_number= ?",false,t.Car_Number).Find(&debttotal)
				for i:=0;i<len(debttotal);i++{
					d.Total_Debt+=debttotal[i].Residual_Amount
				}
				if err := db.Save(&d).Error; err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}


		bookTran := model.Book_Transaction{}
		book_id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("There is error in generating uuid")
		}
		bookTran.Transaction_ID = book_id
		bookTran.Income_ID_Fk = income.Income_ID
		bookTran.Petrol_ID_Fk = p.Petrol_ID
		bookTran.TLiter = income.Liter_Per_Once
		bookTran.TDate = income.Income_Date
		bookTran.Customer_ID_Fk = customer.Customer_ID
		bookTran.Status=true
		if err := db.Save(&bookTran).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		bookarr := []model.Book_Transaction{}
		if err := db.Where(" customer_id_fk= ?", customer.Customer_ID).Find(&bookarr).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		if len(bookarr) == 1 {
			bookTran.Total_Liter = bookTran.TLiter
			if err := db.Save(&bookTran).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if len(bookarr) > 1 {
			book := model.Book_Transaction{}
			if err := db.Where(" customer_id_fk= ?", customer.Customer_ID).Last(&book).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
			}
			bookTotal := model.Book_Transaction{}
			bookid := book.ID - 1
			if err := db.Where(" id= ? and customer_id_fk= ?", bookid,customer.Customer_ID).Find(&bookTotal).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
			}
			total := bookTotal.Total_Liter + book.TLiter
			bookTran.Total_Liter = total
			/*bookTran.Status=false*/
			if err := db.Save(&bookTran).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			/*bookTotal:=[]model.Book_Transaction{}
			db.Where("customer_id_fk= ?",customer.Customer_ID).Find(&bookTotal)
			for i:=0;i<len(bookTotal);i++{
				bookTran.Total_Liter+=bookTotal[i].TLiter
			}
			if err := db.Save(&bookTran).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}*/
		}

		if t.Gift==true {
			gift := getGiftOr404(db, w, r)
			bookfortoal := model.Book_Transaction{}
			if err := db.Where("customer_id_fk= ?", bookTran.Customer_ID_Fk).Last(&bookfortoal).Error; err != nil {
				respondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if bookfortoal.Total_Liter >= gift.Gift_Liter {
				total := bookfortoal.Total_Liter - gift.Gift_Liter
				db.Model(&bookfortoal).Where("customer_id_fk = ?", bookTran.Customer_ID_Fk).Update(model.Book_Transaction{Total_Liter: total})

				if t.Amount >= gift.Gift_Amount {
					/*income.Status = true
					if err := db.Save(&income).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}*/
					giftran := model.GiftTransaction{}
					gf_id, err := uuid.NewV4()
					if err != nil {
						fmt.Println("There is error in generating uuid")
					}
					giftran.Gift_Tran_ID = gf_id
					giftran.Gift_Tran_Car_Number = t.Car_Number
					giftran.Gift_Tran_Date = t.Income_Date
					giftran.Gift_Tran_Liter = gift.Gift_Liter
					giftran.Gift_Tran_Total_Liter = total
					if err := db.Save(&giftran).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
					exp := model.Expense{}
					exp_id, err := uuid.NewV4()
					if err != nil {
						fmt.Print("There is error in generating uuid")
					}
					exp.Exp_ID = exp_id
					exp.Date = t.Income_Date
					exp.Amount = gift.Gift_Amount
					exp.Info = "Gift Liter Expense for " + t.Car_Number
					if err := db.Save(&exp).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
					bookarr := []model.Book_Transaction{}
					db.Where("customer_id_fk= ?", customer.Customer_ID).Find(&bookarr)
					var sum float32
					for i := 0; i < len(bookarr); i++ {
						sum += bookarr[i].TLiter
						if sum > gift.Gift_Amount {
							for j := 0; j < i; j++ {
								db.Where("transaction_id= ?", bookarr[j].Transaction_ID).Delete(&bookarr[j])
							}
							bookarr[i].Total_Liter = sum - gift.Gift_Amount
							/*bookarr[i].Status=true*/
							if err := db.Save(&bookarr[i]).Error; err != nil {
								respondError(w, http.StatusInternalServerError, err.Error())
								return
							}
							booktotal := []model.Book_Transaction{}
							db.Where("customer_id_fk= ? and status= ?", customer.Customer_ID,true).Find(&booktotal)
							for k := 1; k < len(booktotal); k++ {
								booktotal[0].Total_Liter += booktotal[k].TLiter
								//booktotal[k].Status=false
								if err := db.Save(&booktotal[k]).Error; err != nil {
									respondError(w, http.StatusInternalServerError, err.Error())
									return
								}
							}
							break
						} else if sum == gift.Gift_Amount {
							for j := 0; j <= i; j++ {
								db.Where("transaction_id= ?", bookarr[j].Transaction_ID).Delete(&bookarr[j])
							}
							bookarr[i].Total_Liter = 0
							if err := db.Save(&bookarr[i]).Error; err != nil {
								respondError(w, http.StatusInternalServerError, err.Error())
								return
							}
							booktotal := []model.Book_Transaction{}
							db.Where("customer_id_fk= ? and status=?", customer.Customer_ID,true).Find(&booktotal)
							for k := 1; k < len(booktotal); k++ {
								booktotal[0].Total_Liter += booktotal[k].TLiter
								//booktotal[k].Status=false
								if err := db.Save(&booktotal[k]).Error; err != nil {
									respondError(w, http.StatusInternalServerError, err.Error())
									return
								}
							}
							break
						} else if sum < gift.Gift_Amount {
							continue
						}

					}
				}


				if t.Debt == true {
					income.Status=false
					if err := db.Save(&income).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
					d.Paid_Amount += gift.Gift_Amount
					d.Residual_Amount = d.Debt_Amount - d.Paid_Amount
					d.Total_Debt = d.Total_Debt - d.Paid_Amount
					if err := db.Save(&d).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
				}
				if t.Both==true{
					d.Paid_Amount += gift.Gift_Amount
					d.Residual_Amount = d.Debt_Amount - d.Paid_Amount
					d.Total_Debt = d.Total_Debt - d.Paid_Amount
					if err := db.Save(&d).Error; err != nil {
						respondError(w, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}

		}
	}

		fmt.Println("Create a new Transaction")
		respondJSON(w, http.StatusCreated, "Successfully create a new transaction")
}

func GetAllIncomeTransactionPerDay(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	income:=[]model.OneDay_Transaction{}

	/*if err:=db.Where("income_date=?","2014-01-01 23:28:57").Find(&income).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
	}*/

	//db.Raw("select * from  where DATE(date) = '2014-03-19';")
	date:=time.Now()
	res:=strings.SplitAfterN(date.String()," ",2)
	fmt.Println(date,"Date Format")
	fmt.Println(res[0],"Date split")

	startDate:=res[0]+" 00:00:00"
	endDate:=res[0]+" 23:59:59"
	layout := "2006-01-02 15:04:05"
	tstart, err:= time.Parse(layout, startDate)
	tend,err:=time.Parse(layout,endDate)
	fmt.Println(tstart,"Start Date and Time")
	fmt.Println(tend,"End Date and Time")
	if err!=nil{
		respondJSON(w,http.StatusNotFound,err.Error())
	}
	db.Where("income_date between ? and ?",startDate,endDate).Find(&income)

	all := []model.GetAllIncomeTransaction{}
	for i:=0;i<len(income);i++{
		alltran:=model.GetAllIncomeTransaction{}
		alltran.Income_No= int(income[i].ID)
		alltran.Income_Date=income[i].Income_Date
		alltran.Car_Number=income[i].Car_Number
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income[i].Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
		}
		alltran.Petrol_Name=petrol.Petrol_Name
		if income[i].Permission==false{
			alltran.Petrol_Price=petrol.Price
		}else{
			alltran.Petrol_Price=income[i].Low_Price
		}
		alltran.Liter=income[i].Liter_Per_Once
		alltran.Amount=income[i].Amount
		alltran.Paid_Amt=income[i].Paid_Amt
		alltran.Permission=income[i].Permission
		alltran.Debt_Amt=alltran.Amount-alltran.Paid_Amt
		alltran.Status=income[i].Status
		all =append(all, alltran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}
}

func GetAllIncomeTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	income:=[]model.OneDay_Transaction{}

	/*if err:=db.Where("income_date=?","2014-01-01 23:28:57").Find(&income).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
	}*/

	//db.Raw("select * from  where DATE(date) = '2014-03-19';")
	/*date:=time.Now()
	res:=strings.SplitAfterN(date.String()," ",2)
	fmt.Println(date,"Date Format")
	fmt.Println(res[0],"Date split")*/

	db.Find(&income)

	all := []model.GetAllIncomeTransaction{}
	for i:=0;i<len(income);i++{
		alltran:=model.GetAllIncomeTransaction{}
		alltran.Income_No= int(income[i].ID)
		alltran.Income_Date=income[i].Income_Date
		alltran.Car_Number=income[i].Car_Number
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income[i].Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		alltran.Petrol_Name=petrol.Petrol_Name
		if income[i].Permission==false{
			alltran.Petrol_Price=petrol.Price
		}else{
			alltran.Petrol_Price=income[i].Low_Price
		}
		alltran.Liter=income[i].Liter_Per_Once
		alltran.Amount=income[i].Amount
		alltran.Paid_Amt=income[i].Paid_Amt
		alltran.Permission=income[i].Permission
		alltran.Debt_Amt=alltran.Amount-alltran.Paid_Amt
		alltran.Status=income[i].Status
		all =append(all, alltran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}
}

func GetAllIncomeTransactionOnDay(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Print("Hello Transaction on Day")
	vars := mux.Vars(r)

	date := vars["income_date"]
	url.QueryEscape(date)
	fmt.Println(date, "Date from UI")
	startDate:=date+" 00:00:00"
	endDate:=date+" 23:59:59"
	layout := "2006-01-02 15:04:05"
	tstart, err:= time.Parse(layout, startDate)
	tend,err:=time.Parse(layout,endDate)
	fmt.Println(tstart,"Start Date and Time")
	fmt.Println(tend,"End Date and Time")
	if err != nil {
		fmt.Println(err)
	}

	income:=[]model.OneDay_Transaction{}


	if err:=db.Where("income_date between ? and ?",startDate,endDate).Find(&income).Error;err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	fmt.Println("Hello alltran")
	all := []model.GetAllIncomeTransaction{}
	for i:=0;i<len(income);i++{
		alltran:=model.GetAllIncomeTransaction{}
		alltran.Income_No= int(income[i].ID)
		alltran.Income_Date=income[i].Income_Date
		alltran.Car_Number=income[i].Car_Number
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income[i].Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}

		alltran.Petrol_Name=petrol.Petrol_Name
		if income[i].Permission==false{
			alltran.Petrol_Price=petrol.Price
		}else{
			alltran.Petrol_Price=income[i].Low_Price
		}
		alltran.Liter=income[i].Liter_Per_Once
		alltran.Amount=income[i].Amount
		alltran.Paid_Amt=income[i].Paid_Amt
		alltran.Permission=income[i].Permission
		alltran.Debt_Amt=alltran.Amount-alltran.Paid_Amt
		alltran.Status=income[i].Status
		all =append(all, alltran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}
}

func GetAllIncomeTransactionForDebt(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	income:=[]model.OneDay_Transaction{}
	fmt.Println("Hello For query of car number")
	if err:=db.Where("status=?",false).Find(&income).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	fmt.Println("Array of length")
	all := []model.GetAllIncomeTransaction{}
	for i:=0;i<len(income);i++{
		alltran:=model.GetAllIncomeTransaction{}
		alltran.Income_No= int(income[i].ID)
		alltran.Income_Date=income[i].Income_Date
		alltran.Car_Number=income[i].Car_Number
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income[i].Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
			break
		}
		alltran.Petrol_Name=petrol.Petrol_Name
		if income[i].Permission==false{
			alltran.Petrol_Price=petrol.Price
		}else{
			alltran.Petrol_Price=income[i].Low_Price
		}
		alltran.Liter=income[i].Liter_Per_Once
		alltran.Amount=income[i].Amount
		alltran.Paid_Amt=income[i].Paid_Amt
		alltran.Permission=income[i].Permission
		alltran.Debt_Amt=alltran.Amount-alltran.Paid_Amt
		alltran.Status=income[i].Status
		all =append(all, alltran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}
}

func GetAllIncomeTransactionbyCarNo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Transaction By Car Number")

	vars := mux.Vars(r)

	car := vars["car_number"]
	url.QueryEscape(car)
	fmt.Println(car, "car number")
	income:=[]model.OneDay_Transaction{}
	fmt.Println("Hello For query of car number")
	if err:=db.Where("car_number=?",car).Find(&income).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	fmt.Println("Array of length")
	all := []model.GetAllIncomeTransaction{}
	for i:=0;i<len(income);i++{
		alltran:=model.GetAllIncomeTransaction{}
		alltran.Income_No= int(income[i].ID)
		alltran.Income_Date=income[i].Income_Date
		alltran.Car_Number=income[i].Car_Number
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income[i].Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
			break
		}
		alltran.Petrol_Name=petrol.Petrol_Name
		if income[i].Permission==false{
			alltran.Petrol_Price=petrol.Price
		}else{
			alltran.Petrol_Price=income[i].Low_Price
		}
		alltran.Liter=income[i].Liter_Per_Once
		alltran.Amount=income[i].Amount
		alltran.Paid_Amt=income[i].Paid_Amt
		alltran.Permission=income[i].Permission
		alltran.Debt_Amt=alltran.Amount-alltran.Paid_Amt
		alltran.Status=income[i].Status
		all =append(all, alltran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}
}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Update Transaction")
	t := model.GetAllIncomeTransaction{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	fmt.Println(t.Income_No,"Income id")
	income:=model.OneDay_Transaction{}
	if err:=db.Where("id=?",t.Income_No).Find(&income).Error;err!=nil{
		respondJSON(w,http.StatusNotFound,err.Error())
		return
	}
	fmt.Println(income.Amount,"Amount of Income")
	petrol:=model.Petrol_Info{}//the previous petrol info
	if err:=db.Where("petrol_id=? and access=?",income.Petrol_ID_Fk,true).Find(&petrol).Error;err!=nil{
		respondJSON(w,http.StatusNotFound,err.Error())
		return
	}
	p:=model.Petrol_Info{}
	if err:=db.Where("petrol_name=? and access=?",t.Petrol_Name,true).Find(&p).Error;err!=nil{
		respondJSON(w,http.StatusNotFound,err.Error())
		return
	}
	fmt.Println(petrol.Petrol_Name,"Petrol Name")
	if income.Car_Number!=t.Car_Number {
		book := getBookOr404(db, t.Car_Number, w, r)
		if book.Car_Number != "" {
			income.Customer_ID_Fk = book.Customer_ID
			income.Car_Number = t.Car_Number
			b:=model.Book_Transaction{}
			db.Where("customer_id_fk=?",book.Customer_ID).Last(&b)

			book1 := model.Book_Transaction{}
			if err := db.Where("income_id_fk=?", income.Income_ID).Find(&book1).Error; err != nil {
				respondJSON(w, http.StatusNotFound, err.Error())
				return
			}
			book1.Customer_ID_Fk = income.Customer_ID_Fk
			book1.TLiter=t.Liter
			book1.Total_Liter=b.Total_Liter+t.Liter
			db.Save(&book1)

			d:=model.Debt{}
			db.Where("customer_id_fk=?",book.Customer_ID).Last(&d)
			debt := model.Debt{}
			if err := db.Where("tran_id=?", income.Income_ID).Find(&debt).Error; err != nil {
				respondJSON(w, http.StatusNotFound, err.Error())
				return
			}

			debt.Customer_ID_Fk = income.Customer_ID_Fk
			debt.Car_Number = t.Car_Number
			debt.Debt_Date=t.Income_Date
			if t.Paid_Amt!=0{
				debt.Paid_Date=t.Income_Date
				debt.Debt_Amount=t.Amount-t.Paid_Amt
			}else{
				debt.Debt_Amount=t.Amount
				debt.Paid_Amount=0
			}
			debt.Residual_Amount=debt.Debt_Amount
			debt.Total_Debt=d.Total_Debt+debt.Residual_Amount
			db.Save(&debt)
			db.Save(&income)
		} else if book.Car_Number==""{
			income.Car_Number = t.Car_Number
			income.Customer_ID_Fk=uuid.UUID{00000000-0000-0000-0000-000000000000}
			book1 := model.Book_Transaction{}
			if err := db.Where("income_id_fk=?", income.Income_ID).Find(&book1).Error; err != nil {
				respondJSON(w, http.StatusNotFound, err.Error())
				return
			}
			db.Delete(&book1)
			debt := model.Debt{}
			if err := db.Where("tran_id=?", income.Income_ID).Find(&debt).Error; err != nil {
				respondJSON(w, http.StatusNotFound, err.Error())
				return
			}
			fmt.Println(debt.Debt_ID,"Debt Id ")
			debt.Car_Number = t.Car_Number
			debt.Customer_ID_Fk=uuid.UUID{00000000-0000-0000-0000-000000000000}
			db.Save(&debt)
			db.Save(&income)
		}
	}


	if income.Paid_Amt!=t.Paid_Amt{
		income.Paid_Amt=t.Paid_Amt

		debt:=model.Debt{}
		if err:=db.Where("tran_id=?",income.Income_ID).Find(&debt).Error;err!=nil{
			respondJSON(w,http.StatusNotFound,err.Error())
			return
		}

		debt.Paid_Amount=t.Paid_Amt
		debt.Debt_Amount=t.Amount-t.Paid_Amt
		debt.Total_Debt-=debt.Residual_Amount
		debt.Residual_Amount=debt.Debt_Amount
		debt.Total_Debt+=debt.Residual_Amount
		if t.Amount==debt.Paid_Amount{
			income.Status=true
			if err:=db.Save(&income).Error;err!=nil{
				respondJSON(w,http.StatusNotFound,err.Error())
				return
			}
			debt.Debt_Amount=0
			debt.Total_Debt-=debt.Residual_Amount
			debt.Residual_Amount=0
			debt.Paid_Amount=income.Amount
			debt.Paid_Date=t.Income_Date
			debt.Status=true
			t.Status=true
			if err:=db.Save(&debt).Error;err!=nil{
				respondJSON(w,http.StatusNotFound,err.Error())
				return
			}
			db.Delete(&debt)
		}
		if err:=db.Save(&income).Error;err!=nil{
			respondJSON(w,http.StatusNotFound,err.Error())
			return
		}

	}

	if t.Permission==false && petrol.Petrol_Name!=t.Petrol_Name || income.Liter_Per_Once!=t.Liter{
		fmt.Println("Hello from 1st if")
		book:=model.Book_Transaction{}
		db.Where("income_id_fk=?",income.Income_ID).Find(&book)
		book.Total_Liter-=book.TLiter
		book.TLiter=t.Liter
		book.Total_Liter+=book.TLiter
		db.Save(&book)

		if petrol.Price!=p.Price || income.Liter_Per_Once!=t.Liter{
			income.Liter_Per_Once=t.Liter
			income.Petrol_ID_Fk=p.Petrol_ID
			fmt.Println(p.Price,"Price of Liter for second if")
			t.Amount=t.Liter*p.Price

			if t.Status==true{
				t.Paid_Amt=t.Liter*p.Price
				t.Debt_Amt=0
			}else if t.Status==false{
				fmt.Println("Hello from status false")
				t.Debt_Amt=t.Amount-t.Paid_Amt
				fmt.Println(t.Debt_Amt,"Debt Amount")
				debt:=model.Debt{}
				db.Where("tran_id=? and car_number=?",income.Income_ID,t.Car_Number).Find(&debt)
				if debt.Debt_ID==uuid.Nil{
					debt.Debt_ID,_=uuid.NewV4()
					debt.Tran_ID=income.Income_ID
					debt.Customer_ID_Fk=income.Customer_ID_Fk
					debt.Car_Number=t.Car_Number
					debt.Debt_Date=t.Income_Date
				}

				debt.Total_Debt-=debt.Residual_Amount
				debt.Debt_Amount=t.Amount-t.Paid_Amt
				debt.Paid_Amount=t.Paid_Amt
				debt.Residual_Amount=t.Amount-t.Paid_Amt
				debt.Total_Debt+=debt.Residual_Amount
				if err:=db.Save(&debt).Error;err!=nil{
					respondJSON(w,http.StatusNotFound,err.Error())
					return
				}
				if t.Amount==debt.Paid_Amount{
					debt.Debt_Amount=0
					debt.Total_Debt-=debt.Residual_Amount
					debt.Residual_Amount=0
					debt.Paid_Amount=income.Amount
					debt.Paid_Date=t.Income_Date
					debt.Status=true
					t.Status=true
					if err:=db.Save(&debt).Error;err!=nil{
						respondJSON(w,http.StatusNotFound,err.Error())
						return
					}
					db.Delete(&debt)
				}

			}

		}
		income.Amount=t.Amount
		income.Paid_Amt=t.Paid_Amt
		income.Status=t.Status
		if err:=db.Save(&income).Error;err!=nil{
			respondJSON(w,http.StatusNotFound,err.Error())
			return
		}

	}else if t.Permission==true && p.Price!=t.Petrol_Price || income.Liter_Per_Once!=t.Liter || petrol.Petrol_Name==t.Petrol_Name {
		t.Amount=t.Liter*t.Petrol_Price
		if t.Status==true{
			t.Paid_Amt=t.Liter*t.Petrol_Price
			t.Debt_Amt=0
		}else if t.Status==false{
			t.Debt_Amt=t.Amount-t.Paid_Amt
			debt:=model.Debt{}
			if err:=db.Where("tran_id=?",income.Income_ID).Find(&debt).Error;err!=nil{
				respondJSON(w,http.StatusNotFound,err.Error())
				return
			}
			debt.Total_Debt-=debt.Residual_Amount
			debt.Debt_Amount=t.Amount-t.Paid_Amt
			debt.Paid_Amount=t.Paid_Amt
			debt.Residual_Amount=t.Amount-t.Paid_Amt
			debt.Total_Debt+=debt.Residual_Amount
			if err:=db.Save(&debt).Error;err!=nil{
				respondJSON(w,http.StatusNotFound,err.Error())
				return
			}
			if t.Amount==debt.Paid_Amount{
				debt.Status=true
				t.Status=true
				if err:=db.Save(&debt).Error;err!=nil{
					respondJSON(w,http.StatusNotFound,err.Error())
					return
				}
				db.Delete(&debt)
			}

		}

		income.Amount=t.Amount
		if petrol.Petrol_Name!=t.Petrol_Name{
			income.Low_Price=t.Petrol_Price
		}else{
			income.Low_Price=0
		}

		income.Paid_Amt=t.Paid_Amt
		income.Status=t.Status
		if err:=db.Save(&income).Error;err!=nil{
			respondJSON(w,http.StatusNotFound,err.Error())
			return
		}

	}

	respondJSON(w,http.StatusCreated,"Successfully update transaction")
}

func DeleteTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	id := vars["id"]
	url.QueryEscape(id)
	tran:=model.OneDay_Transaction{}
	if err:=db.Where("id=?",id).Find(&tran).Error;err!=nil{
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	book:=model.Book_Transaction{}
	db.Where(model.Book_Transaction{Income_ID_Fk: tran.Income_ID})
	debt:=model.Debt{}
	db.Where(model.Debt{Tran_ID: tran.Income_ID}).Find(&debt)
	if debt.Car_Number!=""{
		db.Delete(debt)
	}
	if book.Transaction_ID!=uuid.Nil{
		db.Delete(book)
	}
	db.Delete(tran)
	respondJSON(w, http.StatusCreated, "Succesfully Deleted The Transaction")
}

func AmtForNormalCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	p:=model.AmtforNormalPrice{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	pprice:=getPetrolInfoOr404(db,p.Petrol_Name,w,r)
	amt:=pprice.Price*p.Liter
	fmt.Println(amt,"Amount for Normal Price")
	respondJSON(w,http.StatusCreated,amt)
}

func AmtForVIPCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the function for VIP customer")
	p:=model.AmtforVIPPrice{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	amt:=p.Petrol_Price*p.Liter
	fmt.Println(amt,"Amount for Normal Price")
	respondJSON(w,http.StatusCreated,amt)
}

func CreateGiftLiterAndPrice(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	gift := model.Gift{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gift); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	g:=[]model.Gift{}
	if err:=db.Find(&g).Error; err!=nil{
		fmt.Println("There is an error in query")
	}
	if len(g)==0{
		id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("There is error in generating uuid")
		}
		gift.Gift_ID=id
		if err := db.Save(&gift).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}else if len(g)>=1{
		if err := db.Delete(&g).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("There is error in generating uuid")
		}
		gift.Gift_ID=id
		if err := db.Save(&gift).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondJSON(w, http.StatusCreated, gift)
}

func GetGiftLiterAndPrice(db *gorm.DB, w http.ResponseWriter, r *http.Request) model.Gift {
	gift := model.Gift{}
	db.Find(&gift)
	respondJSON(w, http.StatusOK, gift)
	return gift
}

func getGiftOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) model.Gift {
	gift := model.Gift{}
	db.Find(&gift)
	return gift
}

func GetAllDebt(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	debt:=[]model.Debt{}
	if err:=db.Where("status=?",false).Find(&debt).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}

	all:=[]model.GetAllDebt{}
	for i:=0;i<len(debt);i++{
		debtTran:=model.GetAllDebt{}
		debtTran.Debt_No=int(debt[i].ID)
		debtTran.Debt_Date=debt[i].Debt_Date
		debtTran.Car_Number=debt[i].Car_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",debt[i].Tran_ID).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		debtTran.Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			debtTran.Petrol_Price=income.Low_Price
		}else{
			debtTran.Petrol_Price=petrol.Price
		}
		debtTran.Liter=income.Liter_Per_Once
		debtTran.Amount=income.Amount
		debtTran.Paid_Date=debt[i].Paid_Date
		debtTran.Paid_Amt=debt[i].Paid_Amount
		debtTran.Debt_Amt=debt[i].Debt_Amount
		debtTran.Total_Debt=debt[i].Total_Debt
		debtTran.Status=debt[i].Status
		all=append(all,debtTran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}

}

func GetAllDebtForToday(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	debt:=[]model.Debt{}

	date:=time.Now()
	res:=strings.SplitAfterN(date.String()," ",2)
	fmt.Println(date,"Date Format")
	fmt.Println(res[0],"Date split")

	startDate:=res[0]+" 00:00:00"
	endDate:=res[0]+" 23:59:59"
	layout := "2006-01-02 15:04:05"
	tstart, err:= time.Parse(layout, startDate)
	tend,err:=time.Parse(layout,endDate)
	fmt.Println(tstart,"Start Date and Time")
	fmt.Println(tend,"End Date and Time")
	if err!=nil{
		respondJSON(w,http.StatusNotFound,err.Error())
	}
	db.Where("debt_date between ? and ?",startDate,endDate).Find(&debt)
	/*if err:=db.Where("status=?",false).Find(&debt).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}*/

	all:=[]model.GetAllDebt{}
	for i:=0;i<len(debt);i++{
		debtTran:=model.GetAllDebt{}
		debtTran.Debt_No=int(debt[i].ID)
		debtTran.Debt_Date=debt[i].Debt_Date
		debtTran.Car_Number=debt[i].Car_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",debt[i].Tran_ID).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		debtTran.Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			debtTran.Petrol_Price=income.Low_Price
		}else{
			debtTran.Petrol_Price=petrol.Price
		}
		debtTran.Liter=income.Liter_Per_Once
		debtTran.Amount=income.Amount
		debtTran.Paid_Date=debt[i].Paid_Date
		debtTran.Paid_Amt=debt[i].Paid_Amount
		debtTran.Debt_Amt=debt[i].Debt_Amount
		debtTran.Total_Debt=debt[i].Total_Debt
		debtTran.Status=debt[i].Status
		all=append(all,debtTran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}

}

func GetAllDebtOnDay(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	debt:=[]model.Debt{}

	vars := mux.Vars(r)

	date := vars["debt_date"]
	url.QueryEscape(date)
	fmt.Println(date, "Date from UI")
	startDate:=date+" 00:00:00"
	endDate:=date+" 23:59:59"
	layout := "2006-01-02 15:04:05"
	tstart, err:= time.Parse(layout, startDate)
	tend,err:=time.Parse(layout,endDate)
	fmt.Println(tstart,"Start Date and Time")
	fmt.Println(tend,"End Date and Time")
	if err != nil {
		fmt.Println(err)
	}


	if err:=db.Where("debt_date between ? and ?",startDate,endDate).Find(&debt).Error;err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}

	all:=[]model.GetAllDebt{}
	for i:=0;i<len(debt);i++{
		debtTran:=model.GetAllDebt{}
		debtTran.Debt_No=int(debt[i].ID)
		debtTran.Debt_Date=debt[i].Debt_Date
		debtTran.Car_Number=debt[i].Car_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",debt[i].Tran_ID).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		debtTran.Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			debtTran.Petrol_Price=income.Low_Price
		}else{
			debtTran.Petrol_Price=petrol.Price
		}
		debtTran.Liter=income.Liter_Per_Once
		debtTran.Amount=income.Amount
		debtTran.Paid_Date=debt[i].Paid_Date
		debtTran.Paid_Amt=debt[i].Paid_Amount
		debtTran.Debt_Amt=debt[i].Debt_Amount
		debtTran.Total_Debt=debt[i].Total_Debt
		debtTran.Status=debt[i].Status
		all=append(all,debtTran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}

}

func GetAllDebtbyCarNo(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	debt:=[]model.Debt{}
	vars := mux.Vars(r)

	car := vars["car_number"]
	url.QueryEscape(car)
	fmt.Println(car, "car number")
	if err:=db.Where("car_number=?",car).Find(&debt).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}

	all:=[]model.GetAllDebt{}
	for i:=0;i<len(debt);i++{
		debtTran:=model.GetAllDebt{}
		debtTran.Debt_No=int(debt[i].ID)
		debtTran.Debt_Date=debt[i].Debt_Date
		debtTran.Car_Number=debt[i].Car_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",debt[i].Tran_ID).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		debtTran.Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			debtTran.Petrol_Price=income.Low_Price
		}else{
			debtTran.Petrol_Price=petrol.Price
		}
		debtTran.Liter=income.Liter_Per_Once
		debtTran.Amount=income.Amount
		debtTran.Paid_Date=debt[i].Paid_Date
		debtTran.Paid_Amt=debt[i].Paid_Amount
		debtTran.Debt_Amt=debt[i].Debt_Amount
		debtTran.Total_Debt=debt[i].Total_Debt
		debtTran.Status=debt[i].Status
		all=append(all,debtTran)
	}
	if len(all)!=0 {
		respondJSON(w, http.StatusCreated, &all)
	}else{
		respondJSON(w,http.StatusCreated,"There is no transaction for today")
	}

}

func UpdateDebt(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	debt:=model.GetAllDebt{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&debt); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	d:=model.Debt{}
	if err:=db.Where("id=? and car_number=?",debt.Debt_No,debt.Car_Number).Find(&d).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	if debt.Debt_Date!=d.Debt_Date{
		d.Debt_Date=debt.Debt_Date
		db.Save(&d)
	}
	if debt.Paid_Amt!=d.Paid_Amount{
		d.Total_Debt-=d.Debt_Amount
		d.Paid_Amount=debt.Paid_Amt
		d.Debt_Amount=debt.Amount-debt.Paid_Amt
		d.Residual_Amount=d.Debt_Amount
		d.Total_Debt+=d.Debt_Amount
		d.Paid_Date=debt.Paid_Date
		if d.Paid_Amount==d.Debt_Amount{
			d.Debt_Amount=0
			d.Paid_Amount=debt.Amount
			d.Residual_Amount=0
			d.Total_Debt-=d.Debt_Amount
			d.Status=true
			db.Delete(&d)
		}
	}
	if debt.Paid_Date!=d.Paid_Date{
		d.Paid_Date= debt.Paid_Date
		db.Save(&d)
	}
	respondJSON(w, http.StatusCreated,"Successfully Updated to Debt Transaction" )

}

func PaidOnDebt(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	paidDebt:=model.PaidOnDebt{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&paidDebt); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	income:=model.Income{}
	debt:=[]model.Debt{}
	if err:=db.Where("car_number=? and debt_date=? and status=?",paidDebt.Car_Number,paidDebt.Debt_Date,false).Find(&debt).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	if paidDebt.Full_Amt==true{
		newid,_:=uuid.NewV4()
		d:=model.Debt{}
		d.Debt_ID=newid
		d.Car_Number=debt[0].Car_Number
		d.Customer_ID_Fk=debt[0].Customer_ID_Fk
		d.Debt_Amount=paidDebt.Residual_Amt
		d.Debt_Date=paidDebt.Debt_Date
		d.Paid_Date=paidDebt.Paid_Date
		d.Paid_Amount=d.Residual_Amount
		d.Total_Debt=debt[len(debt)-1].Total_Debt-d.Residual_Amount
		d.Residual_Amount=0
		d.Status=true
		d.Tran_ID=debt[0].Tran_ID
		db.Save(&d)
		for i:=0;i<len(debt);i++{
			debt[i].Status=true
			income.Income_Model.Debt_Id=append(income.Income_Model.Debt_Id,debt[i].Debt_ID)
			db.Save(&debt[i])
		}
		newIncome_id,_:=uuid.NewV4()
		income.Income_Id=newIncome_id
		income.Date=paidDebt.Paid_Date
		income.Info="Paid On Debt"
		income.Income_Model.Tran_Id=debt[0].Tran_ID
		income.Amount=paidDebt.Residual_Amt
		db.Save(&income)

	}else if paidDebt.Some_Amt==true{
		newid,_:=uuid.NewV4()
		d:=model.Debt{}
		d.Debt_ID=newid
		d.Car_Number=debt[0].Car_Number
		d.Customer_ID_Fk=debt[0].Customer_ID_Fk
		d.Debt_Amount=debt[len(debt)-1].Residual_Amount
		d.Debt_Date=paidDebt.Debt_Date
		d.Paid_Date=paidDebt.Paid_Date
		d.Paid_Amount=paidDebt.Paid_Amt
		d.Total_Debt=debt[len(debt)-1].Total_Debt-d.Paid_Amount
		d.Residual_Amount=d.Debt_Amount-d.Paid_Amount
		if d.Residual_Amount!=0{
			d.Status=false
		}
		d.Tran_ID=debt[0].Tran_ID
		db.Save(&d)
		for i:=0;i<len(debt);i++{
			income.Income_Model.Debt_Id=append(income.Income_Model.Debt_Id,debt[i].Debt_ID)
			debt[i].Status=true
			db.Save(&debt[i])
		}
		newIncome_id,_:=uuid.NewV4()
		income.Income_Id=newIncome_id
		income.Date=paidDebt.Paid_Date
		income.Info="Paid On Debt"
		income.Income_Model.Tran_Id=debt[0].Tran_ID
		income.Amount=paidDebt.Paid_Amt
		db.Save(&income)
	}
	respondJSON(w, http.StatusCreated,"Successfully Paid on Debt" )
}

func GetBookTransaction (db *gorm.DB, w http.ResponseWriter, r *http.Request){
	booktran:=[]model.Book_Transaction{}
	alltran:=[]model.GetAllBookTran{}
	if err:=db.Find(&booktran).Error; err!=nil{
		respondJSON(w, http.StatusCreated,"There is no book transaction!" )
		return
	}
	for i:=0;i<len(booktran);i++{
		book:=model.Customer_Info{}
		if err:=db.Where("customer_id=? and access=?",booktran[i].Customer_ID_Fk,true).Find(&book).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		alltran[i].Car_Number=book.Car_Number
		alltran[i].Phone=book.Phone_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",booktran[i].Income_ID_Fk).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		alltran[i].Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			alltran[i].Petrol_Price=income.Low_Price
		}else{
			alltran[i].Petrol_Price=petrol.Price
		}
		alltran[i].Date=booktran[i].TDate
		alltran[i].Liter=booktran[i].TLiter
		alltran[i].Total_Liter=booktran[i].Total_Liter
		alltran[i].Status=booktran[i].Status
	}
	if len(alltran)>0{
		respondJSON(w, http.StatusCreated,alltran)
	}else{
		respondJSON(w, http.StatusNotFound,"There is no book transaction!" )
	}

}

func GetBookTransactionByCarNo (db *gorm.DB, w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	car := vars["car_number"]
	url.QueryEscape(car)
	fmt.Println(car, "car number")
	book:=model.Customer_Info{}
	if err:=db.Where("car_number=? and access=?",car,true).Find(&book).Error; err!=nil{
		respondJSON(w, http.StatusNotFound,err.Error() )
		return
	}
	booktran:=[]model.Book_Transaction{}
	alltran:=[]model.GetAllBookTran{}
	if err:=db.Where("customer_id_fk=?",book.Customer_ID).Find(&booktran).Error; err!=nil{
		respondJSON(w, http.StatusCreated,"There is no book transaction!" )
		return
	}
	for i:=0;i<len(booktran);i++{
		alltran[i].Car_Number=book.Car_Number
		alltran[i].Phone=book.Phone_Number
		income:=model.OneDay_Transaction{}
		if err:=db.Where("income_id=?",booktran[i].Income_ID_Fk).Find(&income).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		petrol:=model.Petrol_Info{}
		if err:=db.Where("petrol_id=?",income.Petrol_ID_Fk).Find(&petrol).Error; err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		alltran[i].Petrol_Name=petrol.Petrol_Name
		if income.Permission==true{
			alltran[i].Petrol_Price=income.Low_Price
		}else{
			alltran[i].Petrol_Price=petrol.Price
		}
		alltran[i].Date=booktran[i].TDate
		alltran[i].Liter=booktran[i].TLiter
		alltran[i].Total_Liter=booktran[i].Total_Liter
		alltran[i].Status=booktran[i].Status
	}
	if len(alltran)>0{
		respondJSON(w, http.StatusCreated,alltran)
	}else{
		respondJSON(w, http.StatusNotFound,"There is no book transaction!" )
	}

}

func GetBookTransactionByGiftLiter (db *gorm.DB, w http.ResponseWriter, r *http.Request){
	gift:=model.Gift{}
	if err:=db.Find(&gift).Error; err!=nil{
		respondJSON(w, http.StatusCreated,"There is no book transaction!" )
		return
	}
	booktran:=[]model.Book_Transaction{}

	book:=model.GetAllBookTranbyGiftLiter{}
	alltran:=[]model.GetAllBookTranbyGiftLiter{}
	cus_id:=[]int{}

	fmt.Println(gift.Gift_Liter,"Gift Liter")
	db.Where("total_liter >=?",gift.Gift_Liter).Find(&booktran).Group("customer_id_fk")
	for i:=0; i< len(booktran);i++{
		if len(cus_id)==0{
			for j:=i;j<len(booktran);j++{
				if booktran[i].Customer_ID_Fk==booktran[j].Customer_ID_Fk{
					book.Total_Liter= booktran[j].Total_Liter
					customer:=model.Customer_Info{}
					if err:=db.Where("customer_id=? and access=?",booktran[j].Customer_ID_Fk,true).Find(&customer).Error; err!=nil{
						respondJSON(w, http.StatusNotFound,err.Error() )
						return
					}
					book.Car_Number=customer.Car_Number
					book.Date=booktran[j].TDate
					book.Phone= customer.Phone_Number
					cus_id=append(cus_id,int(booktran[i].ID))
				}else {
					book.Total_Liter= booktran[j].Total_Liter
					customer:=model.Customer_Info{}
					if err:=db.Where("customer_id=? and access=?",booktran[j].Customer_ID_Fk,true).Find(&customer).Error; err!=nil{
						respondJSON(w, http.StatusNotFound,err.Error() )
						return
					}
					book.Car_Number=customer.Car_Number
					book.Date=booktran[j].TDate
					book.Phone= customer.Phone_Number
					cus_id=append(cus_id,int(booktran[i].ID))
				}

				alltran=append(alltran,book)
			}
		}else{
			for k:=0;k<len(cus_id);k++{
				if int(booktran[i].ID)==cus_id[k]{
					break
				}else{
					for j:=i;j<len(booktran);j++{
						if booktran[i].Customer_ID_Fk==booktran[j].Customer_ID_Fk{
							book.Total_Liter= booktran[j].Total_Liter
							customer:=model.Customer_Info{}
							if err:=db.Where("customer_id=? and access=?",booktran[j].Customer_ID_Fk,true).Find(&customer).Error; err!=nil{
								respondJSON(w, http.StatusNotFound,err.Error() )
								return
							}
							book.Car_Number=customer.Car_Number
							book.Date=booktran[j].TDate
							book.Phone= customer.Phone_Number
							cus_id=append(cus_id,int(booktran[i].ID))
						}else {
							continue
						}

						alltran=append(alltran,book)

					}
				}
			}
		}
		break
	}


	respondJSON(w,http.StatusCreated,&alltran)


}

func CreateGiftTrans (db *gorm.DB, w http.ResponseWriter, r *http.Request){
	gifttran := model.GiftTran{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gifttran); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("There is error in generating uuid")
	}
	gift:=model.GiftTransaction{}
	gift.Gift_Tran_ID=id
	gift.Gift_Tran_Date=gifttran.GT_Date
	gift.Gift_Tran_Car_Number=gifttran.GT_Car_Number
	if gifttran.GT_GoodsORAmt==true {
		g:=model.Gift{}
		db.Find(&g)
		gift.Gift_Tran_Amount=g.Gift_Amount
		gift.Gift_Tran_Liter=g.Gift_Liter
		cus:=model.Customer_Info{}
		fmt.Println(gift.Gift_Tran_Car_Number,"Car Number from UI")
		if err:=db.Where("car_number=? and access=?",gift.Gift_Tran_Car_Number,true).Find(&cus).Error;err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		book:=model.Book_Transaction{}
		if err:=db.Where("customer_id_fk=? and status=?",cus.Customer_ID,true).Last(&book).Error;err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		gift.Gift_Tran_Total_Liter=book.Total_Liter-g.Gift_Liter
		db.Save(&gift)
		book.Total_Liter=book.Total_Liter-g.Gift_Liter
		if book.Total_Liter<g.Gift_Liter{
			book.Status=false
		}
		db.Save(&book)
	}else {
		g:=model.Gift{}
		db.Find(&g)
		gift.Gift_Tran_Goods=gifttran.GT_Goods
		gift.Gift_Tran_Liter=g.Gift_Liter
		cus:=model.Customer_Info{}
		if err:=db.Where("car_number=? and access=?",gift.Gift_Tran_Car_Number,true).Find(&cus).Error;err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}
		fmt.Println(cus.Car_Number,"Car Number")
		book:=model.Book_Transaction{}
		if err:=db.Where("customer_id_fk=? and status=?",cus.Customer_ID,true).Last(&book).Error;err!=nil{
			respondJSON(w, http.StatusNotFound,err.Error() )
			return
		}

		fmt.Println(book.Total_Liter,"Total Liter")
		gift.Gift_Tran_Total_Liter=book.Total_Liter-g.Gift_Liter
		db.Save(&gift)
		book.Total_Liter=book.Total_Liter-g.Gift_Liter
		if book.Total_Liter<g.Gift_Liter{
			book.Status=false
		}
		db.Save(&book)
	}

		fmt.Println("Create new gift Transaction")
		respondJSON(w, http.StatusCreated, &gift)
}

func GetLastBookTransactionByGiftLiter (db *gorm.DB, w http.ResponseWriter, r *http.Request){
	gift:=model.Gift{}
	if err:=db.Find(&gift).Error; err!=nil{
		respondJSON(w, http.StatusCreated,"There is no book transaction!" )
		return
	}
	//booktran:=[]model.Book_Transaction{}

	book:=model.GetAllBookTranbyGiftLiter{}

	alltran:=[]model.GetAllBookTranbyGiftLiter{}
	//cus_id:=[]int{}
	customer:=[]model.Customer_Info{}
	db.Find(&customer)
	for i:=0;i<len(customer);i++{
		lastgift:=model.Book_Transaction{}
		db.Where("customer_id_fk=?",customer[i].Customer_ID).Last(&lastgift)
		if lastgift.Total_Liter<gift.Gift_Liter{
			continue
		}else if lastgift.Total_Liter >= gift.Gift_Liter{
			book.Total_Liter= lastgift.Total_Liter

			book.Car_Number=customer[i].Car_Number
			book.Date=lastgift.TDate
			book.Phone= customer[i].Phone_Number
			alltran=append(alltran,book)
		}

	}

	respondJSON(w,http.StatusCreated,&alltran)


}

func login(db *gorm.DB,w http.ResponseWriter,r *http.Request){
	fmt.Print("method",r.Method)
	if r. Method=="Post"{
		r.ParseForm()
		username:=r.FormValue("name")
		password:=r.FormValue("password")
		fmt.Println("Username",r.Form["name"])
		fmt.Println("Password",r.Form["password"])
		CheckLoginData(db,w,r,username,password)
	}

}

