// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package i18n

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type I18nMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取系统语言包
func NewI18nMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *I18nMessagesLogic {
	return &I18nMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *I18nMessagesLogic) I18nMessages(req *types.I18nMessagesReq) (resp *types.Envelope, err error) {
	lang, err := normalizeLang(req.Lang)
	if err != nil {
		return nil, err
	}

	data, err := readLanguagePack(lang)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, respx.NewHTTPError(http.StatusNotFound, "language pack not found")
		}
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, respx.NewHTTPError(http.StatusInternalServerError, "language pack parse failed")
	}

	return respx.Envelope(map[string]any{"messages": messages}), nil
}

func normalizeLang(lang string) (string, error) {
	key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-"))
	switch key {
	case "", "en", "en-us":
		return "en_US", nil
	case "zh", "zh-cn":
		return "zh_CN", nil
	default:
		return "", respx.NewHTTPError(http.StatusBadRequest, "unsupported language")
	}
}

func readLanguagePack(lang string) ([]byte, error) {
	filename := lang + ".json"
	for _, base := range []string{
		filepath.Join("etc", "i18n"),
		filepath.Join("backend", "etc", "i18n"),
	} {
		data, err := os.ReadFile(filepath.Join(base, filename))
		if err == nil {
			return data, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	return nil, os.ErrNotExist
}
