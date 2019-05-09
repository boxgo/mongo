package mongo

import (
	"strings"

	"github.com/globalsign/mgo"
	"github.com/prometheus/client_golang/prometheus"
)

type mgoCollector struct {
	clusterDesc      *prometheus.Desc
	masterConnDesc   *prometheus.Desc
	slaveConnDesc    *prometheus.Desc
	sentOpDesc       *prometheus.Desc
	receivedOpDesc   *prometheus.Desc
	receivedDocDesc  *prometheus.Desc
	socketsAliveDesc *prometheus.Desc
	socketsInUseDesc *prometheus.Desc
	socketRefDesc    *prometheus.Desc
}

// NewMgoCollector returns a collector which exports metrics about mgo stats.
func NewMgoCollector(namespace, subsystem string) prometheus.Collector {
	mgo.SetStats(true)

	prefixArr := []string{}
	if namespace != "" {
		prefixArr = append(prefixArr, namespace)
	}
	if subsystem != "" {
		prefixArr = append(prefixArr, subsystem)
	}

	prefix := strings.Join(prefixArr, "_")
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	return &mgoCollector{
		clusterDesc: prometheus.NewDesc(
			prefix+"mgo_cluster",
			"Mgo Cluster from stats",
			nil, nil,
		),
		masterConnDesc: prometheus.NewDesc(
			prefix+"mgo_master_conn",
			"Mgo Master Conn from stats",
			nil, nil,
		),
		slaveConnDesc: prometheus.NewDesc(
			prefix+"mgo_slave_conn",
			"Mgo Slave Conn from stats",
			nil, nil,
		),
		sentOpDesc: prometheus.NewDesc(
			prefix+"mgo_sent_op",
			"Mgo sent op from stats",
			nil, nil,
		),
		receivedOpDesc: prometheus.NewDesc(
			prefix+"mgo_received_op",
			"Mgo received op from stats",
			nil, nil,
		),
		receivedDocDesc: prometheus.NewDesc(
			prefix+"mgo_received_doc",
			"Mgo received doc from stats",
			nil, nil,
		),
		socketsAliveDesc: prometheus.NewDesc(
			prefix+"mgo_sockets_alive",
			"Mgo sockets alive from stats",
			nil, nil,
		),
		socketsInUseDesc: prometheus.NewDesc(
			prefix+"mgo_sockets_in_use",
			"Mgo sockets in use from stats",
			nil, nil,
		),
		socketRefDesc: prometheus.NewDesc(
			prefix+"mgo_socket_ref",
			"Mgo socket ref from stats",
			nil, nil,
		),
	}
}

// Describe returns all descriptions of the collector.
func (c *mgoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.clusterDesc
	ch <- c.masterConnDesc
	ch <- c.slaveConnDesc
	ch <- c.sentOpDesc
	ch <- c.receivedOpDesc
	ch <- c.receivedDocDesc
	ch <- c.socketsAliveDesc
	ch <- c.socketsInUseDesc
	ch <- c.socketRefDesc
}

// Collect returns the current state of all metrics of the collector.
func (c *mgoCollector) Collect(ch chan<- prometheus.Metric) {
	stats := mgo.GetStats()

	ch <- prometheus.MustNewConstMetric(c.clusterDesc, prometheus.GaugeValue, float64(stats.Clusters))
	ch <- prometheus.MustNewConstMetric(c.masterConnDesc, prometheus.GaugeValue, float64(stats.MasterConns))
	ch <- prometheus.MustNewConstMetric(c.slaveConnDesc, prometheus.GaugeValue, float64(stats.SlaveConns))
	ch <- prometheus.MustNewConstMetric(c.sentOpDesc, prometheus.GaugeValue, float64(stats.SentOps))
	ch <- prometheus.MustNewConstMetric(c.receivedOpDesc, prometheus.GaugeValue, float64(stats.ReceivedOps))
	ch <- prometheus.MustNewConstMetric(c.receivedDocDesc, prometheus.GaugeValue, float64(stats.ReceivedDocs))
	ch <- prometheus.MustNewConstMetric(c.socketsAliveDesc, prometheus.GaugeValue, float64(stats.SocketsAlive))
	ch <- prometheus.MustNewConstMetric(c.socketsInUseDesc, prometheus.GaugeValue, float64(stats.SocketsInUse))
	ch <- prometheus.MustNewConstMetric(c.socketRefDesc, prometheus.GaugeValue, float64(stats.SocketRefs))
}
