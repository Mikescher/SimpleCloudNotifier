package models

type TransactionLockMode string //@enum:type

const (
	TLockNone      TransactionLockMode = "NONE"
	TLockRead      TransactionLockMode = "READ"
	TLockReadWrite TransactionLockMode = "READ_WRITE"
)
