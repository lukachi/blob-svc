/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type RefreshRequest struct {
	Key
	Attributes RefreshRequestAttributes `json:"attributes"`
}
type RefreshRequestResponse struct {
	Data     RefreshRequest `json:"data"`
	Included Included       `json:"included"`
}

type RefreshRequestListResponse struct {
	Data     []RefreshRequest `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustRefreshRequest - returns RefreshRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRefreshRequest(key Key) *RefreshRequest {
	var refreshRequest RefreshRequest
	if c.tryFindEntry(key, &refreshRequest) {
		return &refreshRequest
	}
	return nil
}
