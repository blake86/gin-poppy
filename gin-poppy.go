package gin_poppy

import (
	"expvar"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	// DefaultPrefix url prefix of gin-poppy
	DefaultPrefix = "/debug/poppy"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func Register(r *gin.Engine, prefixOptions ...string) {

	pop := NewPoppy()
	r.Use(PoppyHandler(pop))

	prefix := getPrefix(prefixOptions...)
	//fmt.Print(prefix)

	//展示相关的handler绑定
	prefixRouter := r.Group(prefix)
	{
		prefixRouter.GET("/", pop.OutputHandler)
		prefixRouter.GET("/raw", pop.RawDataHandler)
		prefixRouter.GET("/expvar", func(c *gin.Context) {
			w := c.Writer
			c.Header("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte("{\n"))
			first := true
			expvar.Do(func(kv expvar.KeyValue) {
				if !first {
					w.Write([]byte(",\n"))
				}
				first = false
				fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
			})
			w.Write([]byte("\n}\n"))
			c.AbortWithStatus(200)
		})
	}
}

type inputStat struct {
	path    string
	code    int
	latency int64
}

type Poppy struct {
	inputCh chan inputStat
	Store   PoppyStore
}

func PoppyHandler(pop *Poppy) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		pop.inputCh <- inputStat{
			path:    c.FullPath(),
			code:    c.Writer.Status(),
			latency: int64(time.Since(t)),
		}

	}
}

func NewPoppy() *Poppy {
	pop := &Poppy{}
	pop.inputCh = make(chan inputStat)
	pop.Store = NewPoppyStore()

	go pop.inputReceiver()

	return pop
}

func (p *Poppy) inputReceiver() {
	for {
		select {
		case input := <-p.inputCh:
			p.inputHanlder(input)
		}
	}
}

func (p *Poppy) inputHanlder(stat inputStat) {
	p.Store.Add(stat.path, stat.code, stat.latency)
}

func (p *Poppy) OutputHandler(ctx *gin.Context) {
	p.GenerateResult()
	ctx.JSON(http.StatusOK,
		gin.H{
			"poppy": p.GenerateResult(),
		})
}

func (p *Poppy) GenerateResult() PoppyStat {
	return p.Store.GenerateStat()
}

func (p *Poppy) GenerateRawData() PoppyRawData {
	return p.Store.GenerateRawData()
}

func (p *Poppy) RawDataHandler(ctx *gin.Context) {
	fmt.Println(p.Store)
	ctx.JSON(http.StatusOK,
		gin.H{
			"raw": gin.H{
				"UriStore":    p.GenerateRawData(),
				"GlobalStore": p.Store.GlobalStore,
			},
		})
}
