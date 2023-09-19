/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Key struct {
	ID           string `json:"id"`
	ResourceType string `json:"type"`
}

func (r *Key) GetKey() Key {
	return *r
}

func (r Key) GetKeyP() *Key {
	return &r
}

func (r Key) AsRelation() *Relation {
	return &Relation{
		Data: r.GetKeyP(),
	}
}
