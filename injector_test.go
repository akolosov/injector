package injector

import (
  "testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInjectors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Injector Suite")
}

var TestInjector InjectorInterface
var FirstObject *FirstObjectStruct

type FirstObjectStruct struct {
  InjectorInterface
  TestInfection string     `infection:"test_infection"`
  TestInjection string     `injection:"test_injection"`
}

var _ = BeforeSuite(func() {
  TestInjector = NewInjector()
  FirstObject = &FirstObjectStruct{InjectorInterface: TestInjector, TestInfection: "test_infection"}
  FirstObject.Register("test_injection", "test_injection")
  FirstObject.Inject(FirstObject)
})

var _ = Describe("Injector", func() {
  Describe("#Injector", func() {
    It("Should be not nil", func() {
      Expect(TestInjector).ShouldNot(BeNil())
    })

    It("Should be singleton", func() {
      Expect(TestInjector).Should(Equal(NewInjector()))
    })

    It("Should give a value", func() {
      TestInjector.Register("testString", "test")
      TestInjector.Register("testBool", true)
      TestInjector.Register("testInt", 42)

      Expect(TestInjector.Invoke("testString")).Should(ContainSubstring("test"))
      Expect(TestInjector.Invoke("testBool")).Should(BeTrue())
      Expect(TestInjector.Invoke("testInt")).Should(Equal(42))

      TestInjector.Unregister("testString")
      Expect(TestInjector.Invoke("testString")).Should(BeNil())
    })
  })

  Describe("#Infection", func() {
    It("Should be infect value", func() {
      Expect(TestInjector.Invoke("test_infection")).Should(Equal("test_infection"))
    })

    It("Should be run for each object", func() {
      TestInjector.Each(func (obj interface{}) {
        Expect(obj).ShouldNot(BeNil())
      })
    })
  })

  Describe("#Injection", func() {
    It("Should be inject value", func() {
      Expect(TestInjector.Invoke("test_injection")).Should(Equal("test_injection"))
    })

    It("Should be run for each object", func() {
      TestInjector.Each(func (obj interface{}) {
        Expect(obj).ShouldNot(BeNil())
      })
    })
  })

  Describe("Test Non Struct", func() {
    It("Should be nil", func(){
      Expect(TestInjector.Inject("testing 1 2 3")).Should(BeNil())
    })
  })
})