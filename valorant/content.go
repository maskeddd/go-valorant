package valorant

import (
	"context"
	"net/http"
)

type ContentService service

type Locale string

const (
	LocaleArAE Locale = "ar-AE"
	LocaleDeDE Locale = "de-DE"
	LocaleEnGB Locale = "en-GB"
	LocaleEnUS Locale = "en-US"
	LocaleEsES Locale = "es-ES"
	LocaleEsMX Locale = "es-MX"
	LocaleFrFR Locale = "fr-FR"
	LocaleIdID Locale = "id-ID"
	LocaleItIT Locale = "it-IT"
	LocaleJaJP Locale = "ja-JP"
	LocaleKoKR Locale = "ko-KR"
	LocalePlPL Locale = "pl-PL"
	LocalePtBR Locale = "pt-BR"
	LocaleRuRU Locale = "ru-RU"
	LocaleThTH Locale = "th-TH"
	LocaleTrTR Locale = "tr-TR"
	LocaleViVN Locale = "vi-VN"
	LocaleZhCN Locale = "zh-CN"
	LocaleZhTW Locale = "zh-TW"
)

type Content struct {
	Version    string        `json:"version"`
	Characters []ContentItem `json:"characters"`
	Maps       []struct {
		ContentItem
		AssetPath string `json:"assetPath"`
	} `json:"maps"`
	Chromas    []ContentItem `json:"chromas"`
	Skins      []ContentItem `json:"skins"`
	SkinLevels []ContentItem `json:"skinLevels"`
	Equips     []ContentItem `json:"equips"`
	GameModes  []struct {
		ContentItem
		AssetPath string `json:"assetPath"`
	} `json:"gameModes"`
	Sprays       []ContentItem `json:"sprays"`
	SprayLevels  []ContentItem `json:"sprayLevels"`
	Charms       []ContentItem `json:"charms"`
	CharmLevels  []ContentItem `json:"charmLevels"`
	PlayerCards  []ContentItem `json:"playerCards"`
	PlayerTitles []ContentItem `json:"playerTitles"`
	Acts         []Act         `json:"acts"`
	Ceremonies   []ContentItem `json:"ceremonies"`
}

type ContentItem struct {
	Name           string          `json:"name"`
	LocalizedNames *LocalizedNames `json:"localizedNames"`
	ID             string          `json:"id"`
	AssetName      string          `json:"assetName"`
}

type Act struct {
	Name           string          `json:"name"`
	LocalizedNames *LocalizedNames `json:"localizedNames"`
	ID             string          `json:"id"`
	IsActive       bool            `json:"is_active"`
}

type LocalizedNames struct {
	ArAE string `json:"ar-AE"`
	DeDE string `json:"de-DE"`
	EnGB string `json:"en-GB"`
	EnUS string `json:"en-US"`
	EsES string `json:"es-ES"`
	EsMX string `json:"es-MX"`
	FrFR string `json:"fr-FR"`
	IdID string `json:"id-ID"`
	ItIT string `json:"it-IT"`
	JaJP string `json:"ja-JP"`
	KoKR string `json:"ko-KR"`
	PlPL string `json:"pl-PL"`
	PtBR string `json:"pt-BR"`
	RuRU string `json:"ru-RU"`
	ThTH string `json:"th-TH"`
	TrTR string `json:"tr-TR"`
	ViVN string `json:"vi-VN"`
	ZhCN string `json:"zh-CN"`
	ZhTW string `json:"zh-TW"`
}

type ContentListOptions struct {
	Locale Locale `url:"locale,omitempty"`
}

// ListContents lists game content, optionally filtered by locale.
//
// Valorant API docs: https://developer.riotgames.com/apis#val-content-v1/GET_getContent
func (s *ContentService) ListContents(ctx context.Context, opts *ContentListOptions) (*Content, *http.Response, error) {
	u := "content/v1/contents"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var content *Content
	resp, err := s.client.Do(ctx, req, &content)
	if err != nil {
		return nil, resp, err
	}

	return content, resp, nil
}
