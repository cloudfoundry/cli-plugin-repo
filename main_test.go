package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"
	"io/ioutil"
	"time"
)

var buildPath string

var _ = SynchronizedBeforeSuite(func() []byte {
	path, buildErr := gexec.Build("github.com/cloudfoundry-incubator/cli-plugin-repo")
	Expect(buildErr).NotTo(HaveOccurred())
	return []byte(path)
}, func(data []byte) {
	buildPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("Integration", func() {
	var (
		session *gexec.Session
		err     error
		port    string
	)

	BeforeEach(func() {
		port = strconv.Itoa(8080 + GinkgoParallelNode())
		session, err = gexec.Start(
			exec.Command(buildPath, "-p", port),
			GinkgoWriter,
			GinkgoWriter,
		)
		Expect(err).NotTo(HaveOccurred())

		time.Sleep(time.Second)
	})

	AfterEach(func() {
		session.Kill()
	})

	Describe("/", func() {
		It("returns HTML we expect", func() {
			client := http.DefaultClient
			response, err := client.Get("http://127.0.0.1:" + port)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(BeSuccessful())

			b, err := ioutil.ReadFile("fixtures/index.html")
			Expect(err).NotTo(HaveOccurred())

			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(string(b)))
		})
	})

	Describe("/list", func() {
		It("returns json that looks like we expect it", func() {
			client := http.DefaultClient
			response, err := client.Get("http://127.0.0.1:" + port + "/list")
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(BeSuccessful())

			b, err := ioutil.ReadFile("fixtures/repo-index-response.json")
			Expect(err).NotTo(HaveOccurred())

			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal(string(b)))
		})
	})
})

func BeSuccessful() types.GomegaMatcher {
	return &SuccessfulHTTPResponseMatcher{}
}

type SuccessfulHTTPResponseMatcher struct{}

func (matcher *SuccessfulHTTPResponseMatcher) Match(actual interface{}) (success bool, err error) {
	response, ok := actual.(*http.Response)
	if !ok {
		return false, fmt.Errorf("SuccessfulHTTPResponseMatcher matcher expects an http.Response")
	}

	return (response.StatusCode >= 200) && (response.StatusCode < 400), nil
}

func (matcher *SuccessfulHTTPResponseMatcher) FailureMessage(actual interface{}) (message string) {
	response := actual.(*http.Response)

	return fmt.Sprintf("Expected Status Code\n\t%d\nto be successful (2XX or 3XX)", response.StatusCode)
}

func (matcher *SuccessfulHTTPResponseMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	response := actual.(*http.Response)

	return fmt.Sprintf("Expected Status Code\n\t%d\nto not be successful (1XX, 4XX, 5XX)", response.StatusCode)
}
