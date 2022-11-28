// package events contains message related types and functions to be consumed by wallenda
package events

type Message struct {
	Name      string
	Group     string
	Resources ResourceConfig
}

type ResourceConfig struct {
	CPU    string
	Memory string
}
