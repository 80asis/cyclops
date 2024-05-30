// entity/entity.go
package entity

// Entity represents an entity with its properties
type Entity struct {
	EntityID   string `default:""`
	EntityKind string `default:""`
	OpType     string `default:""`
}
