// Пакет предоставляет функции Dependency Injection/Infection (DI/I)
package injector

import (
  "reflect"
)

// Интерфейсный типа, предоставляющий доступ к основным функциям модуля
// Анонимно встраивается во все структуры, где необходимо делать DI/I
type InjectorInterface interface {
  Invoke(string) interface{}
  Register(string, interface{})
  Unregister(string)
  Inject(interface{}) error
  Each(func(interface{}))
}

// Основное KV-хранилище для DI/I
type Injector map[string]interface{}

// Singleton-объект для DI/I механизма
var masterInjector Injector

// Возвращает ссылку на созданное KV-хранилище в виде интерфейса вызовов
func NewInjector() InjectorInterface {
  if masterInjector == nil {
    masterInjector = make(Injector)
  }
  return masterInjector
}

// Добавляет новый объект в KV-хранилище по ключу key
// Если объект уже существует - замещает его новым объектом
func (in Injector) Register(key string, value interface{}) {
  in[key] = value
}

// Удаляет существующий объект по ключу key из KV-хранилища
func (in Injector) Unregister(key string) {
  delete(in, key)
}

// Возвращает сохраненный объект по ключу key из KV-хранилища
// Если объекта не существует - возвращает nil
func (in Injector) Invoke(key string) interface{} {
  return in[key]
}

// Функция осуществляет вызов функции each для каждого объекта в K/V хранилище
func (in Injector) Each(each func(interface{})) {
  for _, obj := range in {
    each(obj)
  }
}

// Функция производит инъекцию и заражение объектов из/в KV-хранилище в структуру val
// Инъекция (injection) происходит по тэгу `injection:"ключ-в-KV-хранилище"` в структуре val
// Заражение - внедрение в KV-хранилище объектов из структуры val
// Заражение (infection) происходит по тэгу `infection:"ключ-в-KV-хранилище"` в структуре val
// Возвращает ошибку, если что-то пошло не так
func (in Injector) Inject(val interface{}) (err error) {

  v := reflect.ValueOf(val)

  for v.Kind() == reflect.Ptr {
    v = v.Elem()
  }

  if v.Kind() != reflect.Struct {
    return nil
  }

  t := v.Type()
  for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    // Ищем поля для инъекции
    tag := t.Field(i).Tag.Get("injection")
    // Если поле структуры может принять инъекцию и полученный тэг не пустой...
    if f.CanSet() && len(tag) > 0 {
      // ... и значение объекта не nil...
      if s := in.Invoke(tag); s != nil {
        // ... записываем значение из KV-хранилища в поле структуры
        f.Set(reflect.ValueOf(s))
      }
    // Ищем поля для заражения
    } else if tag := t.Field(i).Tag.Get("infection"); len(tag) > 0 {
      // Если полученный тэг не пустой и значение поля можно преобразовать в interface{}
      if f.CanInterface() {
        // ... записываем значение поля структуры в KV-хранилище
        in.Register(tag, f.Interface())
      }
    }
  }

  return
}

