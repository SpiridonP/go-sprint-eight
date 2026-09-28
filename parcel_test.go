package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, number)
	parcel.Number = number

	// get
	data, err := store.Get(number)
	assert.NoError(t, err)

	expected := parcel
	expected.Number = number
	assert.Equal(t, expected, data)

	// delete
	err = store.Delete(number)
	assert.NoError(t, err)

	_, err = store.Get(number)
	assert.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, number)
	parcel.Number = number

	// set address
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	assert.NoError(t, err)

	// check
	data, err := store.Get(number)
	assert.NoError(t, err)

	expected := parcel
	expected.Number = number
	expected.Address = newAddress
	assert.Equal(t, expected, data)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, number)
	parcel.Number = number

	// set status
	newStatus := ParcelStatusRegistered
	err = store.SetStatus(number, newStatus)
	assert.NoError(t, err)

	// check
	data, err := store.Get(number)
	assert.NoError(t, err)

	expected := parcel
	expected.Number = number
	expected.Status = newStatus
	assert.Equal(t, expected, data)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		expected, ok := parcelMap[parcel.Number]
		assert.True(t, ok)

		assert.Equal(t, expected, parcel)
	}
}
