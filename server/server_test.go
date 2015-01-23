package server_test

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/server/fakes"

	. "github.com/cloudfoundry-incubator/cli-plugin-repo/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	var (
		repoServer  RepoServer
		fakeHandles *fakes.FakeServerHandles
	)

	BeforeEach(func() {
		fakeHandles = &fakes.FakeServerHandles{}
		repoServer = NewRepoServer(0, "0.0.0.0", fakeHandles)
		go repoServer.Serve()
	})

	AfterEach(func() {
		repoServer.Stop()
	})

	It("has a /list endpoint", func() {
		_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/list", repoServer.Port()))
		Ω(err).ToNot(HaveOccurred())
		Ω(fakeHandles.ListPluginsCallCount()).To(Equal(1))
	})

})
