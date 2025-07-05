package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Usuario     UserModel
	Permissions PermissionModel
	Autor       AutorModel
	Book        BookModel
	Editorial   EditorialModel
	Fine        FineModel
	Loan        LoanModel
	Reservation ReservationModel
	User_       UserModel_
}

func NewModels(db *sql.DB) Models {
	return Models{
		Usuario:     UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Autor:       AutorModel{DB: db},
		Book:        BookModel{DB: db},
		Editorial:   EditorialModel{DB: db},
		Fine:        FineModel{DB: db},
		Loan:        LoanModel{DB: db},
		Reservation: ReservationModel{DB: db},
		User_:       UserModel_{DB: db},
	}
}
