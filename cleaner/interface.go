package cleaner

type Cleaner interface {
	UpdateTTL() error
}
