package sitemap

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ReanSn0w/gew/v2/pkg/view"
	"github.com/ReanSn0w/sitemap/pkg/xml"
)

type Sitemap interface {
	Get(perpage, page uint) []Element

	Index(perpage uint) view.View
	Page(perpage, page uint) view.View

	Handler() http.Handler

	Run(timeout time.Duration)
	Stop()
}

type SitemapBuilder func() ([]Element, error)

func NewSitemap(baseurl string, builders ...SitemapBuilder) Sitemap {
	return &sitemap{
		Baseurl:  baseurl,
		Elements: make([]Element, 0),

		mutex:    sync.RWMutex{},
		done:     make(chan bool),
		builders: builders,
	}
}

type sitemap struct {
	Baseurl  string
	Elements []Element

	mutex    sync.RWMutex
	timer    *time.Timer
	done     chan bool
	builders []SitemapBuilder
}

func (sm *sitemap) Run(timeout time.Duration) {
	log.Println("Запуск создания карты сайта")

	sm.timer = time.NewTimer(time.Second * 3)
	log.Println("Первое создание карты будет через 3 секунды")

	go func() {
		for {
			select {
			case <-sm.timer.C:
				sm.build()
				sm.timer.Reset(timeout)
			case <-sm.done:
				log.Println("Создание карты сайта завершено ocтановлено")
				sm.timer.Stop()
				sm.timer = nil
				return
			}
		}
	}()
}

// Метод для построения карты сайта
func (sm *sitemap) Stop() {
	sm.done <- true

	for {
		time.Sleep(time.Second)
		if sm.IsWorked() {
			continue
		}
	}
}

// Метод для определиния статуса работы системы построения карты сайта
func (sm *sitemap) IsWorked() bool {
	return sm.timer != nil
}

// Вернет документ для рендеринга через пакет xml
func (sm *sitemap) Index(perpage uint) view.View {
	pages := (len(sm.Elements) / int(perpage)) + 1
	if len(sm.Elements)/int(perpage) != 0 {
		pages += 1
	}

	lastmod := time.Now().Format("2006-01-02T15:04:05-07:00")

	return view.Group(
		xml.String("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"),
		xml.NewNode(
			false,
			"sitemapindex",
			xml.Parameters{
				"xmlns:xsi":          "http://www.w3.org/2001/XMLSchema-instance",
				"xmlns":              "http://www.sitemaps.org/schemas/sitemap/0.9",
				"xsi:schemaLocation": "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd",
			},
		)(view.For(pages, func(i int) view.View {
			return xml.NewNode(false, "sitemap", nil)(
				xml.NewNode(true, "loc", nil)(xml.Url(fmt.Sprintf("%s/sitemap.xml?count=%v&page=%v", sm.Baseurl, perpage, i+1))),
				xml.NewNode(true, "lastmod", nil)(xml.String(lastmod)),
			)
		})),
	)
}

// Вернет документ для рендеринга через пакет xml
func (sm *sitemap) Page(perpage, page uint) view.View {
	elements := sm.Get(perpage, page)

	return view.Group(
		xml.String("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"),
		xml.NewNode(
			false,
			"urlset",
			xml.Parameters{"xmlns": "http://www.sitemaps.org/schemas/sitemap/0.9"},
		)(view.For(len(elements), func(i int) view.View {
			return elements[i]
		})),
	)
}

// Метод для вывода страници с елементами списка
// Сделан публичным для теста
func (sm *sitemap) Get(perpage, page uint) []Element {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	count := len(sm.Elements)
	fisrtPageIndex := int(perpage * (page - 1))
	lastPageIndex := int(perpage * page)

	if fisrtPageIndex < 0 {
		fisrtPageIndex = 0
	}

	if count <= fisrtPageIndex {
		return []Element{}
	}

	if count-1 < lastPageIndex {
		return sm.Elements[fisrtPageIndex:]
	}

	return sm.Elements[fisrtPageIndex:lastPageIndex]
}

// Метод для обновления списка элементов
func (sm *sitemap) build() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	batch := make([]Element, 0)

	for _, builder := range sm.builders {
		elements, err := builder()
		if err != nil {
			log.Println("Ошибка при обновлении индекса:", err)
			continue
		}

		batch = append(batch, elements...)
	}

	sm.Elements = batch
}
