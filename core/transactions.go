package core

import (
	"errors"
	"os"
	"path"

	"github.com/zyrouge/pho/utils"
)

type PendingInstallation struct {
	InvolvedDirs  []string
	InvolvedFiles []string
}

type Transactions struct {
	PendingInstallations map[string]PendingInstallation
}

func GetTransactionsPath() (string, error) {
	xdgConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	transtionsPath := path.Join(xdgConfigDir, AppCodeName, "transactions.json")
	return transtionsPath, nil
}

func GetTransactions() (*Transactions, error) {
	transtionsPath, err := GetTransactionsPath()
	if err != nil {
		return nil, err
	}
	transactions, err := utils.ReadJsonFile[Transactions](transtionsPath)
	if errors.Is(err, os.ErrNotExist) {
		transactions = &Transactions{
			PendingInstallations: map[string]PendingInstallation{},
		}
		err = nil
	}
	return transactions, err
}

func SaveTransactions(transactions *Transactions) error {
	transtionsPath, err := GetTransactionsPath()
	if err != nil {
		return err
	}
	err = utils.WriteJsonFileAtomic[Transactions](transtionsPath, transactions)
	if err != nil {
		return err
	}
	return nil
}

type UpdateTransactionFunc func(transactions *Transactions) error

func UpdateTransactions(performer UpdateTransactionFunc) error {
	transactions, err := GetTransactions()
	if err != nil {
		return err
	}
	err = performer(transactions)
	if err != nil {
		return err
	}
	return SaveTransactions(transactions)
}
