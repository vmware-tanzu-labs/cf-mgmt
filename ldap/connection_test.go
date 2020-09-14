package ldap_test

import (
	"errors"

	l "github.com/go-ldap/ldap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap/fakes"
)

var _ = Describe("NewRefreshableConnection", func() {
	It("wraps a connection in a refresher", func() {
		_, err := ldap.NewRefreshableConnection(func() (ldap.Connection, error) {
			return &fakes.FakeConnection{}, nil
		})

		Expect(err).ShouldNot(HaveOccurred())
	})

	When("the connection cannot be created", func() {
		It("passes the error through", func() {
			_, err := ldap.NewRefreshableConnection(func() (ldap.Connection, error) {
				return nil, errors.New("some error")
			})

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError("some error"))
		})
	})
})

var _ = Describe("RefreshableConnection_Search", func() {
	var (
		rc                          *ldap.RefreshableConnection
		err                         error
		createConnectionCallCounter = new(int)
	)

	BeforeEach(func() {
		*createConnectionCallCounter = 0
	})

	newRCWithClosing := func(b bool) (*ldap.RefreshableConnection, error) {
		return ldap.NewRefreshableConnection(
			withCallCounter(createConnectionCallCounter, func() (ldap.Connection, error) {
				fakeConn := &fakes.FakeConnection{}
				fakeConn.IsClosingReturns(b)
				return fakeConn, nil
			}),
		)
	}

	When("the connection is not closing", func() {
		It("just returns the search results", func() {
			Expect(*createConnectionCallCounter).Should(Equal(0))

			rc, err = newRCWithClosing(false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(1))

			_, err = rc.Search(&l.SearchRequest{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(1))
		})
	})

	When("the connection is closing", func() {
		It("tries to refresh the connection", func() {
			Expect(*createConnectionCallCounter).Should(Equal(0))

			rc, err = newRCWithClosing(true)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(1))

			_, err = rc.Search(&l.SearchRequest{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(2))
		})
	})

	When("refresh connection returns an error", func() {
		It("passes the error through", func() {
			const errorMsg = "error refreshing"

			throwError := false
			rc, err = ldap.NewRefreshableConnection(func() (ldap.Connection, error) {
				if throwError {
					return nil, errors.New(errorMsg)
				}

				fc := &fakes.FakeConnection{}
				fc.IsClosingReturns(true)

				return fc, nil
			})

			throwError = true
			_, err = rc.Search(&l.SearchRequest{})
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError(errorMsg))
		})
	})
})

var _ = Describe("RefreshableConnection_RefreshConnection", func() {})

func withCallCounter(callCounter *int, createConnection func() (ldap.Connection, error)) func() (ldap.Connection, error) {
	return func() (ldap.Connection, error) {
		(*callCounter)++
		return createConnection()
	}
}
