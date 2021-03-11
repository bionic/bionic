package describer

type Import interface {
	ImportDescription() string
}

type Export interface {
	ExportDescription() string
}
