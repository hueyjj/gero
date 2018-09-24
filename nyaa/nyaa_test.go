package nyaa

import "testing"

func TestQuery(t *testing.T) {
	Query("https://nyaa.si/?q=psycho+pass&f=0&c=0_0&page=rss")
}
