package web

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"
)

type CacheStore struct {
	Lock  *sync.RWMutex
	Cache map[string]Page
}

func NewCache() CacheStore {
	return CacheStore{Lock: new(sync.RWMutex), Cache: make(map[string]Page)}
}

func (m CacheStore) Put(k string, v Page) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Cache[k]; !ok {
		m.Cache[k] = v
	}
}
func (m CacheStore) Get(k string) (v Page, inCache bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Cache[k]; ok {
		v = m.Cache[k]
		inCache = true
	}
	return
}

type TemplateCacheStore struct {
	Lock  *sync.RWMutex
	Cache map[string]*template.Template
}

func NewTemplateCache() TemplateCacheStore {
	return TemplateCacheStore{Lock: new(sync.RWMutex), Cache: make(map[string]*template.Template)}
}

func (m TemplateCacheStore) Put(k string, v *template.Template) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Cache[k]; !ok {
		m.Cache[k] = v
	}
}
func (m TemplateCacheStore) Get(k string) (v *template.Template, inCache bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Cache[k]; ok {
		v = m.Cache[k]
		inCache = true
	}
	return
}
func (m TemplateCacheStore) JGet(k string) (v *template.Template) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Cache[k]; ok {
		v = m.Cache[k]
	}
	return
}

var TemplateFuncStore template.FuncMap

func Netadd(x, v float64) float64 {
	return v + x
}

func Netsubs(x, v float64) float64 {
	return v - x
}

func Netmultiply(x, v float64) float64 {
	return v * x
}

func Netdivided(x, v float64) float64 {
	return v / x
}

type NoStruct struct {
	/* emptystruct */
}

func NetsessionGet(key string, s *sessions.Session) string {
	return s.Values[key].(string)
}

func UrlAtZ(url, base string) (isURL bool) {
	isURL = strings.Index(url, base) == 0
	return
}

func NetsessionDelete(s *sessions.Session) string {
	//keys := make([]string, len(s.Values))

	//i := 0
	for k := range s.Values {
		// keys[i] = k.(string)
		NetsessionRemove(k.(string), s)
		//i++
	}

	return ""
}

func NetsessionRemove(key string, s *sessions.Session) string {
	delete(s.Values, key)
	return ""
}
func NetsessionKey(key string, s *sessions.Session) bool {
	if _, ok := s.Values[key]; ok {
		//do something here
		return true
	}

	return false
}

func NetsessionGetInt(key string, s *sessions.Session) interface{} {
	return s.Values[key]
}

func NetsessionSet(key string, value string, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}
func NetsessionSetInt(key string, value interface{}, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}

func Netimportcss(s string) string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%%s\" /> ", s)
}

func Netimportjs(s string) string {
	return fmt.Sprintf("<script type=\"text/javascript\" src=\"%%s\" ></script> ", s)
}

func Formval(s string, r *http.Request) string {
	return r.FormValue(s)
}

func Equalz(args ...interface{}) bool {
	if args[0] == args[1] {
		return true
	}
	return false
}
func Nequalz(args ...interface{}) bool {
	if args[0] != args[1] {
		return true
	}
	return false
}

func Netlt(x, v float64) bool {
	if x < v {
		return true
	}
	return false
}
func Netgt(x, v float64) bool {
	if x > v {
		return true
	}
	return false
}
func Netlte(x, v float64) bool {
	if x <= v {
		return true
	}
	return false
}

func GetLine(fname string, match string) int {
	intx := 0
	file, err := os.Open(fname)
	if err != nil {
		color.Red("Could not find a source file")
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		intx = intx + 1
		if strings.Contains(scanner.Text(), match) {

			return intx
		}

	}

	return -1
}
func Netgte(x, v float64) bool {
	if x >= v {
		return true
	}
	return false
}

type Page struct {
	Title      string
	Body       []byte
	IsResource bool
	R          *http.Request
	Session    *sessions.Session
}
