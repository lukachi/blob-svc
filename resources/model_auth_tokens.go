/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AuthTokens struct {
	Key
	Attributes AuthTokensAttributes `json:"attributes"`
}
type AuthTokensResponse struct {
	Data     AuthTokens `json:"data"`
	Included Included   `json:"included"`
}

type AuthTokensListResponse struct {
	Data     []AuthTokens `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustAuthTokens - returns AuthTokens from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAuthTokens(key Key) *AuthTokens {
	var authTokens AuthTokens
	if c.tryFindEntry(key, &authTokens) {
		return &authTokens
	}
	return nil
}
