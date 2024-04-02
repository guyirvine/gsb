package gsb

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

func (h *Host) getHandlerName(rawName string) string {
	// len("Handler") = 7

	log.Debug("handlerLoader.getHandlerName.rawName: ", rawName)
	if len(rawName) <= 7 {
		log.Fatalf("Handler name too short. Suffix must be 'Handler'. name: %s", rawName)
	}
	if rawName[len(rawName)-7:] != "Handler" {
		log.Fatalf("Handler not named correctly. Suffix must be 'Handler'. name: %s", rawName)
	}
	name := rawName[0 : len(rawName)-7]

	return name
}

func (h *Host) addHandler(handlerName string, ha Handler) *HandlerDefinition {
	name := h.getHandlerName(handlerName)
	hd := &HandlerDefinition{name, ha, []APRDefinition{}}
	h.handlers[name] = hd

	return hd
}

func (h *Host) initHandlers() {
	h.handlers = make(map[string]*HandlerDefinition)
}

// WireUpHandler injects AppResources by name
func (h *Host) WireUpHandler(handlerName string, hd *HandlerDefinition) {
	ha := hd.handler
	e := reflect.ValueOf(ha).Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name

		if varName == "Host" {
			val := reflect.ValueOf(h)
			e.Field(i).Set(val)
			continue
		}

		// TODO type check before injection
		// varType := e.Type().Field(i).Type

		apr := h.aprList[varName]
		if apr != nil {
			//val := reflect.ValueOf(apr).Elem()
			val := reflect.ValueOf(apr)
			e.Field(i).Set(val)

			hd.aprList = append(hd.aprList, apr)

			log.Infof("Handler loading, %s. Injected APR: %s", handlerName, varName)
		}
	}
}

func (h *Host) LoadHandler(ha Handler) {
	handlerName := getTypeName(ha)
	log.Info("Handler loading. ", handlerName)
	hd := h.addHandler(handlerName, ha)
	h.WireUpHandler(handlerName, hd)
	err := ha.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("handlerLoader.LoadHandler.Handler, %s. APR Count: %d", handlerName, len(hd.aprList))
	log.Info("Handler loaded. ", handlerName)
}

func (h *Host) loadHandlers() {
	h.initHandlers()
	h.LoadHandler(new(DummyHandler))
}
