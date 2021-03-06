package api_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rsamban/om/api"
	"github.com/rsamban/om/api/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StemcellService", func() {
	Describe("Upload", func() {
		var (
			client *fakes.HttpClient
			bar    *fakes.Progress
		)

		BeforeEach(func() {
			client = &fakes.HttpClient{}
			bar = &fakes.Progress{}
		})

		It("makes a request to upload the stemcell to the OpsManager", func() {
			client.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			}, nil)

			bar.NewBarReaderReturns(strings.NewReader("some other content"))
			service := api.NewUploadStemcellService(client, bar)

			output, err := service.Upload(api.StemcellUploadInput{
				ContentLength: 10,
				Stemcell:      strings.NewReader("some content"),
				ContentType:   "some content-type",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(api.StemcellUploadOutput{}))

			request := client.DoArgsForCall(0)
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.Path).To(Equal("/api/v0/stemcells"))
			Expect(request.ContentLength).To(Equal(int64(10)))
			Expect(request.Header.Get("Content-Type")).To(Equal("some content-type"))

			body, err := ioutil.ReadAll(request.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(Equal("some other content"))

			newReaderContent, err := ioutil.ReadAll(bar.NewBarReaderArgsForCall(0))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(newReaderContent)).To(Equal("some content"))
			Expect(bar.SetTotalArgsForCall(0)).To(BeNumerically("==", 10))
			Expect(bar.KickoffCallCount()).To(Equal(1))
			Expect(bar.EndCallCount()).To(Equal(1))
		})

		Context("when an error occurs", func() {
			Context("when the client errors before the request", func() {
				It("returns an error", func() {
					client.DoReturns(&http.Response{}, errors.New("some client error"))
					service := api.NewUploadStemcellService(client, bar)

					_, err := service.Upload(api.StemcellUploadInput{})
					Expect(err).To(MatchError("could not make api request to stemcells endpoint: some client error"))
				})
			})

			Context("when the api returns a non-200 status code", func() {
				It("returns an error", func() {
					client.DoReturns(&http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(strings.NewReader("{}")),
					}, nil)
					service := api.NewUploadStemcellService(client, bar)

					_, err := service.Upload(api.StemcellUploadInput{})
					Expect(err).To(MatchError(ContainSubstring("request failed: unexpected response")))
				})
			})
		})
	})
})
