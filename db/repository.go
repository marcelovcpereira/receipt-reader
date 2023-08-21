package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PurchaseRepository struct {
	db *sqlx.DB
}

func NewPurchaseRepository(db *sqlx.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (repository PurchaseRepository) GetAll() ([]*Purchase, error) {
	var purchases []*Purchase
	if err := repository.db.Select(&purchases, "SELECT * FROM purchases"); err != nil {
		return nil, err
	}
	if len(purchases) == 0 {
		return make([]*Purchase, 0), nil
	}

	return purchases, nil
}

func (repository PurchaseRepository) Get(name string) (*Purchase, error) {
	var services []Purchase
	err := repository.db.Select(&services, "SELECT * FROM services WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, nil
	}

	service := services[0]
	return &service, nil
}

func (repository PurchaseRepository) Save(service *Purchase) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback() // safe to call on a committed tnx

	if err := repository.insertPurchase(tx, service); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repository PurchaseRepository) insertPurchase(tx *sqlx.Tx, service *Purchase) error {
	_, err := tx.NamedExec("INSERT INTO services (name, type, last_updated) VALUES (:name, :type, :last_updated) ON CONFLICT ON CONSTRAINT services_name_key DO UPDATE SET type = :type, last_updated = :last_updated", service)
	if err != nil {
		return err
	}

	// get service Id
	err = tx.Get(service, "SELECT * FROM services WHERE name = $1", service.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Saving service Id: %d, Name: %s. Total Components: %d", service.Id, service.Name, len(service.Components))

	return nil
}

func (repository PurchaseRepository) delete(tx *sqlx.Tx, purchase Purchase) error {
	_, err := tx.Exec("DELETE FROM purchases WHERE name = $1", purchase.Name)
	if err != nil {
		return err
	}

	return nil
}
