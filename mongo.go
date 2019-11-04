package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	// Mongo mongodb
	Mongo struct {
		URI string `config:"uri"`

		client *mongo.Client
		name   string
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

// ConfigWillLoad before load
func (m *Mongo) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad after load
func (m *Mongo) ConfigDidLoad(context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.URI))
	if err != nil {
		panic(fmt.Sprintf("NewClient mongodb [%s] error: %#v", m.URI, err))
	}

	m.client = client
}

// Serve Connect and ping to mongodb
func (m *Mongo) Serve(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Connect(ctx); err != nil {
		return fmt.Errorf("Connect to mongodb [%s] error: %#v", m.URI, err)
	}

	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("Ping to mongodb [%s] error: %#v", m.URI, err)
	}

	logger.Infof("Connect to mongodb [%s] success", m.URI)

	return nil
}

// Shutdown Disconnect to mongodb
func (m *Mongo) Shutdown(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		panic(fmt.Sprintf("Disconnect to mongodb [%s] error: %#v", m.URI, err))
	}

	logger.Infof("Disconnect to mongodb [%s] success", m.URI)

	return nil
}

// Client Get mongodb client
func (m *Mongo) Client() *mongo.Client {
	return m.client
}

// New mongodb instance
func New(name string) *Mongo {
	mg := &Mongo{
		name: name,
	}

	return mg
}
