package injector

import (
	"fmt"
	"html/template"
)

func ScriptInjector() {
	fmt.Println("Injecting")

	tmpl, err := template.New("name").Parse()
}
