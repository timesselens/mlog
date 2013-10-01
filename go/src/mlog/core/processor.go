package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mlog/io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// Processors, running instances (goroutines)
type Processors struct {
	mu   *sync.RWMutex
	List []io.Processor
}

var processors = Processors{List: make([]io.Processor, 0), mu: new(sync.RWMutex)}

// Factories
type ProcessorFactory struct {
	factory  map[string]reflect.Type
	mu       *sync.RWMutex
	basePath string
}

var processorFactory = ProcessorFactory{
	factory: map[string]reflect.Type{
		"stdin":     reflect.TypeOf(io.Stdin{}),
		"syslog":    reflect.TypeOf(io.SyslogListener{}),
		"grep":      reflect.TypeOf(io.Grep{}),
		"debug":     reflect.TypeOf(io.Debug{}),
		"websocket": reflect.TypeOf(io.WebSocket{}),
	},
	mu: new(sync.RWMutex),
}

func (procs *Processors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	basePath := "/"
	parts := strings.SplitN(r.URL.Path[len(basePath):], "/", 5)

	fmt.Println(r.Method, r.URL.Path, parts)
	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/":
			enc := json.NewEncoder(w)
			enc.Encode(procs.GetList())
		}
	case "POST":
		switch r.URL.Path {
		case "/":
			body, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			x, err := CreateProcessor(string(body))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			log.Printf("new proc %T %s\n", x, x)
			procs.Append(x)
		}
	}
}

func PipeProcessor(i, j int) (err error) {
	var s, t io.Processor
	processors.mu.RLock()
	defer processors.mu.RUnlock()
	if i > len(processors.List) {
		return errors.New("i out of bound")
	}
	if j > len(processors.List) {
		return errors.New("j out of bound")
	}
	if i == j {
		return errors.New("i and j can not be the same")
	}
	s = processors.List[i]
	t = processors.List[j]
	if s != nil && t != nil {
		fmt.Printf("piping [%d] %s to [%d] %s\n", i, s, j, t)
		io.Pipe(s, t)
	}
	return nil
}

func GetProcessorList() (list []io.Processor) {
	processors.mu.RLock()
	defer processors.mu.RUnlock()
	return processors.List
}

func AppendProcessor(p io.Processor) {
	processors.Append(p)
}

func (procs *Processors) Append(p io.Processor) {
	procs.mu.Lock()
	defer procs.mu.Unlock()
	processors.List = append(processors.List, p)
}

func (procs *Processors) GetList() (list []io.Processor) {
	procs.mu.RLock()
	defer procs.mu.RUnlock()
	return processors.List
}

func DeleteProcessor(i int) (p io.Processor, err error) {
	return processors.Delete(i)
}

func (procs *Processors) Delete(i int) (p io.Processor, err error) {
	procs.mu.Lock()
	defer procs.mu.Unlock()

	if i < 0 || i >= len(procs.List) {
		return nil, errors.New("i out of bounds")
	}

	tbd := processors.List[i]
	processors.List = append(processors.List[0:i], processors.List[i+1:]...)
	tbd.Exit()
	tbd.Close()
	return tbd, err
}

func CreateProcessor(d string) (p io.Processor, err error) {
	return processorFactory.CreateProcessor(d)
}

func (pf *ProcessorFactory) CreateProcessor(d string) (p io.Processor, err error) {
	if d[0:1] == "{" {
		return pf.createProcessorFromJSON(d)
	} else {
		return pf.createProcessorFromURL(d)
	}
}

func (pf *ProcessorFactory) createProcessorFromURL(s string) (p io.Processor, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.New("could not parse url from data")
	}
	name := u.Scheme
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	T := pf.factory[name]

	if T == nil {
		return nil, errors.New(fmt.Sprintf("factory for %s not found", name))
	}

	ptr := reflect.New(T)
	proc := ptr.Elem()
	pp := ptr.Interface().(io.Processor)
	err = pp.Init()
	if err != nil {
		fmt.Println("unable to initialize regex processor", err)
	} else {
		fmt.Printf("created proc %#v\n", proc.Interface())
	}
	return pp, err
}

func (pf *ProcessorFactory) createProcessorFromJSON(data string) (p io.Processor, err error) {

	var store interface{}
	json.Unmarshal([]byte(data), &store)

	m, ok := store.(map[string]interface{})

	if !ok {
		return nil, errors.New("could not parse JSON structure as map")
	}

	name, ok := m["type"].(string)

	if !ok || len(name) < 1 {
		return nil, errors.New("could not get type from data")
	}

	pf.mu.RLock()
	defer pf.mu.RUnlock()
	T := pf.factory[name]

	if T == nil {
		return nil, errors.New(fmt.Sprintf("factory for %s not found", name))
	}

	ptr := reflect.New(T)
	proc := ptr.Elem()

	if _, ok := ptr.Interface().(io.Processor); !ok {
		return nil, errors.New(fmt.Sprintf("type returned by factory %s is not an io.Processor but %T", name, ptr.Interface()))
	}

	for k, v := range m {
		if k == "type" {
			continue
		}

		k = strings.ToUpper(k[0:1]) + k[1:]
		field := proc.FieldByName(k)

		if !field.IsValid() {
			log.Printf("field %s of %s is not a valid field (for value '%s')", k, proc.Type(), v)
			continue
		}

		if !field.CanSet() {
			log.Printf("field %s of %s is not settable (for value '%s')", k, proc.Type(), v)
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if str, ok := v.(string); ok {
				field.SetString(str)
			}
		case reflect.Bool:
			if b, ok := v.(bool); ok {
				field.SetBool(b)
			}
		case reflect.Float64:
			if i, ok := v.(float64); ok {
				field.SetFloat(i)
			}
		case reflect.Int:
			if i, ok := v.(int64); ok {
				field.SetInt(i)
			}
		}
	}
	pp := ptr.Interface().(io.Processor)
	err = pp.Init()
	if err != nil {
		fmt.Println("unable to initialize regex processor", err)
	} else {
		fmt.Printf("created proc %#v\n", proc.Interface())
	}
	return pp, err
}
