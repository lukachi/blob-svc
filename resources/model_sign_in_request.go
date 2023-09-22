/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SignInRequest struct {
	Key
	Attributes SignInRequestAttributes `json:"attributes"`
}
type SignInRequestResponse struct {
	Data     SignInRequest `json:"data"`
	Included Included      `json:"included"`
}

type SignInRequestListResponse struct {
	Data     []SignInRequest `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustSignInRequest - returns SignInRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSignInRequest(key Key) *SignInRequest {
	var signInRequest SignInRequest
	if c.tryFindEntry(key, &signInRequest) {
		return &signInRequest
	}
	return nil
}
