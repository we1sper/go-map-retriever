package mapretriever

import (
	"fmt"
	"testing"
)

func prepareTestData() interface{} {
	return map[string]interface{}{
		"author": "welsper",
		"emails": []string{
			"welsper@qq.com",
			"welsper@nit.com",
			"welsper@nekomi.com",
		},
		"details": map[string]interface{}{
			"country": "China",
			"age":     18,
			"job": map[string]string{
				"company": "China Merchants Bank",
				"salary":  "100,000/m",
			},
		},
		"others": map[interface{}]string{
			1:      "test",
			"key1": "value1",
		},
	}
}

func TestMapRetriever(t *testing.T) {
	retriever := NewMapRetriever(prepareTestData())

	fmt.Println(retriever.Get("details", "country").Trace())
	fmt.Println(retriever.Get("details", "country").Debug())
	fmt.Println(retriever.Get("details", "country").Unsafe().String())

	fmt.Println(retriever.Get("emails").At(10).Trace())
	fmt.Println(retriever.Get("emails").At(10).Debug())
	fmt.Println(retriever.Get("emails").At(10).Unsafe().String())

	fmt.Println(retriever.Get("emails").Head().Trace())
	fmt.Println(retriever.Get("emails").Head().Debug())
	fmt.Println(retriever.Get("emails").Head().Unsafe().String())

	fmt.Println(retriever.Get("emails").Tail().Trace())
	fmt.Println(retriever.Get("emails").Tail().Debug())
	fmt.Println(retriever.Get("emails").Tail().Unsafe().String())

	fmt.Println(retriever.Get("details", "job", "company").Trace())
	fmt.Println(retriever.Get("details", "job", "company").Debug())
	fmt.Println(retriever.Get("details", "job", "company").Unsafe().String())

	fmt.Println(retriever.Get("others", 1).Trace())
	fmt.Println(retriever.Get("others", 1).Debug())
	fmt.Println(retriever.Get("others", 1).Unsafe().String())

	fmt.Println(retriever.Get("others", "key1").Trace())
	fmt.Println(retriever.Get("others", "key1").Debug())
	fmt.Println(retriever.Get("others", "key1").Unsafe().String())
}
