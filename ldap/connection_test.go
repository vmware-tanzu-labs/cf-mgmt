package ldap_test

import (
	"crypto/tls"
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

var _ = Describe("RefreshableConnection_RefreshConnection", func() {
	var (
		rc                          *ldap.RefreshableConnection
		err                         error
		createConnectionCallCounter = new(int)
	)

	BeforeEach(func() {
		*createConnectionCallCounter = 0
	})

	newRC := func() (*ldap.RefreshableConnection, error) {
		return ldap.NewRefreshableConnection(
			withCallCounter(createConnectionCallCounter, func() (ldap.Connection, error) {
				fakeConn := &fakes.FakeConnection{}
				return fakeConn, nil
			}),
		)
	}

	When("refreshConnection does not return an error", func() {
		It("creates a new connection", func() {
			Expect(*createConnectionCallCounter).Should(Equal(0))

			rc, err = newRC()
			connBeforeRefreshConnection := rc.Connection
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(1))

			err = rc.RefreshConnection()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*createConnectionCallCounter).Should(Equal(2))
			Expect(rc.Connection).ShouldNot(BeIdenticalTo(connBeforeRefreshConnection))
		})
	})

	When("refreshConnection returns an error", func() {
		It("passes the error through", func() {
			const errorMsg = "error refreshing"

			throwError := false
			rc, err = ldap.NewRefreshableConnection(func() (ldap.Connection, error) {
				if throwError {
					return nil, errors.New(errorMsg)
				}

				fakeConn := &fakes.FakeConnection{}
				return fakeConn, nil
			})

			throwError = true
			err = rc.RefreshConnection()
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError(errorMsg))
		})
	})
})

func withCallCounter(callCounter *int, createConnection func() (ldap.Connection, error)) func() (ldap.Connection, error) {
	return func() (ldap.Connection, error) {
		(*callCounter)++
		return createConnection()
	}
}

var _ = Describe("MapTLSVersion", func() {
	When("when 1.0", func() {
		It("returns VersionTLS10", func() {
			val, err := ldap.MapTLSVersion("1.0")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(uint16(tls.VersionTLS10)))
		})
	})
	When("when 1.1", func() {
		It("returns VersionTLS11", func() {
			val, err := ldap.MapTLSVersion("1.1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(uint16(tls.VersionTLS11)))
		})
	})
	When("when 1.2", func() {
		It("returns VersionTLS12", func() {
			val, err := ldap.MapTLSVersion("1.2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(uint16(tls.VersionTLS12)))
		})
	})
	When("when 1.3", func() {
		It("returns VersionTLS13", func() {
			val, err := ldap.MapTLSVersion("1.3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(uint16(tls.VersionTLS13)))
		})
	})

	When("when unknown value", func() {
		It("returns error", func() {
			_, err := ldap.MapTLSVersion("1.3.1")
			Expect(err).Should(HaveOccurred())
		})
	})
})
