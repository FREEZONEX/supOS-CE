package adapter

type TimeSequenceDataStorageAdapter interface {
	DataStorageAdapter

	GetStreamHandler() StreamHandler
	ExecSQL(sql string) (string, error)
}
