package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Employee struct {
	gorm.Model
	Name   string `gorm:"unique" json:"name"`
	City   string `json:"city"`
	Age    int    `json:"age"`
	Status bool   `json:"status"`
}
type Petrol_Info struct {
	gorm.Model
	Petrol_ID uuid.UUID `json:"petrol_id" gorm:"primary_key";type:char(36);`
	Petrol_Name string `json:"petrol_name"`
	Detail string `json:"detail"`
	Price float32 `json:"price"`
	Liter float32 `json:"liter"`
	Color string `json:"color"`
	Access bool `json:"access"'`
}

type Gift struct{
	gorm.Model
	Gift_ID uuid.UUID `json:"gift_id" gorm:"primary_key";type:char(36);`
	Gift_Liter float32 `json:"gift_liter"`
	Gift_Amount float32 `json:"gift_amount"`
}

/*type Petrol_Price struct {
	gorm.Model
	Price_ID uuid.UUID `json:"price_id" gorm:"primary_key";type:char(36);`
	Petrol_ID_Fk uuid.UUID `gorm:"ForeignKey:Petrol_ID" json:"petrol_id_fk"`
	Price float32 `json:"price"`
	Liter float32 `json:"liter"`
}*/

type Customer_Info struct {
	gorm.Model
	Customer_ID uuid.UUID `json:"customer_id" gorm:"primary_key";type:char(36);`
	Customer_Name string `json:"customer_name"`
	Car_Type string `json:"car_type"`
	Car_Number string `json:"car_number"`
	NRC string `json:"nrc"`
	Phone_Number string`json:"phone_number"`
	Owner_Name string `json:"owner_name"`
	Owner_Phone string `json:"owner_phone"`
	Access bool `json:"access"'`
	Permission bool `json:"permission"`
}

type OneDay_Transaction struct {
	gorm.Model
	Income_ID uuid.UUID `json:"income_id" gorm:"primary_key";type:char(36);`
	Income_Date time.Time `json:"income_date"`
	Petrol_ID_Fk uuid.UUID `gorm:"ForeignKey:Petrol_ID" json:"petrol_id_fk"`
	Liter_Per_Once float32 `json:"liter_per_once"`
	Car_Number string `json:"car_number"`
	Amount float32 `json:"amount"`
	Paid_Amt float32 `json:"paid_amt"`
	Info string`json:"info"`
	Status bool `json:"status"` //for debt or paid
	Permission bool `json:"permission"`//permission for price change
	Low_Price float32 `json:"lowprice"'`
	Customer_ID_Fk uuid.UUID `gorm:"ForeignKey:Customer_ID" json:"customer_id_fk"`
}

type Book_Transaction struct {
	gorm.Model
	Transaction_ID uuid.UUID `json:"transaction_id" gorm:"primary_key";type:char(36);`
	TDate time.Time `json:"tdate"`
	Petrol_ID_Fk uuid.UUID `gorm:"ForeignKey:Petrol_ID" json:"petrol_id_fk"`
	TLiter float32 `json:"tliter"`
	Total_Liter float32 `json:"total_liter"`
	Income_ID_Fk uuid.UUID `gorm:"ForeignKey:Income_ID" json:"income_id_fk"`
	Customer_ID_Fk uuid.UUID `gorm:"ForeignKey:Customer_ID" json:"customer_id_fk"`
	Status bool `json:"status"`// If status is true, it didn't take out the gift liter
}

type Boucher struct {
	gorm.Model
	Boucher_ID uuid.UUID `json:"boucher_id" gorm:"primary_key";type:char(36);`
	BDate time.Time `json:"bdate"`
	Customer_ID_Fk uuid.UUID `gorm:"ForeignKey:Customer_ID" json:"customer_id_fk"`
	Car_Number string `json:"car_number"`
	Petrol_ID_Fk uuid.UUID `gorm:"ForeignKey:Petrol_ID" json:"petrol_id_fk"`
	BAmount float32 `json:"bamount"`
	BLiter float32 `json:"bliter"`
	Status bool `json:"status"`
}

type Debt struct {
	gorm.Model
	Debt_ID uuid.UUID `json:"debt_id" gorm:"primary_key";type:char(36);`
	Customer_ID_Fk uuid.UUID `gorm:"ForeignKey:Customer_ID" json:"customer_id_fk"`
	Car_Number string `json:"car_number"`
	Debt_Amount float32 `json:"debt_amount"`
	Debt_Date time.Time `json:"debt_date"`
	Paid_Date time.Time `json:"paid_date"`
	Paid_Amount float32 `json:"paid_amount"`
	Residual_Amount float32 `json:"residual_amount"`
	Total_Debt float32 `json:"total_debt"`
	Status bool `json:"status"`
	Tran_ID uuid.UUID `json:"tran_id"`
}

type SignUp struct {
	gorm.Model
	SignUp_ID uuid.UUID `json:"signup_id" gorm:"primary_key";type:char(36);`
	SName string `json:"sname"`
	SPassword string `json:"spassword"`
	ReEnter_Password string `json:"reenter_password"`
	BirthDate time.Time`json:"birthdate"`
	Email string `json:"email"`
	PhoneNo string `json:"phone_no"`
}

type LogIn struct {
	gorm.Model
	Login_ID uuid.UUID `json:"login_id" gorm:"primary_key";type:char(36);`
	UserName string `json:"username"`
	Password string `json:"password"`
	SignUp_ID_Fk uuid.UUID  `gorm:"ForeignKey:SignUp_ID" json:"signup_id_fk"`
	Access bool `json:"access"`

}

