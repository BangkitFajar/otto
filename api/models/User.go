package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nama      string    `gorm:"size:255;not null;unique" json:"nama"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Result struct {
	ID      uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Nama    string  `gorm:"size:255;not null;unique" json:"nama"`
	Email   string  `gorm:"size:100;not null;unique" json:"email"`
	Balance Balance `gorm:"size:100;not null;unique" json:"balance`
}

type Balance struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserId    uint32    `gorm:"size:255;not null" json:"user_id"`
	Saldo     string    `gorm:"not null;" json:"saldo"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Status    string    `gorm:"default:" json:"status"`
}

type Results struct {
	ID      uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nama    string    `gorm:"size:255;not null;unique" json:"nama"`
	Email   string    `gorm:"size:100;not null;unique" json:"email"`
	Balance []Balance `gorm:"size:100;not null;unique" json:"balance`
}

type Paymentmethod struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserId    uint32    `gorm:"size:255;not null" json:"user_id"`
	Nama      string    `gorm:"nama" json:"nama"`
	VA        string    `gorm:"va" json:"va"`
	Status    string    `gorm:"status" json:"status"`
	Nominal   string    `gorm:"nominal" json:"nominal"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (p *Paymentmethod) SavePaymentMethod(db *gorm.DB) (*Paymentmethod, error) {
	var err error
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &Paymentmethod{}, err
	}
	return p, nil
}

func (b *Balance) SaveBalance(db *gorm.DB) (*Balance, error) {
	// log.Fatal(*b)
	var err error
	err = db.Debug().Create(&b).Error
	if err != nil {
		return &Balance{}, err
	}
	return b, nil
}

func (b *Balance) FindAllBalance(db *gorm.DB, uid uint32) (*[]Balance, error) {
	var err error
	balances := []Balance{}
	err = db.Debug().Model(&Balance{}).Where("user_id = ?", uid).Where("status != ?", "aktif").Limit(100).Find(&balances).Error
	if err != nil {
		return &[]Balance{}, err
	}
	return &balances, err
}

func (b *Balance) FindBalanceByID(db *gorm.DB, uid uint32) (*Balance, error) {
	var err error
	err = db.Debug().Model(Balance{}).Where("user_id = ?", uid).Take(&b).Error
	if err != nil {
		return &Balance{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Balance{}, errors.New("User Not Found")
	}
	return b, err
}

func (b *Balance) GetLatestBalance(db *gorm.DB, uid uint32) (*Balance, error) {
	var err error
	err = db.Debug().Model(&Balance{}).Where("user_id = ?", uid).Order("id desc, id").Limit(1).Take(&b).Error
	// db.Order("age desc, name").Find(&users)
	if err != nil {
		return &Balance{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Balance{}, errors.New("User Not Found")
	}
	return b, err
}

func (p *Paymentmethod) FindBalanceByVA(db *gorm.DB, va string) (*Paymentmethod, error) {
	var err error
	err = db.Debug().Model(Balance{}).Where("va = ?", va).Take(&p).Error
	if err != nil {
		return &Paymentmethod{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Paymentmethod{}, errors.New("Number VA Not Found")
	}
	return p, err
}

func (p *Paymentmethod) UpdatePayment(db *gorm.DB, uid uint32) (*Paymentmethod, error) {

	db = db.Debug().Model(&Paymentmethod{}).Where("id = ?", uid).Take(&Paymentmethod{}).UpdateColumns(
		map[string]interface{}{
			"Status":    p.Status,
			"UpdatedAt": time.Now(),
		},
	)
	if db.Error != nil {
		return &Paymentmethod{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Paymentmethod{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Paymentmethod{}, err
	}
	return p, nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nama = html.EscapeString(strings.TrimSpace(u.Nama))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nama == "" {
			return errors.New("Required Nama")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nama == "" {
			return errors.New("Required Nama")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}
