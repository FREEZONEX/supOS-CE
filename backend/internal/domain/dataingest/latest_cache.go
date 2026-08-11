package dataingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"backend/internal/infra/cache"
)

const (
	latestCacheKeyPrefix = "cache:uns:latest"

	qualityGood      = "Good"
	qualityUncertain = "Uncertain"
	qualityBad       = "Bad"
)

const latestCacheUpsertScript = `
local current = redis.call("GET", KEYS[1])
local incoming = tonumber(ARGV[2]) or 0
if current and current ~= false then
  local ok, currentObj = pcall(cjson.decode, current)
  if ok and currentObj then
    local currentTs = tonumber(currentObj.timeStamp or currentObj.updateTime) or 0
    if currentTs > incoming then
      return 0
    end
  end
end
redis.call("SET", KEYS[1], ARGV[1])
return 1
`

type latestCache struct {
	redis *cache.Client
}

func newLatestCache(redis *cache.Client) *latestCache {
	if redis == nil {
		return nil
	}
	return &latestCache{redis: redis}
}

func (c *latestCache) Get(ctx context.Context, nodeID int64) (map[string]any, bool, error) {
	if c == nil || c.redis == nil || nodeID <= 0 {
		return nil, false, nil
	}
	raw, err := c.redis.GetString(ctx, latestCacheKey(nodeID))
	if errors.Is(err, cache.ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, false, err
	}
	normalizeLatestPayload(payload)
	if len(payload) == 0 {
		return nil, false, nil
	}
	return payload, true, nil
}

func (c *latestCache) UpsertRecords(ctx context.Context, records []Record) error {
	if c == nil || c.redis == nil || len(records) == 0 {
		return nil
	}
	latestByID := make(map[int64]map[string]any, len(records))
	for _, record := range records {
		if record.UnsID <= 0 {
			continue
		}
		payload := latestPayloadFromRecord(record)
		if len(payload) == 0 {
			continue
		}
		current := latestByID[record.UnsID]
		if current == nil || latestPayloadTimestamp(payload) >= latestPayloadTimestamp(current) {
			latestByID[record.UnsID] = payload
		}
	}
	for nodeID, payload := range latestByID {
		if err := c.UpsertPayload(ctx, nodeID, payload); err != nil {
			return err
		}
	}
	return nil
}

func (c *latestCache) UpsertPayload(ctx context.Context, nodeID int64, payload map[string]any) error {
	if c == nil || c.redis == nil || nodeID <= 0 || len(payload) == 0 {
		return nil
	}
	normalizeLatestPayload(payload)
	ts := latestPayloadTimestamp(payload)
	if ts <= 0 {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.redis.EvalInt64(ctx, latestCacheUpsertScript, []string{latestCacheKey(nodeID)}, string(body), strconv.FormatInt(ts, 10))
	return err
}

func (c *latestCache) DeleteIDs(ctx context.Context, ids ...int64) error {
	if c == nil || c.redis == nil || len(ids) == 0 {
		return nil
	}
	keys := make([]string, 0, len(ids))
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		keys = append(keys, latestCacheKey(id))
	}
	return c.redis.Del(ctx, keys...)
}

func latestCacheKey(nodeID int64) string {
	return fmt.Sprintf("%s:%d", latestCacheKeyPrefix, nodeID)
}

func latestPayloadFromRecord(record Record) map[string]any {
	payload := livePayloadFromRecord(record)
	if len(payload) == 0 {
		return nil
	}
	ts := latestPayloadTimestamp(payload)
	payload["timeStamp"] = ts
	payload["quality"] = qualityFromRecord(record)
	return payload
}

func normalizeLatestPayload(payload map[string]any) {
	if payload == nil {
		return
	}
	data, dataOK := payload["data"]
	body, bodyOK := payload["payload"]
	if !dataOK && bodyOK {
		payload["data"] = body
		data = body
		dataOK = true
	}
	if !bodyOK && dataOK {
		payload["payload"] = data
	}
	ts := latestPayloadTimestamp(payload)
	if ts > 0 {
		payload["timeStamp"] = ts
		if openapiAny := payload["updateTime"]; openapiAny == nil || anyInt64(openapiAny) == 0 {
			payload["updateTime"] = ts
		}
	}
	fields := latestPayloadFields(payload)
	if len(fields) > 0 {
		payload["fields"] = fields
		if _, ok := latestPayloadDT(payload); !ok && ts > 0 {
			dt := make(map[string]int64, len(fields))
			for _, field := range fields {
				dt[field] = ts
			}
			payload["dt"] = dt
		}
	}
	payload["quality"] = qualityFromAny(payload["quality"])
}

