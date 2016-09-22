package resources

type Transaction interface {
	GetValue(bucketName, key []byte) []byte
	SetValue(bucketName, key, value []byte) error
	ModifyValue(bucketName, key []byte, modify func ([]byte) []byte) error
	EnsureValue(bucketName, key, value []byte) error
	EnsureEntity(bucketName, key []byte, entity Entity) error
	AddEntity(bucketName []byte, entity Entity) error
	DeleteEntity(bucketName, id []byte) error
	RetrieveEntity(bucketName, id []byte, entity Entity) error
	RetrieveEntities(bucketName []byte, getObject func (string) Entity) error
	ModifyEntity(bucketName, key []byte, entity Entity, modifyFunc func ()) error
	MapEntities(bucketName []byte, entity Entity, mapFunc func ()) error
	InitializeBucket(bucketName []byte) error
	FilterEntities(bucketName []byte, addEntity func (), getNewEntity func () Entity, filterFunc func () bool) error
	Execute()
	View()
	Add(exec func () error)
}

type AnybarManager interface {
	RemoveAndQuit(id string, t Transaction)
	AddToActivePorts(title, icon string, id string, t Transaction)
	EnsureActivePorts(activePorts ActivePorts)
	StartWithIcon(port int, title, icon string)
	StartNew(port int, title string)
	ChangeIcon(port int, colour string)
	GetNewPort(activePorts []*ActivePort) int
	Quit(port int)
	GetActivePorts(t Transaction) ActivePorts
	Ping(port int) bool
}

type Entity interface {
	SetId(string)
	GetId() string
	Load(Transaction) error
}

type Item interface {
	GetItem() *AlfredItem
}