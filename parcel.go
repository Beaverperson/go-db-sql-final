package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at)"+
			"VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Printf("error while inserting parcel to DB\n%v", err)
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error while inserting parcel to DB\n%v", err)
		return 0, err
	}
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number",
		sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Printf("error while selecting client (%d) from DB\n%v", client, err)
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		bufferPars := Parcel{}
		err := rows.Scan(
			&bufferPars.Number,
			&bufferPars.Client,
			&bufferPars.Status,
			&bufferPars.Address,
			&bufferPars.CreatedAt)
		if err != nil {
			fmt.Printf("error while rows iteration for client (%d) from DB\n%v", client, err)
			return nil, err
		}
		res = append(res, bufferPars)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error after rows iteration for client (%d) from DB\n%v", client, err)
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Printf("error while updating status for parcel(%d)\n%v", number, err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Printf("error while setting new addres for parcel(%d)\n%v", number, err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	targetParcel, err := s.Get(number)
	if err != nil {
		// дебажное кукареку в консоль реализовано в методе GET структуры ParcelStore
		return err
	}
	if targetParcel.Status != ParcelStatusRegistered {
		fmt.Printf("attempt to delete sent/delivered parcel number: %d\n", number)
		return fmt.Errorf("attempt to delete sent/delivered parcel number: %d\n", number)
	}
	_, err = s.db.Exec(`DELETE FROM parcel WHERE number = :number`,
		sql.Named("number", number))
	if err != nil {
		fmt.Printf("error while deleting parcel (%d) from DB\n%v", number, err)
		return err
	}
	return nil
}