func latestPayloadTimestamp(payload map[string]any) int64 {
	if payload == nil {
		return 0
	}
	for _, key := range []string{"timeStamp", "updateTime", "timestamp"} {
		if ts := anyInt64(payload[key]); ts > 0 {
			return ts
		}
	}
	if dt, ok := latestPayloadDT(payload); ok {
		var latest int64
		for _, ts := range dt {
			if ts > latest {
				latest = ts
			}
		}
		return latest
	}
	return 0
}

func latestPayloadFields(payload map[string]any) []string {
	if payload == nil {
		return nil
	}
	switch raw := payload["fields"].(type) {
	case []string:
		out := append([]string(nil), raw...)
		sort.Strings(out)
		return out
	case []any:
		out := make([]string, 0, len(raw))
		for _, item := range raw {
			if field := strings.TrimSpace(fmt.Sprint(item)); field != "" {
				out = append(out, field)
			}
		}
		sort.Strings(out)
		return out
	}
	body, _ := payload["data"].(map[string]any)
	if len(body) == 0 {
		body, _ = payload["payload"].(map[string]any)
	}
	if len(body) == 0 {
		return nil
	}
	out := make([]string, 0, len(body))
	for field := range body {
		out = append(out, field)
	}
	sort.Strings(out)
	return out
}

func latestPayloadDT(payload map[string]any) (map[string]int64, bool) {
	if payload == nil {
		return nil, false
	}
	switch raw := payload["dt"].(type) {
	case map[string]int64:
		return raw, true
	case map[string]any:
		out := make(map[string]int64, len(raw))
		for key, value := range raw {
			out[key] = anyInt64(value)
		}
		return out, len(out) > 0
	}
	return nil, false
}

func qualityFromRecord(record Record) string {
	for _, key := range []string{metricQualityColumn, "quality", "_qos"} {
		if raw, ok := record.RawValues[key]; ok {
			return qualityFromAny(raw)
		}
	}
	return qualityFromAny(record.Quality)
}

func qualityFromAny(value any) string {
	switch v := value.(type) {
	case nil:
		return qualityGood
	case string:
		return normalizeQuality(v)
	case []byte:
		return normalizeQuality(string(v))
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return qualityFromNumber(n)
		}
		return normalizeQuality(v.String())
	case int:
		return qualityFromNumber(int64(v))
	case int8:
		return qualityFromNumber(int64(v))
	case int16:
		return qualityFromNumber(int64(v))
	case int32:
		return qualityFromNumber(int64(v))
	case int64:
		return qualityFromNumber(v)
	case uint:
		return qualityFromNumber(int64(v))
	case uint8:
		return qualityFromNumber(int64(v))
	case uint16:
		return qualityFromNumber(int64(v))
	case uint32:
		return qualityFromNumber(int64(v))
	case uint64:
		return qualityFromNumber(int64(v))
	case float32:
		return qualityFromNumber(int64(v))
	case float64:
		return qualityFromNumber(int64(v))
	default:
		return normalizeQuality(fmt.Sprint(v))
	}
}

func normalizeQuality(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "0", "1", strings.ToLower(qualityGood):
		return qualityGood
	case "2", strings.ToLower(qualityUncertain):
		return qualityUncertain
	case "3", strings.ToLower(qualityBad):
		return qualityBad
	default:
		return qualityGood
	}
}

func qualityFromNumber(value int64) string {
	switch value {
	case 2:
		return qualityUncertain
	case 3:
		return qualityBad
	default:
		return qualityGood
	}
}

func anyInt64(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case json.Number:
		n, _ := v.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		return n
	default:
		return 0
	}
}
