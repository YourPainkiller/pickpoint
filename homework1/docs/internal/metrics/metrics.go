package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	contactLabel = "contact"
	handlerLabel = "handler"
	codeLabel    = "code"
)

var (
	acceptedOrdersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_accepted_orders_total",
		Help: "total number of accepted orders",
	})

	okRespByHandlerTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pvz_ok_response_by_handler_total",
		Help: "total number of ok responses in handler total",
	}, []string{
		handlerLabel,
	})

	badRespByHandlerTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pvz_bad_response_by_handler_total",
		Help: "total number of bad responses in handler total with resp code",
	}, []string{
		handlerLabel,
		codeLabel,
	})
)

func AddAcceptedOrder() {
	acceptedOrdersTotal.Inc()
}

func IncOkRespByHandler(handler string) {
	okRespByHandlerTotal.With(prometheus.Labels{
		handlerLabel: handler,
	}).Inc()
}

func IncBadRespByHandler(handler string, code int) {
	badRespByHandlerTotal.With(prometheus.Labels{
		handlerLabel: handler,
		codeLabel:    fmt.Sprint(code),
	}).Inc()
}
