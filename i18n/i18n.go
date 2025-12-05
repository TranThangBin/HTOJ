package i18n

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Locale string

type Bundle struct {
	locales     map[string]map[string]string
	defaultLang string
	localesList []Locale
}

type contextKey struct{}

type contextValue struct {
	translator func(string) string
	locale     string
	locales    []Locale
}

// LoadDir loads all translation files from a directory
func LoadDir(dir string, defaultLang string) (*Bundle, error) {
	bundle := &Bundle{
		locales:     make(map[string]map[string]string),
		defaultLang: defaultLang,
		localesList: make([]Locale, 0),
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		lang := strings.TrimSuffix(entry.Name(), ext)
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", entry.Name(), err)
		}

		var translations map[string]string
		if ext == ".json" {
			if err := json.Unmarshal(data, &translations); err != nil {
				return nil, fmt.Errorf("unmarshal %s: %w", entry.Name(), err)
			}
		} else {
			if err := yaml.Unmarshal(data, &translations); err != nil {
				return nil, fmt.Errorf("unmarshal %s: %w", entry.Name(), err)
			}
		}

		bundle.locales[lang] = translations
		bundle.localesList = append(bundle.localesList, Locale(lang))
	}

	return bundle, nil
}

// HasLocale checks if a locale exists
func (b *Bundle) HasLocale(lang string) bool {
	_, ok := b.locales[lang]
	return ok
}

// GetTranslations returns all translations for a specific locale
func (b *Bundle) GetTranslations(lang string) map[string]string {
	if translations, ok := b.locales[lang]; ok {
		return translations
	}
	return b.locales[b.defaultLang]
}

// GetETag generates an ETag hash for cache validation
func (b *Bundle) GetETag(lang string) string {
	translations := b.GetTranslations(lang)
	data, _ := json.Marshal(translations)
	hash := sha256.Sum256(data)
	return `"` + hex.EncodeToString(hash[:8]) + `"` // Use first 8 bytes for shorter ETag
}

// Default returns the default language
func (b *Bundle) Default() string {
	return b.defaultLang
}

// Locales returns all available locales
func (b *Bundle) Locales() []Locale {
	return b.localesList
}

// Translator returns a translation function for a specific locale
func (b *Bundle) Translator(lang string) func(string) string {
	translations := b.locales[lang]
	if translations == nil {
		translations = b.locales[b.defaultLang]
	}

	return func(key string) string {
		if val, ok := translations[key]; ok {
			return val
		}
		return key
	}
}

// Current returns the current locale info
func (b *Bundle) Current(lang string) Locale {
	return Locale(lang)
}

// WithContext adds i18n context
func WithContext(ctx context.Context, translator func(string) string, locale string, locales []Locale) context.Context {
	return context.WithValue(ctx, contextKey{}, &contextValue{
		translator: translator,
		locale:     locale,
		locales:    locales,
	})
}

// FromContext retrieves i18n from context
func FromContext(ctx context.Context) (func(string) string, string, []Locale) {
	val := ctx.Value(contextKey{})
	if val == nil {
		return nil, "", nil
	}
	cv := val.(*contextValue)
	return cv.translator, cv.locale, cv.locales
}
