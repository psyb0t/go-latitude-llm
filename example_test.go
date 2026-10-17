package latitude_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-latitude-llm"
)

// This README example needs live API credentials, so it only compiles.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := latitude.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("LATITUDE_API_KEY"),
			)

			return nil
		},
	)

	client, err := latitude.NewClientWithResponses(
		"https://api.latitude.so",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListProjectsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
