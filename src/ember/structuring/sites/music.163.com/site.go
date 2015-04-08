package structuring

func (p *Song) Run(appender Appender) (err error) {
	// TODO
	// extrat infos
	// get similars
	// new song task -> chan
	return
}

type Song struct {
	url string
}

func (p *Site) NewTask(url string, typ string) Task {
	switch typ {
	case "song":
		return &Song{url}
	}
	return nil
}

func New() *Site {
	return &Site{}
}

type Site struct {
}

type Task interface {
	Run(appender Appender) error
}
type Appender func(url string, typ string, weight int) error
