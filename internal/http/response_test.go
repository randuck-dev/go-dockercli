package http

import "testing"

func TestContentType(t *testing.T) {

	t.Run("application/json is set", func(t *testing.T) {
		resp := Response{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		content_type, err := resp.ContentType()

		if err != nil {
			t.Errorf("unexpected error occured: %s", err)
		}

		if content_type != "application/json" {
			t.Errorf("got %s want %s", content_type, "application/json")
		}
	})

	t.Run("header is not set", func(t *testing.T) {
		resp := Response{}

		_, err := resp.ContentType()

		if err != ErrNoContentTypefound {
			t.Errorf("got %s want %s", err, ErrNoContentTypefound)
		}
	})
}
