package db

import (
	"fmt"
	"sync"

	cmn "github.com/dominonetwork/spacex/libs/common"
)

//----------------------------------------
// debugDB

type debugDB struct {
	label string
	db    DB
}

// For printing all operationgs to the console for debugging.
func NewDebugDB(label string, db DB) debugDB {
	return debugDB{
		label: label,
		db:    db,
	}
}

// Implements atomicSetDeleter.
func (ddb debugDB) Mutex() *sync.Mutex { return nil }

func (ddb debugDB) Dir() string { return string("") }

// Implements DB.
func (ddb debugDB) Get(key []byte) (value []byte) {
	defer func() {
		fmt.Printf("%v.Get(%v) %v\n", ddb.label,
			cmn.ColoredBytes(key, cmn.Cyan, cmn.Blue),
			cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	}()
	value = ddb.db.Get(key)
	return
}
func (ddb debugDB) Load(key []byte) (value []byte, err error) {
	defer func() {
		fmt.Printf("%v.Get(%v) %v\n", ddb.label,
			cmn.ColoredBytes(key, cmn.Cyan, cmn.Blue),
			cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	}()
	return ddb.db.Load(key)
}

// Implements DB.
func (ddb debugDB) Has(key []byte) (has bool) {
	defer func() {
		fmt.Printf("%v.Has(%v) %v\n", ddb.label,
			cmn.ColoredBytes(key, cmn.Cyan, cmn.Blue), has)
	}()
	return ddb.db.Has(key)
}
func (ddb debugDB) Exist(key []byte) (has bool, err error) {
	defer func() {
		fmt.Printf("%v.Has(%v) %v\n", ddb.label,
			cmn.ColoredBytes(key, cmn.Cyan, cmn.Blue), has)
	}()
	return ddb.db.Exist(key)
}

// Implements DB.
func (ddb debugDB) Set(key []byte, value []byte) {
	fmt.Printf("%v.Set(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	ddb.db.Set(key, value)
}
func (ddb debugDB) Put(key []byte, value []byte) error {
	fmt.Printf("%v.Set(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	return ddb.db.Put(key, value)
}

// Implements DB.
func (ddb debugDB) SetSync(key []byte, value []byte) {
	fmt.Printf("%v.SetSync(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	ddb.db.SetSync(key, value)
}

// Implements atomicSetDeleter.
func (ddb debugDB) SetNoLock(key []byte, value []byte) {
	fmt.Printf("%v.SetNoLock(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	ddb.db.(atomicSetDeleter).SetNoLock(key, value)
}

// Implements atomicSetDeleter.
func (ddb debugDB) SetNoLockSync(key []byte, value []byte) {
	fmt.Printf("%v.SetNoLockSync(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	ddb.db.(atomicSetDeleter).SetNoLockSync(key, value)
}

// Implements DB.
func (ddb debugDB) Delete(key []byte) {
	fmt.Printf("%v.Delete(%v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	ddb.db.Delete(key)
}
func (ddb debugDB) Del(key []byte) error {
	fmt.Printf("%v.Delete(%v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	return ddb.db.Del(key)
}

// Implements DB.
func (ddb debugDB) DeleteSync(key []byte) {
	fmt.Printf("%v.DeleteSync(%v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	ddb.db.DeleteSync(key)
}

// Implements atomicSetDeleter.
func (ddb debugDB) DeleteNoLock(key []byte) {
	fmt.Printf("%v.DeleteNoLock(%v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	ddb.db.(atomicSetDeleter).DeleteNoLock(key)
}

// Implements atomicSetDeleter.
func (ddb debugDB) DeleteNoLockSync(key []byte) {
	fmt.Printf("%v.DeleteNoLockSync(%v)\n", ddb.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	ddb.db.(atomicSetDeleter).DeleteNoLockSync(key)
}

// Implements DB.
func (ddb debugDB) Iterator(start, end []byte) Iterator {
	fmt.Printf("%v.Iterator(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(start, cmn.Cyan, cmn.Blue),
		cmn.ColoredBytes(end, cmn.Cyan, cmn.Blue))
	return NewDebugIterator(ddb.label, ddb.db.Iterator(start, end))
}

func (ddb debugDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	end := PrefixToEnd(prefix)
	fmt.Printf("%v.Iterator(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(prefix, cmn.Cyan, cmn.Blue),
		cmn.ColoredBytes(end, cmn.Cyan, cmn.Blue))
	return NewDebugIterator(ddb.label, ddb.db.Iterator(prefix, end))
}

// Implements DB.
func (ddb debugDB) ReverseIterator(start, end []byte) Iterator {
	fmt.Printf("%v.ReverseIterator(%v, %v)\n", ddb.label,
		cmn.ColoredBytes(start, cmn.Cyan, cmn.Blue),
		cmn.ColoredBytes(end, cmn.Cyan, cmn.Blue))
	return NewDebugIterator(ddb.label, ddb.db.ReverseIterator(start, end))
}

// Implements DB.
// Panics if the underlying db is not an
// atomicSetDeleter.
func (ddb debugDB) NewBatch() Batch {
	fmt.Printf("%v.NewBatch()\n", ddb.label)
	return NewDebugBatch(ddb.label, ddb.db.NewBatch())
}

// Implements DB.
func (ddb debugDB) Close() {
	fmt.Printf("%v.Close()\n", ddb.label)
	ddb.db.Close()
}

// Implements DB.
func (ddb debugDB) Print() {
	ddb.db.Print()
}

// Implements DB.
func (ddb debugDB) Stats() map[string]string {
	return ddb.db.Stats()
}

//----------------------------------------
// debugIterator

type debugIterator struct {
	label string
	itr   Iterator
}

// For printing all operationgs to the console for debugging.
func NewDebugIterator(label string, itr Iterator) debugIterator {
	return debugIterator{
		label: label,
		itr:   itr,
	}
}

// Implements Iterator.
func (ditr debugIterator) Domain() (start []byte, end []byte) {
	defer func() {
		fmt.Printf("%v.itr.Domain() (%X,%X)\n", ditr.label, start, end)
	}()
	start, end = ditr.itr.Domain()
	return
}

// Implements Iterator.
func (ditr debugIterator) Valid() (ok bool) {
	defer func() {
		fmt.Printf("%v.itr.Valid() %v\n", ditr.label, ok)
	}()
	ok = ditr.itr.Valid()
	return
}

// Implements Iterator.
func (ditr debugIterator) Next() bool {
	fmt.Printf("%v.itr.Next()\n", ditr.label)
	return ditr.itr.Next()
}

func (ditr debugIterator) Seek(key []byte) bool {
	fmt.Printf("%v.itr.Seek(%v)\n", ditr.label, key)
	return ditr.itr.Seek(key)
}

// Implements Iterator.
func (ditr debugIterator) Key() (key []byte) {
	key = ditr.itr.Key()
	fmt.Printf("%v.itr.Key() %v\n", ditr.label,
		cmn.ColoredBytes(key, cmn.Cyan, cmn.Blue))
	return
}

// Implements Iterator.
func (ditr debugIterator) Value() (value []byte) {
	value = ditr.itr.Value()
	fmt.Printf("%v.itr.Value() %v\n", ditr.label,
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	return
}

// Implements Iterator.
func (ditr debugIterator) Close() {
	fmt.Printf("%v.itr.Close()\n", ditr.label)
	ditr.itr.Close()
}

//----------------------------------------
// debugBatch

type debugBatch struct {
	label string
	bch   Batch
}

// For printing all operationgs to the console for debugging.
func NewDebugBatch(label string, bch Batch) debugBatch {
	return debugBatch{
		label: label,
		bch:   bch,
	}
}

// Implements Batch.
func (dbch debugBatch) Set(key, value []byte) {
	fmt.Printf("%v.batch.Set(%v, %v)\n", dbch.label,
		cmn.ColoredBytes(key, cmn.Yellow, cmn.Blue),
		cmn.ColoredBytes(value, cmn.Green, cmn.Blue))
	dbch.bch.Set(key, value)
}

// Implements Batch.
func (dbch debugBatch) Delete(key []byte) {
	fmt.Printf("%v.batch.Delete(%v)\n", dbch.label,
		cmn.ColoredBytes(key, cmn.Red, cmn.Yellow))
	dbch.bch.Delete(key)
}

// Implements Batch.
func (dbch debugBatch) Write() {
	fmt.Printf("%v.batch.Write()\n", dbch.label)
	dbch.bch.Write()
}
func (dbch debugBatch) Commit() error {
	fmt.Printf("%v.batch.Commit()\n", dbch.label)
	return dbch.bch.Commit()
}

// Implements Batch.
func (dbch debugBatch) WriteSync() {
	fmt.Printf("%v.batch.WriteSync()\n", dbch.label)
	dbch.bch.WriteSync()
}

func (dbch debugBatch) ValueSize() int {
	fmt.Printf("%v.batch.ValueSize()\n", dbch.label)
	return dbch.bch.ValueSize()
}

func (dbch debugBatch) Reset() {
	fmt.Printf("%v.batch.Reset()\n", dbch.label)
	dbch.bch.Reset()
}
