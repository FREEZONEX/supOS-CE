package dataingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"backend/internal/config"
	"backend/internal/infra/cache"
	"backend/internal/infra/eventstream"
	"backend/internal/repo"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

type Service struct {
	cfg             config.DataIngestConf
	resolver        *DefinitionResolver
	writer          *SinkWriter
	latest          *latestCache
	queue           chan []Record
	stop            chan struct{}
	done            chan struct{}
	client          mqtt.Client
	publishOverride func(topic string, qos byte, retain bool, payload []byte) error
	started         atomic.Bool
	written         atomic.Int64
	dropped         atomic.Int64
	mu              sync.Mutex
}

type Stats struct {
	Written int64 `json:"written"`
	Dropped int64 `json:"dropped"`
	Queued  int   `json:"queued"`
}

const (
	defaultHistoryLimit  = 1000
	maxHistoryLimit      = 2000
	defaultHistoryWindow = 5 * time.Minute
)

func New(ctx context.Context, cfg config.DataIngestConf, sinkDBURL string, cacheClients ...*cache.Client) (*Service, error) {
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 20000
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 5000
	}
	if cfg.FlushIntervalMs <= 0 {
		cfg.FlushIntervalMs = 1000
	}
	var cacheClient *cache.Client
	if len(cacheClients) > 0 {
		cacheClient = cacheClients[0]
	}
	s := &Service{
		cfg:      cfg,
		resolver: NewDefinitionResolver(repo.NewUnsRepo(ctx)),
		latest:   newLatestCache(cacheClient),
		queue:    make(chan []Record, cfg.QueueSize),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	if strings.TrimSpace(sinkDBURL) != "" {
		writer, err := NewSinkWriter(ctx, sinkDBURL, cfg.BatchSize)
		if err != nil {
			return nil, err
		}
		s.writer = writer
	}
	return s, nil
}

func (s *Service) Resolver() *DefinitionResolver {
	if s == nil {
		return nil
	}
	return s.resolver
}

func (s *Service) EnsurePhysicalForNode(ctx context.Context, node repo.UnsNode) error {
	return s.EnsurePhysicalForNodes(ctx, []repo.UnsNode{node})
}

func (s *Service) EnsurePhysicalForNodes(ctx context.Context, nodes []repo.UnsNode) error {
	if s == nil || s.writer == nil || s.writer.pool == nil || len(nodes) == 0 {
		return nil
	}
	// dataingest 关闭时同时跳过 UNS 节点触发的物理表结构维护（ALTER 加列）。
	// 否则 outbox 的 uns.physical.ensure 任务仍会在节点创建/变更时持续触发
	// uns_timeserial 动态加列，绕过 DATAINGEST_ENABLED 开关，在小内存主机上引发 OOM。
	if !s.cfg.Enabled {
		return nil
	}
	records := make([]Record, 0, len(nodes))
	for _, node := range nodes {
		if node.Type != 2 {
			continue
		}
		def := definitionFromNode(node)
		if !def.SaveToDB {
			continue
		}
		fields := dataFields(def.Fields, def.TopicType)
		if len(fields) == 0 {
			fields = []Field{{Name: "value", Type: "string"}}
		}
		records = append(records, Record{
			UnsID:     def.ID,
			Alias:     def.Alias,
			Namespace: def.Namespace,
			TopicType: def.TopicType,
			Fields:    fields,
			Time:      time.Now(),
		})
	}
	if len(records) == 0 {
		return nil
	}
	_, err := s.writer.ensurePhysicalForRecords(ctx, records)
	return err
}

func (s *Service) Start(ctx context.Context) error {
	if s == nil {
		return nil
	}
	if !s.started.CompareAndSwap(false, true) {
		return nil
	}
	go s.flushLoop()
	if !s.cfg.Enabled {
		return nil
	}
	client, err := s.connectMQTT(ctx)
	if err != nil {
		s.started.Store(false)
		close(s.stop)
		<-s.done
		return err
	}
	s.client = client
	return nil
}

func (s *Service) Close() {
	if s == nil {
		return
	}
	if s.started.CompareAndSwap(true, false) {
		close(s.stop)
		if s.client != nil {
			s.client.Disconnect(500)
		}
		<-s.done
	}
	if s.writer != nil {
		s.writer.Close()
	}
}

