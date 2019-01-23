package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/lovoo/goka/kafka"
	"github.com/tidwall/gjson"
)

type marble struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}

type schema string

const marbleSchema schema = `{
	"type":"struct",
	"fields":[
		{"type":"string", "optional": false, "field":"docType"},
		{"type":"string", "optional": false, "field":"name"},
		{"type":"string", "optional": false, "field":"color"},
		{"type":"int32",  "optional": false, "field":"size"},
		{"type":"string", "optional": false, "field":"owner"}
	]}`

type kMsg struct {
	Schema  *json.RawMessage `json:"schema,omitempty"`
	Payload *json.RawMessage `json:"payload,omitempty"`
}

// Processor ...
type Processor interface {
	Process(ctx goka.Context, msg interface{})
	Run() error
	Close()
}

type processor struct {
	brokers        []string
	group          goka.Group
	inTopic        goka.Stream
	outTopicPrefix string
	gokaProc       *goka.Processor
	stop           chan struct{}
}

func newProcessor(brokers []string, group, inTopic, outTopicPrefix string) (*processor, error) {
	var err error
	p := processor{
		brokers:        brokers,
		group:          goka.Group(group),
		inTopic:        goka.Stream(inTopic),
		outTopicPrefix: outTopicPrefix,
		stop:           make(chan struct{}),
	}

	config := kafka.NewConfig()
	config.Consumer.Offsets.Initial = kafka.OffsetOldest

	g := goka.DefineGroup(p.group,
		goka.Input(p.inTopic, new(codec.String), p.Process),
		goka.Persist(new(codec.String)),
	)
	p.gokaProc, err = goka.NewProcessor(p.brokers, g,
		goka.WithConsumerBuilder(kafka.ConsumerBuilderWithConfig(config)))
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *processor) Close() error {
	if p.stop != nil {
		p.stop <- struct{}{}
		close(p.stop)
	}
	return nil
}

func newkMsg(s schema, data interface{}) (*kMsg, error) {
	sb := json.RawMessage(s)
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	payload := json.RawMessage(b)
	msg := &kMsg{
		Schema:  &sb,
		Payload: &payload,
	}
	return msg, nil
}

func (p *processor) Process(ctx goka.Context, msg interface{}) {

	log.Println("Key: ", ctx.Key(), "\nValue: ", msg)
	data := msg.(string)
	docType := gjson.Get(data, "docType").String()

	if docType == "marble" {
		var m marble
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			ctx.Fail(err)
		}
		msg, err := newkMsg(marbleSchema, &m)
		if err != nil {
			ctx.Fail(err)
		}
		if err := p.emit("marble", m.Name, msg); err != nil {
			ctx.Fail(err)
		}
	}

	ctx.SetValue(ctx.Key())
}

func (p *processor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Gracefully stop the processor!
	rerr := make(chan error)
	go func() {
		if err := p.gokaProc.Run(ctx); err != nil {
			rerr <- err
		}
	}()
	select {
	case <-p.stop:
	case err := <-rerr:
		return err
	}
	return nil
}

func (p *processor) emit(typeName, key string, msg *kMsg) error {
	outTopic := goka.Stream(p.outTopicPrefix + "-" + typeName)
	emitter, err := goka.NewEmitter(p.brokers, outTopic, new(codec.String))
	if err != nil {
		return err
	}
	defer emitter.Finish()

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = emitter.EmitSync(key, string(b))
	if err != nil {
		return err
	}
	return nil
}
