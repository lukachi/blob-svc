/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SignUpRequest struct {
	Key
	Attributes SignUpRequestAttributes `json:"attributes"`
}
type SignUpRequestResponse struct {
	Data     SignUpRequest `json:"data"`
	Included Included      `json:"included"`
}

type SignUpRequestListResponse struct {
	Data     []SignUpRequest `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustSignUpRequest - returns SignUpRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSignUpRequest(key Key) *SignUpRequest {
	var signUpRequest SignUpRequest
	if c.tryFindEntry(key, &signUpRequest) {
		return &signUpRequest
	}
	return nil
}
