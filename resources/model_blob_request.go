/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type BlobRequest struct {
	Key
	Attributes BlobRequestAttributes `json:"attributes"`
}
type BlobRequestResponse struct {
	Data     BlobRequest `json:"data"`
	Included Included    `json:"included"`
}

type BlobRequestListResponse struct {
	Data     []BlobRequest `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustBlobRequest - returns BlobRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustBlobRequest(key Key) *BlobRequest {
	var blobRequest BlobRequest
	if c.tryFindEntry(key, &blobRequest) {
		return &blobRequest
	}
	return nil
}
