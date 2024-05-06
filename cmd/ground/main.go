package main

import (
	"fmt"
	"reflect"
	"schedule_task_command/entry/e_command_template"
)

func main() {
	d := "aaa"
	template := &e_command_template.CommandTemplateCreate{
		Name:        "Example Template",
		Visible:     true,
		Protocol:    e_command_template.Http,
		Timeout:     30,
		Description: &d,
		Host:        "localhost",
		Port:        "8080",
	}
	inspectStruct(template)
}

func inspectStruct(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Iterate over each field of the struct.
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fmt.Println(val)
		fmt.Println(field)
		fmt.Println(val.Type().Field(i))
		typeField := val.Type().Field(i)
		fmt.Printf("Field Name: %s\n", typeField.Name)
		fmt.Printf("Field Type: %s\n", field.Type())
		fmt.Printf("Value: %v\n", field.Interface())

		// Check if the field is a zero value.
		if field.IsZero() {
			fmt.Println("Value is zero")
		} else {
			fmt.Println("Value is not zero")
		}

		// Additional processing based on field type.
		switch field.Kind() {
		case reflect.String:
			fmt.Println("This is a string")
		case reflect.Bool:
			fmt.Println("This is a boolean")
		case reflect.Int, reflect.Int32, reflect.Int64:
			fmt.Println("This is an integer")
		case reflect.Ptr:
			fmt.Println(field.Kind())
			val2 := field.Elem()
			fmt.Println(val2.Kind())
		case reflect.Struct:
			fmt.Println("This is a struct")
			// Optionally, recursively inspect nested structs.
			// inspectStruct(field.Interface())
		}
		fmt.Println()
	}
}
