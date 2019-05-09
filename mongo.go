package mongo

import (
	"context"

	"github.com/boxgo/box/minibox"
	"github.com/boxgo/metrics"
	"github.com/globalsign/mgo"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Mongo mongodb数据库
	Mongo struct {
		Metrics   bool    `config:"metrics" desc:"default is false"`
		URI       string  `json:"uri"`
		DB        string  `json:"db"`
		PoolLimit uint    `json:"poolLimit"`
		Batch     uint    `json:"batch"`
		Prefetch  float64 `json:"prefetch"`
		Mode      uint    `json:"mode"`

		name    string
		session *mgo.Session
		metrics *metrics.Metrics
	}
)

var (
	// Default mongodb
	Default = New("mongo")
)

// Name mongodb config name
func (m *Mongo) Name() string {
	return m.name
}

// Exts app
func (m *Mongo) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{m.metrics}
}

// ConfigWillLoad before load
func (m *Mongo) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad after load
func (m *Mongo) ConfigDidLoad(context.Context) {
	if m.PoolLimit == 0 {
		m.PoolLimit = 200
	}

	if m.Batch == 0 {
		m.Batch = 50
	}

	if m.Prefetch <= 0 {
		m.Prefetch = 0.20
	}

	if m.Metrics {
		prometheus.MustRegister(NewMgoCollector(m.metrics.Namespace, m.metrics.Subsystem))
	}
}

// Serve start
func (m *Mongo) Serve(ctx context.Context) error {
	m.GetSession()

	return nil
}

// Shutdown end
func (m *Mongo) Shutdown(ctx context.Context) error {
	if m.session != nil {
		m.session.Close()
	}

	return nil
}

// GetSession get session
func (m *Mongo) GetSession() *mgo.Session {
	if m.session != nil {
		return m.session
	}

	sess, err := mgo.Dial(m.URI)
	if err != nil {
		panic(err)
	}

	m.session = sess
	m.session.SetMode(mgo.Mode(m.Mode), true)
	m.session.SetPoolLimit(200)
	m.session.SetBatch(50)
	m.session.SetPrefetch(0.20)

	return m.session
}

// GetDB get selected db
func (m *Mongo) GetDB(db string) *mgo.Database {
	return m.GetSession().DB(db)
}

// GetDefaultDB get default db
func (m *Mongo) GetDefaultDB() *mgo.Database {
	return m.GetSession().DB(m.DB)
}

// New mongodb
func New(name string, ms ...*metrics.Metrics) *Mongo {
	mg := &Mongo{
		name: name,
	}

	if len(ms) == 0 {
		mg.metrics = metrics.Default
	} else {
		mg.metrics = ms[0]
	}

	return mg
}
