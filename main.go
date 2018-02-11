package main

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

var Total int64

type Addon struct {
	Name     string
	Cost     int64
	Quantity int64
	Rate     int64
}

func (a *Addon) Render() *vecty.HTML {
	return elem.Button(
		vecty.Markup(
			vecty.Class("pure-button"),
			vecty.MarkupIf(Total < a.Cost,
				vecty.Class("pure-button-disabled")),
			event.Click(func(e *vecty.Event) {
				if Total < a.Cost {
					return
				}
				atomic.AddInt64(&Total, -a.Cost)
				atomic.AddInt64(&a.Quantity, 1)
			}),
		),
		vecty.Text("Buy "+a.Name),
	)
}

var Addons = []Addon{
	Addon{
		Name: "TI-89",
		Cost: 10,
		Rate: 1,
	},
	Addon{
		Name: "FPGA",
		Cost: 100,
		Rate: 5,
	},
	Addon{
		Name: "GPU",
		Cost: 200,
		Rate: 8,
	},
}

func main() {
	vecty.AddStylesheet("https://unpkg.com/purecss@1.0.0/build/pure-min.css")
	vecty.AddStylesheet("https://unpkg.com/purecss@1.0.0/build/grids-responsive-min.css")
	vecty.RenderBody(&PageView{})
}

type PageView struct {
	vecty.Core
	Ticker *time.Ticker
}

func (p *PageView) Render() vecty.ComponentOrHTML {

	rate := int64(0)
	total := strconv.FormatInt(Total, 10) + " CookieCoins"
	vecty.SetTitle(total)
	if p.Ticker == nil {
		p.Start()
	}

	tr := []vecty.MarkupOrChild{}

	for i := range Addons {
		rate += Addons[i].Quantity * Addons[i].Rate
		tr = append(tr, elem.TableRow(
			elem.TableData(vecty.Text(strconv.FormatInt(Addons[i].Quantity, 10))),
			elem.TableData(vecty.Text(Addons[i].Name)),
			elem.TableData(vecty.Text(strconv.FormatInt(Addons[i].Cost, 10))),
			elem.TableData(Addons[i].Render()),
		))
	}

	t := elem.Table(vecty.Markup(
		vecty.Class("pure-table"),
		vecty.Style("width", "100%"),
	),
		elem.TableHead(elem.TableRow(
			elem.TableHeader(vecty.Text("#")),
			elem.TableHeader(vecty.Text("Name")),
			elem.TableHeader(vecty.Text("Cost")),
			elem.TableHeader(vecty.Text("Buy")),
		)),
		elem.TableBody(tr...),
	)

	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("pure-g"),
				vecty.Style("width", "790px"),
				vecty.Style("margin-left", "auto"),
				vecty.Style("margin-right", "auto"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("pure-u-1"),
					vecty.Class("pure-u-sm-1-2"),
					vecty.Style("text-align", "center"),
				),
				elem.Heading1(
					vecty.Text(total+" "),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("pure-u-1"),
					vecty.Class("pure-u-sm-1-2"),
					vecty.Style("text-align", "center"),
				),
				elem.Heading1(
					vecty.Text(strconv.FormatInt(rate, 10)+" CookieCoins/s"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("pure-u-1"),
					vecty.Style("text-align", "center"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Class("pure-button"),
						vecty.Style("margin", "1em"),
						event.Click(func(e *vecty.Event) {
							atomic.AddInt64(&Total, 1)
							vecty.Rerender(p)
						}),
					),
					vecty.Text("Mine a CookieCoin"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("pure-u-1"),
				),
				t,
			),
		),
	)
}

func (p *PageView) Start() {
	p.Ticker = time.NewTicker(time.Second)
	go func() {
		for range p.Ticker.C {
			rate := int64(0)
			for _, a := range Addons {
				rate += a.Quantity * a.Rate
			}
			atomic.AddInt64(&Total, rate)
			vecty.Rerender(p)
		}
	}()
}
