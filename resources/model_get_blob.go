/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GetBlob struct {
	Key
	Attributes    GetBlobAttributes    `json:"attributes"`
	Relationships GetBlobRelationships `json:"relationships"`
}
type GetBlobResponse struct {
	Data     GetBlob  `json:"data"`
	Included Included `json:"included"`
}

type GetBlobListResponse struct {
	Data     []GetBlob `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustGetBlob - returns GetBlob from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGetBlob(key Key) *GetBlob {
	var getBlob GetBlob
	if c.tryFindEntry(key, &getBlob) {
		return &getBlob
	}
	return nil
}
