package logger

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/core/consensus/agreement"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire/topics"
)

// Make sure LogProcessor implements the TopicProcessor interface
var _ wire.TopicProcessor = (*LogProcessor)(nil)

const MonitorTopic = "monitor_topic"

const (
	ErrWriter byte = iota
	ErrLog
	ErrOther
)

// LogProcessor is a TopicProcessor that intercepts messages on the gossip to create statistics and push the to the monitoring process
// It creates a new instance of logrus and writes on a io.Writer (preferrably UNIX sockets but any kind of connection will do)
type (
	LogProcessor struct {
		*log.Logger
		lastInfo *blockInfo
		p        wire.EventPublisher
		entry    *log.Entry
	}

	blockInfo struct {
		t time.Time
		*agreement.Agreement
	}
)

func New(p wire.EventPublisher, w io.WriteCloser, formatter log.Formatter) *LogProcessor {
	logger := log.New()
	logger.Out = w
	if formatter == nil {
		logger.SetFormatter(&log.JSONFormatter{})
	}
	entry := logger.WithFields(log.Fields{
		"process": "monitor",
	})
	return &LogProcessor{
		p:      p,
		Logger: logger,
		entry:  entry,
	}
}

// Deprecated. MIght be quite useless
func (l *LogProcessor) Wire(w io.WriteCloser) {
	_ = l.Close()
	l.Out = w
}

func (l *LogProcessor) Close() error {
	return l.Out.(io.WriteCloser).Close()
}

func (l *LogProcessor) Send(entry *log.Entry) error {
	formatted, err := l.Formatter.Format(entry)
	if err != nil {
		return err
	}

	if _, err = l.Out.Write(formatted); err != nil {
		return err
	}

	return nil
}

// Process creates a copy of the message, checks the topic header
func (l *LogProcessor) Process(buf *bytes.Buffer) (*bytes.Buffer, error) {
	var newBuf bytes.Buffer
	r := io.TeeReader(buf, &newBuf)

	topic, err := topics.Extract(r)
	if err != nil {
		return nil, err
	}

	evt, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	switch topic {
	case topics.Agreement:
		go l.PublishRoundEvent(evt)
	case topics.Candidate:
		go l.PublishCandidateEvent(evt)
	}

	return &newBuf, nil
}

func (l *LogProcessor) ReportError(bErr byte, err error) {
	b := bytes.NewBuffer([]byte{bErr})
	l.p.Publish(MonitorTopic, b)
}

func (l *LogProcessor) WithTime(fields log.Fields) *log.Entry {
	entry := l.entry.WithField("time", time.Now())
	if fields == nil {
		return entry
	}
	return entry.WithFields(fields)
}

func (l *LogProcessor) WithError(err error) *log.Entry {
	return l.Logger.WithError(err).WithTime(time.Now())
}