func (s *Service) Stats() Stats {
	if s == nil {
		return Stats{}
	}
	return Stats{Written: s.written.Load(), Dropped: s.dropped.Load(), Queued: len(s.queue)}
}

func (s *Service) LastPayload(ctx context.Context, node repo.UnsNode) (map[string]any, error) {
	if s == nil || node.Type != 2 {
		return nil, nil
	}
	if s.latest != nil {
		payload, ok, err := s.latest.Get(ctx, node.ID)
		if err == nil && ok {
			return payload, nil
		}
		if err != nil {
			logx.WithContext(ctx).Errorf("dataingest latest cache read failed node=%d err=%v", node.ID, err)
		}
	}
	payload, err := s.lastPayloadFromDB(ctx, node)
	if err != nil || len(payload) == 0 {
		return payload, err
	}
	if s.latest != nil {
		if err := s.latest.UpsertPayload(ctx, node.ID, payload); err != nil {
			logx.WithContext(ctx).Errorf("dataingest latest cache backfill failed node=%d err=%v", node.ID, err)
		}
	}
	return payload, nil
}

func (s *Service) lastPayloadFromDB(ctx context.Context, node repo.UnsNode) (map[string]any, error) {
	if s == nil || s.writer == nil || s.writer.pool == nil || node.Type != 2 {
		return nil, nil
	}
	def := definitionFromNode(node)
	fields := dataFields(def.Fields, def.TopicType)
	if len(fields) == 0 {
		fields = []Field{{Name: "value", Type: "string"}}
	}
	alias := strings.Trim(def.Alias, "/")
	if alias == "" {
		alias = strings.Trim(def.Namespace, "/")
	}
	if alias == "" {
		return nil, nil
	}
	timeColumn := sinkTimeColumn(def)
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY %s DESC LIMIT 1`, sinkQualifiedName(alias), quoteIdent(timeColumn))
	rows, err := s.writer.pool.Query(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	values, err := rows.Values()
	if err != nil {
		return nil, err
	}
	desc := rows.FieldDescriptions()
	row := make(map[string]any, len(desc))
	for i, field := range desc {
		row[string(field.Name)] = values[i]
	}
	latest := valueTimeFromAny(row[timeColumn])
	data := payloadFromSinkRow(row, fields, def.TopicType)
	if len(data) == 0 {
		return nil, nil
	}
	dt := make(map[string]int64, len(data))
	keys := make([]string, 0, len(data))
	for key := range data {
		dt[key] = latest.UnixMilli()
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return map[string]any{
		"data":       data,
		"payload":    data,
		"dt":         dt,
		"timeStamp":  latest.UnixMilli(),
		"updateTime": latest.UnixMilli(),
		"fields":     keys,
		"quality":    qualityFromSinkRow(row, def.TopicType),
	}, nil
}

func (s *Service) ClearLatest(ctx context.Context, nodes []repo.UnsNode) error {
	if s == nil || s.latest == nil || len(nodes) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(nodes))
	for _, node := range nodes {
		if node.Type == 2 {
			ids = append(ids, node.ID)
		}
	}
	return s.latest.DeleteIDs(ctx, ids...)
}

func (s *Service) History(ctx context.Context, node repo.UnsNode, startMs, endMs int64, limit int) (map[string]any, error) {
	return s.QueryHistory(ctx, node, HistoryQuery{
		StartMs:    startMs,
		EndMs:      endMs,
		Auto:       true,
		PointLimit: limit,
	})
}

func (s *Service) WriteOpenAPIValue(
	ctx context.Context,
	node repo.UnsNode,
	value any,
	timestampMs int64,
	qos byte,
	retain bool,
) error {
	if s == nil || node.Type != 2 {
		return nil
	}
	def := definitionFromNode(node)
	def.SaveToDB = true
	payloadValue := value
	if timestampMs > 0 {
		if body, ok := value.(map[string]any); ok {
			withTime := make(map[string]any, len(body)+1)
			for key, item := range body {
				withTime[key] = item
			}
			if _, ok := withTime["timeStamp"]; !ok {
				withTime["timeStamp"] = timestampMs
			}
			payloadValue = withTime
		} else {
			payloadValue = map[string]any{
				"value":     value,
				"timeStamp": timestampMs,
			}
		}
	}
	payload, err := json.Marshal(payloadValue)
	if err != nil {
		return err
	}
	records, err := NormalizePayload(ctx, def, payload)
	if err != nil {
		return err
	}
	if err := s.writeLatestCache(ctx, records); err != nil {
		logx.WithContext(ctx).Errorf("dataingest openapi latest cache write failed node=%d err=%v", node.ID, err)
	}
	publishPayloadRecords(records)
	if s.writer != nil {
		if err := s.writer.WriteRecords(ctx, records); err != nil {
			return err
		}
	}
	// Publish the same logical value to MQTT so the OPC UA server (and other
	// MQTT consumers) see the update. For a plain string without an explicit
	// timestamp we publish the raw string to match typical UNS MQTT payload
	// semantics; otherwise we publish the JSON representation.
	mqttPayload := payload
	if timestampMs <= 0 {
		if str, ok := value.(string); ok {
			mqttPayload = []byte(str)
		}
	}
	return s.publishOpenAPIWrite(node.Namespace, qos, retain, mqttPayload)
}

func (s *Service) publishOpenAPIWrite(
	topic string,
	qos byte,
	retain bool,
	payload []byte,
) error {
	topic = strings.Trim(topic, "/")
	if topic == "" {
		return errors.New("mqtt topic is empty")
	}
	if s.publishOverride != nil {
		return s.publishOverride(topic, qos, retain, payload)
	}
	if s.client == nil {
		return errors.New("mqtt client is not initialized")
	}
	if !s.client.IsConnectionOpen() {
		return errors.New("mqtt broker is not connected")
	}
	token := s.client.Publish(topic, qos, retain, payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("mqtt publish timeout for topic %s", topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish failed for topic %s: %w", topic, err)
	}
	return nil
}

func (s *Service) Consume(ctx context.Context, topic string, payload []byte) error {
	if s == nil || s.resolver == nil || strings.HasPrefix(topic, "$") {
		return nil
	}
	def, err := s.resolver.Resolve(ctx, topic)
	if err != nil {
		if errors.Is(err, ErrDefinitionNotFound) {
			return nil
		}
		return err
	}
	records, err := NormalizePayload(ctx, def, payload)
	if err != nil || len(records) == 0 {
		return err
	}
	if err := s.writeLatestCache(ctx, records); err != nil {
		logx.WithContext(ctx).Errorf("dataingest latest cache write failed topic=%s err=%v", topic, err)
	}
	publishPayloadRecords(records)
	if !def.SaveToDB || s.writer == nil {
		return nil
	}
	select {
	case s.queue <- records:
	default:
		s.dropped.Add(int64(len(records)))
		return errors.New("dataingest queue full")
	}
	return nil
}

func (s *Service) writeLatestCache(ctx context.Context, records []Record) error {
	if s == nil || s.latest == nil || len(records) == 0 {
		return nil
	}
	return s.latest.UpsertRecords(ctx, records)
}

func publishPayloadRecords(records []Record) {
	for _, record := range records {
		data := livePayloadFromRecord(record)
		if len(data) == 0 {
			continue
		}
		keys := []string{fmt.Sprintf("id:%d", record.UnsID)}
		if alias := strings.Trim(record.Alias, "/"); alias != "" {
			keys = append(keys, "alias:"+alias)
		}
		if namespace := strings.Trim(record.Namespace, "/"); namespace != "" {
			keys = append(keys, "namespace:"+namespace)
		}
		eventstream.Publish(eventstream.Event{
			Type: "uns.payload",
			Keys: keys,
			Data: data,
		})
	}
}

func livePayloadFromRecord(record Record) map[string]any {
	if record.TopicType != metricTopicType {
		return liveEventPayloadFromRecord(record)
	}
	fields := recordFields(record)
	if len(fields) == 0 {
		return nil
	}
	ts := record.Time.UTC().UnixMilli()
	data := make(map[string]any, len(fields))
	dt := make(map[string]int64, len(fields))
	names := make([]string, 0, len(fields))
	for _, field := range fields {
		raw, ok := record.Values[field.Name]
		if !ok {
			continue
		}
		value, err := typedValueByType(field.Type, raw)
		if err != nil {
			value = raw
		}
		data[field.Name] = value
		dt[field.Name] = ts
		names = append(names, field.Name)
	}
	if len(data) == 0 {
		return nil
	}
	sort.Strings(names)
	return map[string]any{
		"data":       data,
		"payload":    data,
		"dt":         dt,
		"updateTime": ts,
		"fields":     names,
	}
}

func liveEventPayloadFromRecord(record Record) map[string]any {
	data := eventPayloadValues(record)
	if len(data) == 0 {
		return nil
	}
	ts := record.Time.UTC().UnixMilli()
	dt := make(map[string]int64, len(data))
	names := make([]string, 0, len(data))
	for name := range data {
		dt[name] = ts
		names = append(names, name)
	}
	sort.Strings(names)
	return map[string]any{
		"data":       data,
		"payload":    data,
		"dt":         dt,
		"updateTime": ts,
		"fields":     names,
	}
}

func (s *Service) connectMQTT(ctx context.Context) (mqtt.Client, error) {
	topics := splitTopics(nonEmpty(s.cfg.MqttTopic, "#"))
	if len(topics) == 0 {
		topics = []string{"#"}
	}
	// 参考 share/clients/mqttClient：每次尝试新建 client，失败递增退避重试 5 次，仍失败则返回错误
	const tryTimes = 5
	var (
		client mqtt.Client
		err    error
	)
	for i := tryTimes; i > 0; i-- {
		client, err = s.newMQTTClient(topics)
		if err == nil {
			break
		}
		logx.Errorf("dataingest mqtt connect failed, retries left: %d, err: %v", i-1, err)
		if i > 1 {
			time.Sleep(time.Second * time.Duration(tryTimes) / time.Duration(i))
		}
	}
	if err != nil {
		return nil, err
	}
	if err := s.subscribeMQTTTopics(ctx, client, topics); err != nil {
		client.Disconnect(0)
		return nil, err
	}
	return client, nil
}

func (s *Service) newMQTTClient(topics []string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range splitBrokers(s.cfg.MqttBrokers) {
		opts.AddBroker(broker)
	}
	if len(opts.Servers) == 0 {
		return nil, errors.New("dataingest mqtt brokers is empty")
	}
	opts.SetClientID(nonEmpty(s.cfg.MqttClientID, "edge-backend-sink"))
	opts.SetUsername(s.cfg.MqttUsername)
	opts.SetPassword(s.cfg.MqttPassword)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetCleanSession(false)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	// 与下方 WaitTimeout(5s) 的单次尝试预算对齐，避免超时后底层 dial 仍在进行，
	// 旧尝试稍后连上时与同 clientID 的新 client 互相挤占
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		go s.subscribeMQTTTopicsWithRetry(context.Background(), client, topics)
	})
	opts.SetReconnectingHandler(func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		logx.Infof("dataingest mqtt reconnecting")
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logx.Errorf("dataingest mqtt connection lost: %v", err)
	})
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(5 * time.Second) {
		client.Disconnect(0)
		return nil, errors.New("connect mqtt timeout")
	}
	if err := token.Error(); err != nil {
		client.Disconnect(0)
		return nil, err
	}
	return client, nil
}

func (s *Service) subscribeMQTTTopics(ctx context.Context, client mqtt.Client, topics []string) error {
	for _, topic := range topics {
		token := client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
			_ = s.Consume(context.Background(), msg.Topic(), msg.Payload())
		})
		if !token.WaitTimeout(10 * time.Second) {
			return fmt.Errorf("subscribe mqtt %s timeout", topic)
		}
		if err := token.Error(); err != nil {
			return fmt.Errorf("subscribe mqtt %s: %w", topic, err)
		}
	}
	_ = ctx
	return nil
}

func (s *Service) subscribeMQTTTopicsWithRetry(ctx context.Context, client mqtt.Client, topics []string) {
	for {
		if err := s.subscribeMQTTTopics(ctx, client, topics); err == nil {
			return
		} else {
			logx.Errorf("dataingest mqtt subscribe failed, retrying: %v", err)
		}
		select {
		case <-s.stop:
			return
		case <-time.After(3 * time.Second):
		}
	}
}

func (s *Service) flushLoop() {
	defer close(s.done)
	ticker := time.NewTicker(time.Duration(s.cfg.FlushIntervalMs) * time.Millisecond)
	defer ticker.Stop()
	batch := make([]Record, 0, s.cfg.BatchSize)
	flush := func() {
		if len(batch) == 0 || s.writer == nil {
			return
		}
		records := batch
		batch = make([]Record, 0, s.cfg.BatchSize)
		if err := s.writer.WriteRecords(context.Background(), records); err != nil {
			s.dropped.Add(int64(len(records)))
			return
		}
		s.written.Add(int64(len(records)))
	}
	for {
		select {
		case <-s.stop:
			for {
				select {
				case records := <-s.queue:
					batch = append(batch, records...)
				default:
					flush()
					return
				}
			}
		case <-ticker.C:
			flush()
		case records := <-s.queue:
			batch = append(batch, records...)
			if len(batch) >= s.cfg.BatchSize {
				flush()
			}
		}
	}
}

func historyFields(fields []Field) []map[string]any {
	out := make([]map[string]any, 0, len(fields))
	for _, field := range fields {
		out = append(out, map[string]any{
			"name": field.Name,
			"type": field.Type,
		})
	}
	return out
}

func sinkTimeColumn(def *Definition) string {
	return metricTimeColumn
}

func payloadFromSinkRow(row map[string]any, fields []Field, topicType int16) map[string]any {
	if topicType != metricTopicType {
		return eventPayloadFromSinkJSON(row[eventJSONColumn])
	}
	payload := make(map[string]any, len(fields))
	for _, field := range fields {
		if value, ok := row[field.Name]; ok {
			payload[field.Name] = value
		}
	}
	return payload
}

func qualityFromSinkRow(row map[string]any, topicType int16) string {
	if topicType == metricTopicType {
		return qualityFromAny(row[metricQualityColumn])
	}
	return qualityFromEventJSON(row[eventJSONColumn])
}

func qualityFromEventJSON(value any) string {
	raw := eventPayloadFromSinkJSON(value)
	for _, key := range []string{metricQualityColumn, "quality", "_qos"} {
		if value, ok := raw[key]; ok {
			return qualityFromAny(value)
		}
	}
	return qualityGood
}

func eventPayloadFromSinkJSON(value any) map[string]any {
	var decoded any
	switch v := value.(type) {
	case nil:
		return map[string]any{}
	case []byte:
		decoded, _ = decodeJSONValue(v)
	case string:
		decoded, _ = decodeJSONValue([]byte(v))
	case json.RawMessage:
		decoded, _ = decodeJSONValue(v)
	case map[string]any:
		return clonePayloadMap(v)
	default:
		data, err := json.Marshal(v)
		if err == nil {
			decoded, _ = decodeJSONValue(data)
		}
	}
	raw, _ := decoded.(map[string]any)
	return clonePayloadMap(raw)
}

func historyWindow(startMs, endMs int64) (time.Time, time.Time) {
	endAt := time.Now().UTC()
	if endMs > 0 {
		endAt = time.UnixMilli(endMs).UTC()
	}
	startAt := endAt.Add(-defaultHistoryWindow)
	if startMs > 0 {
		startAt = time.UnixMilli(startMs).UTC()
	}
	if startAt.After(endAt) {
		startAt = endAt.Add(-defaultHistoryWindow)
	}
	return startAt, endAt
}

func historyLimit(limit int) int {
	if limit <= 0 {
		return defaultHistoryLimit
	}
	if limit > maxHistoryLimit {
		return maxHistoryLimit
	}
	return limit
}

func valueTimeFromAny(value any) time.Time {
	switch v := value.(type) {
	case time.Time:
		return v.UTC()
	case string:
		if t, ok := parseTime(v); ok {
			return t
		}
	}
	return time.Now().UTC()
}

func splitBrokers(raw string) []string {
	return splitList(raw)
}

func splitTopics(raw string) []string {
	return splitList(raw)
}

func splitList(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ';' })
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func nonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
