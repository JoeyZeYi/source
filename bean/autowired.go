package bean

import (
	"fmt"
	"reflect"
	"sync"
)

// Container bean 容器接口
type Container interface {
	Autowired(p interface{})
	Injection(i interface{})
	InjectionDefault(i interface{})
	Check() error
	AddInitializingBean(handler InitializingBeanHandler)
	AddInitializingBeanFunc(handler func())
}

// InitializingBeanHandler 注入完成后的回调
type InitializingBeanHandler interface {
	AfterPropertiesSet()
}

type initializingBeanFunc func()

func (fun initializingBeanFunc) AfterPropertiesSet() {
	fun()
}

type container struct {
	autowiredInterfaces map[string][]reflect.Value
	cache               map[string]reflect.Value
	defaultInterface    map[string]reflect.Value
	initializingBean    []InitializingBeanHandler
	lock                sync.Mutex
	once                sync.Once
}

func (m *container) AddInitializingBean(handler InitializingBeanHandler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if need := m.check(); len(need) == 0 {
		handler.AfterPropertiesSet()
		return
	} else {
	}
	m.initializingBean = append(m.initializingBean, handler)

}

func (m *container) AddInitializingBeanFunc(handler func()) {
	m.AddInitializingBean(initializingBeanFunc(handler))
}

// NewContainer 创建新的 bean 容器
func NewContainer() Container {
	return &container{
		autowiredInterfaces: make(map[string][]reflect.Value),
		cache:               make(map[string]reflect.Value),
		defaultInterface:    make(map[string]reflect.Value),
		lock:                sync.Mutex{},
		once:                sync.Once{},
	}
}

func (m *container) add(key string, v reflect.Value) {

	if e, has := m.cache[key]; has {
		v.Set(e)
		return
	}

	if ed, has := m.defaultInterface[key]; has {
		v.Set(ed)
	}
	m.autowiredInterfaces[key] = append(m.autowiredInterfaces[key], v)
}

func (m *container) set(key string, v reflect.Value) {

	values := m.autowiredInterfaces[key]
	delete(m.autowiredInterfaces, key)

	for _, e := range values {
		e.Set(v)
	}
}
func (m *container) check() []string {

	r := make([]string, 0, len(m.autowiredInterfaces))
	for v := range m.autowiredInterfaces {
		if _, has := m.defaultInterface[v]; !has {
			r = append(r, v)
		}
	}

	return r
}

// Autowired  声明
func (m *container) Autowired(p interface{}) {
	v, e := pkg(p)
	m.lock.Lock()

	m.add(v, e)
	m.lock.Unlock()
}

// Injection 注入
func (m *container) Injection(i interface{}) {

	pkgV, v := pkg(i)

	m.lock.Lock()

	m.cache[pkgV] = v
	m.set(pkgV, v)
	m.lock.Unlock()
}

// InjectionDefault 注入默认值， 这里注入的只会对没有被其他注入对接口生效
func (m *container) InjectionDefault(i interface{}) {

	pkgV, v := pkg(i)
	m.lock.Lock()
	m.defaultInterface[pkgV] = v
	m.lock.Unlock()

}

func (m *container) injectionAll() {

	cache := m.cache
	for k, v := range cache {
		m.set(k, v)
	}

	defaults := m.defaultInterface
	for k, v := range defaults {
		m.set(k, v)
	}

}
func pkg(i interface{}) (string, reflect.Value) {
	v := reflect.ValueOf(i)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	pkgV := key(v.Type())
	if pkgV == "" {
		panic("invalid interface")
	}
	return pkgV, v
}
func key(t reflect.Type) string {
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.String())
}

// Check 检查是否实现相关dao类
func (m *container) Check() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	var err error = nil
	m.once.Do(func() {

		m.injectionAll()
		rs := m.check()
		if len(rs) > 0 {
			err = fmt.Errorf("need:%v", rs)
			return
		}

		m.dispatchAfterPropertiesSet()
	})

	return err
}

func (m *container) dispatchAfterPropertiesSet() {

	beanHandlers := m.initializingBean
	for _, h := range beanHandlers {
		h.AfterPropertiesSet()
	}
}