type Expense struct {
	gorm.Model
	Exp_ID uuid.UUID `json:"exp_id" gorm:"primary_key";type:char(36);`
	Date time.Time `json:"date"`
	Info string `json:"info"`
	Amount float32 `json:"amount"`
}

type GiftTransaction struct {
	gorm.Model
	Gift_Tran_ID uuid.UUID `json:"gift_tran_id" gorm:"primary_key";type:char(36);`
	Gift_Tran_Date time.Time `json:"gift_tran_date"`
	Gift_Tran_Car_Number string `json:"gift_tran_car_number"`
	Gift_Tran_Goods string `json:"gift_tran_goods"`
	Gift_Tran_Liter float32 `json:"gift_tran_liter"`
	Gift_Tran_Total_Liter float32 `json:"gift_tran_total_liter"`
	Gift_Tran_Amount float32 `json:"gift_tran_amt"`
}

type Transaction struct {
	Income_Date time.Time `json:"income_date"`
	Petrol_Name string `json:"petrol_name"`
	Liter float32 `json:"liter"`
	Car_Number string `json:"car_number"`
	Debt bool `json:"debt"`
	Paid bool `json:"paid"`
	Both bool `json:"both"`
	Gift bool `json:"gift"`
	Permission bool `json:"permission"`
	LowPrice float32 `json:"lowprice"`
	Amount float32 `json:"amount"`
	Info string `json:"info"`
	Debt_Amt float32 `json:"debt_amt"`
	Paid_Amt float32 `json:"paid_amt"`
	Residual_Amt float32 `json:"residual_amt"`
}

type GetAllIncomeTransaction struct {
	Income_No int `json:"income_no"`
	Income_Date time.Time `json:"income_date"`
	Car_Number string `json:"car_number"`
	Petrol_Name string `json:"petrol_name"`
	Petrol_Price float32 `json:"petrol_price"`
	Liter float32 `json:"liter"`
	Amount float32 `json:"amount"`
	Paid_Amt float32 `json:"paid_amt"`
	Debt_Amt float32 `json:"debt_amt"`
	Permission bool `json:"permission"`
	Status bool `json:"status"`
}


type GetAllBookTran struct {
	Car_Number string `json:"car_number"`
	Phone string `json:"phone"`
	Date time.Time `json:"date"`
	Petrol_Name string `json:"petrol_name"`
	Petrol_Price float32 `json:"petrol_price"`
	Liter float32 `json:"liter"`
	Total_Liter float32 `json:"total_liter"`
	Status bool `json:"status"`
}

type GetAllBookTranbyGiftLiter struct {
	Car_Number string `json:"car_number"`
	Phone string `json:"phone"`
	Date time.Time `json:"date"`
	Total_Liter float32 `json:"total_liter"`
}

type GetAllDebt struct {
	Debt_No int `json:"debt_no"`
	Debt_Date time.Time `json:"debt_date"`
	Car_Number string `json:"car_number"`
	Petrol_Name string `json:"petrol_name"`
	Petrol_Price float32 `json:"petrol_price"`
	Liter float32 `json:"liter"`
	Amount float32 `json:"amount"`
	Paid_Date time.Time`json:"paid_date"`
	Paid_Amt float32 `json:"paid_amt"`
	Debt_Amt float32 `json:"debt_amt"`
	Total_Debt float32 `json:"total_debt"`
	Status bool `json:"status"`
}

type PaidOnDebt struct {
	Car_Number string `json:"car_number"`
	Debt_Date time.Time `json:"debt_date"`
	Paid_Date time.Time `json:"paid_date"`
	Residual_Amt float32 `json:"residual_amt"`
	Full_Amt bool `json:"full_amt"`
	Some_Amt bool `json:"some_amt"`
	Paid_Amt float32 `json:"paid_amt"`
}

type Income struct {
	gorm.Model
	Income_Id uuid.UUID `json:"income_id" gorm:"primary_key";type:char(36);`
	Date time.Time `json:"income_date"`
	Info string `json:"info"`
	Income_Model Debt_Model `gorm:"type:json"`
	Amount float32 `json:"amount"`
}

type Debt_Model struct {
	Tran_Id uuid.UUID `json:"tran_id"'`
	Debt_Id []uuid.UUID `json:"debt_id"`
}

type AmtforNormalPrice struct{
	Petrol_Name string `json:"petrol_name"`
	Liter float32 `json:"liter"`
}

type GiftTran struct {
	GT_Date time.Time `json:"gt_date"`
	GT_Car_Number string `json:"gt_car_number"`
	GT_Goods string `json:"gt_goods"`
	GT_Amount float32 `json:"gt_amt"`
	GT_GoodsORAmt bool `json:"gt_goodsoramt"`
}

type AmtforVIPPrice struct{
	Petrol_Price float32 `json:"petrol_price"`
	Liter float32 `json:"liter"`
}

func (e *Employee) Disable() {
	e.Status = false
}

func (p *Employee) Enable() {
	p.Status = true
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Employee{},&Petrol_Info{},&Customer_Info{},&Book_Transaction{},&Boucher{},&Debt{},&OneDay_Transaction{},&SignUp{},&LogIn{},&Gift{},&Expense{},&GiftTransaction{},&Income{})
	return db
}
