// Code generated by boltbindings generator. DO NOT EDIT.

package store

import (
	bbolt "github.com/etcd-io/bbolt"
	storage "github.com/stackrox/rox/generated/storage"
	storecache "github.com/stackrox/rox/pkg/storecache"
)

type Store interface {
	DeleteServiceAccount(id string) error
	GetServiceAccount(id string) (*storage.ServiceAccount, bool, error)
	GetServiceAccounts(ids []string) ([]*storage.ServiceAccount, []int, error)
	ListServiceAccounts() ([]*storage.ServiceAccount, error)
	UpsertServiceAccount(serviceaccount *storage.ServiceAccount) error
}

func New(db *bbolt.DB, cache storecache.Cache) (Store, error) {
	return newStore(db, cache)
}
