package shortcut

import (
	"os"
	"time"

	"log/slog"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/lederniermetre/shortcut/pkg/shortcut/gen/client"
)

func GetClient() *apiclient.ShortcutAPI {
	// Hack to parse end_date "2023-01-19"
	strfmt.DateTimeFormats = append(strfmt.DateTimeFormats, time.DateOnly)

	// create the transport
	transport := httptransport.New("api.app.shortcut.com", "", nil)
	return apiclient.New(transport, strfmt.Default)
}

func GetAuth() runtime.ClientAuthInfoWriter {
	if os.Getenv("SHORTCUT_API_TOKEN") == "" {
		slog.Error("SHORTCUT_API_TOKEN is empty")
		os.Exit(1)
	}

	return httptransport.APIKeyAuth("Shortcut-Token", "header", os.Getenv("SHORTCUT_API_TOKEN"))
}
