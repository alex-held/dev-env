package action

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl/internal/logging"

	"github.com/alex-held/devctl/internal/devctlpath"
	"github.com/alex-held/devctl/internal/sdkman"
	"github.com/alex-held/devctl/internal/testutils"
)

func testExists(g *goblin.G, fs afero.Fs, expected, msg string) {
	g.Helper()
	exists, err := afero.Exists(fs, expected)
	if err != nil {
		g.Fatalf("error occurred while testing whether file or dir exists; path=%s; error=%v\n", expected, err)
	}
	Expect(exists).Should(BeTrue(), "%s; path=%s", msg, expected)
}

type ActionTestFixture struct {
	g        *goblin.G
	actions  *Actions
	fs       afero.Fs
	logger   *logging.Logger
	mux      *http.ServeMux
	teardown testutils.Teardown
	context  context.Context
	client   *sdkman.Client
	pather   devctlpath.Pather
}

func SetupFixtureDeps(g *goblin.G, fs afero.Fs, pather devctlpath.Pather, logger *logging.Logger, td func()) (fixture *ActionTestFixture) {
	const baseURLPath = "/2"
	mux := http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: ClientIn.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "ClientIn.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)
	teardown := func() {
		td()
		server.Close()
	}
	client := sdkman.NewSdkManClient(
		sdkman.URLOptions(server.URL+baseURLPath),
		sdkman.FileSystemOption(fs),
		sdkman.HTTPClientOption(&http.Client{}),
	)

	actions := NewActions(WithFs(fs), WithSdkmanClient(client), WithPather(pather), WithLogger(logger))

	fixture = &ActionTestFixture{
		g:        g,
		actions:  actions,
		logger:   logger,
		mux:      mux,
		pather:   pather,
		fs:       fs,
		teardown: teardown,
		context:  context.Background(),
		client:   client,
	}
	return fixture
}

func SetupFixture(g *goblin.G) (fixture *ActionTestFixture) {
	return SetupFixtureDeps(g, afero.NewMemMapFs(), devctlpath.NewPather(), testutils.NewLogger(), func() {})
}

func TestNewActions(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	fs := afero.NewMemMapFs()
	pather := devctlpath.NewPather()
	client := sdkman.NewSdkManClient()

	g.Describe("NewActions", func() {
		g.It("WithFs", func() {
			actions := NewActions(WithFs(fs))
			Expect(actions.Options.Fs).Should(Equal(fs))
		})

		g.It("WithPather", func() {
			actions := NewActions(WithPather(pather))
			Expect(actions.Options.Pather).Should(Equal(pather))
		})

		g.It("WithSdkmanClient", func() {
			actions := NewActions(WithSdkmanClient(client))
			Expect(actions.Options.Client).Should(Equal(client))
		})
	})
}
