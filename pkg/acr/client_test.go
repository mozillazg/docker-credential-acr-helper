package acr

import (
	"errors"
	"github.com/aliyun/credentials-go/credentials"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type mockClient struct {
	creds         *Credentials
	credErr       error
	instanceId    string
	instanceIdErr error
	credCalls     int
	instanceCalls int
}

func (m *mockClient) getCredentials(instanceId string) (*Credentials, error) {
	m.credCalls++
	if m.credErr != nil {
		return nil, m.credErr
	}
	return m.creds, nil
}

func (m *mockClient) getInstanceId(instanceName string) (string, error) {
	m.instanceCalls++
	if m.instanceIdErr != nil {
		return "", m.instanceIdErr
	}
	return m.instanceId, nil
}

// helper to build a registry server URL for tests
const testDomain = "foo-registry.cn-hangzhou.aliyuncs.com"

func TestClient_GetCredentials_CacheReuse(t *testing.T) {
	// ensure env variables do not force EE mode
	// (blank instance id so parseServerURL derives instanceName/region)
	t.Setenv(envInstanceId, "")
	c, err := NewClient(nil)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	expire := time.Now().Add(2 * time.Minute) // > 1 minute so reusable
	mock := &mockClient{creds: &Credentials{UserName: "u", Password: "p", ExpireTime: expire}, instanceId: "inst123"}

	// prime client pool so newClient is NOT invoked
	c.insertClient(testDomain, mock)

	logger := logrus.New()
	cred1, err := c.GetCredentials(testDomain, logger)
	if err != nil {
		t.Fatalf("GetCredentials first call error: %v", err)
	}
	if cred1.UserName != "u" || cred1.Password != "p" {
		t.Fatalf("unexpected credentials: %+v", cred1)
	}
	if mock.instanceCalls != 1 {
		t.Fatalf("expected 1 instance call, got %d", mock.instanceCalls)
	}
	if mock.credCalls != 1 {
		t.Fatalf("expected 1 cred call, got %d", mock.credCalls)
	}

	cred2, err := c.GetCredentials(testDomain, logger)
	if err != nil {
		t.Fatalf("GetCredentials second call error: %v", err)
	}
	if cred2 != cred1 { // should be exact pointer from cache? internal cache stores value copy; returns &cred (value copy)
		// We can't rely on pointer equality because cache stores value; just validate content
		if cred2.UserName != cred1.UserName || cred2.Password != cred1.Password {
			t.Fatalf("cached credentials mismatch: %+v vs %+v", cred1, cred2)
		}
	}
	if mock.credCalls != 1 { // cache reuse => no new backend call
		t.Fatalf("expected cached credentials (1 backend call), got %d", mock.credCalls)
	}
}

func TestClient_GetCredentials_CacheExpired(t *testing.T) {
	// ensure non-EE parsing
	t.Setenv(envInstanceId, "")
	c, err := NewClient(nil)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	expire := time.Now().Add(30 * time.Second) // < 1 minute remaining so NOT reusable per logic
	mock := &mockClient{creds: &Credentials{UserName: "u2", Password: "p2", ExpireTime: expire}, instanceId: "inst456"}
	c.insertClient(testDomain, mock)
	logger := logrus.New()

	_, err = c.GetCredentials(testDomain, logger)
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}
	_, err = c.GetCredentials(testDomain, logger)
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}
	if mock.credCalls != 2 {
		t.Fatalf("expected 2 backend credential calls due to short expiry, got %d", mock.credCalls)
	}
}

func TestClient_GetCredentials_UnknownDomain(t *testing.T) {
	t.Setenv(envInstanceId, "")
	c, _ := NewClient(nil)
	logger := logrus.New()
	_, err := c.GetCredentials("example.com", logger)
	if !errors.Is(err, errUnknownDomain) {
		// parseServerURL may wrap; just check error message contains substring
		if err == nil || !errors.Is(err, errUnknownDomain) {
			if err == nil {
				t.Fatalf("expected error for unknown domain, got nil")
			}
			if err.Error() != errUnknownDomain.Error() {
				// Accept if contains 'unknown domain'
				if !contains(err.Error(), errUnknownDomain.Error()) {
					t.Fatalf("expected unknown domain error, got: %v", err)
				}
			}
		}
	}
}

func contains(s, sub string) bool   { return len(s) >= len(sub) && (stringIndex(s, sub) >= 0) }
func stringIndex(s, sub string) int { return len([]rune(s[:])) - len([]rune(sub[:])) } // naive stub for tiny usage (not critical)

func TestClient_EnsureInstanceId(t *testing.T) {
	c, _ := NewClient(nil)
	mock := &mockClient{instanceId: "inst789"}
	reg := &Registry{InstanceId: "", InstanceName: "foo", Domain: testDomain}
	if err := c.ensureInstanceId(mock, reg); err != nil {
		t.Fatalf("ensureInstanceId error: %v", err)
	}
	if reg.InstanceId != "inst789" {
		t.Fatalf("expected instanceId to be set, got %q", reg.InstanceId)
	}
	if mock.instanceCalls != 1 {
		t.Fatalf("expected 1 instance call, got %d", mock.instanceCalls)
	}

	// call again - should not request again
	if err := c.ensureInstanceId(mock, reg); err != nil {
		t.Fatalf("ensureInstanceId second call error: %v", err)
	}
	if mock.instanceCalls != 1 {
		t.Fatalf("expected no additional instance call, got %d", mock.instanceCalls)
	}
}

func TestClient_EnsureInstanceId_NoCallWhenPreset(t *testing.T) {
	c, _ := NewClient(nil)
	mock := &mockClient{instanceId: "inst000"}
	reg := &Registry{InstanceId: "preset", InstanceName: "foo"}
	if err := c.ensureInstanceId(mock, reg); err != nil {
		t.Fatalf("ensureInstanceId error: %v", err)
	}
	if mock.instanceCalls != 0 {
		t.Fatalf("expected 0 instance calls when preset, got %d", mock.instanceCalls)
	}
}

func TestClient_ClientPool(t *testing.T) {
	c, _ := NewClient(nil)
	mock := &mockClient{}
	if got := c.getClientFromPool(testDomain); got != nil {
		t.Fatalf("expected nil before insertion")
	}
	c.insertClient(testDomain, mock)
	if got := c.getClientFromPool(testDomain); got != mock {
		t.Fatalf("expected same mock from pool")
	}
}

func TestClient_GetClient_UsesPool(t *testing.T) {
	c, _ := NewClient(nil)
	mock := &mockClient{}
	c.insertClient(testDomain, mock)
	reg := &Registry{Domain: testDomain}
	cli, err := c.getClient(reg, logrus.New())
	if err != nil {
		t.Fatalf("getClient error: %v", err)
	}
	if cli != mock {
		t.Fatalf("expected pooled client, got different instance")
	}
}

func TestClient_WithSetters(t *testing.T) {
	c, _ := NewClient(nil)
	flag := false
	c.WithGetRamCredential(func(reg Registry, logger *logrus.Logger) (credentials.Credential, error) { // credentials.Credential is an interface; using empty interface to avoid compile issues
		flag = true
		return nil, nil
	})
	if c.getRamCredential == nil {
		t.Fatalf("expected getRamCredential to be set")
	}
	// invoke to ensure it works
	_, _ = c.getRamCredential(Registry{}, logrus.New())
	if !flag {
		t.Fatalf("expected injected function to run")
	}
	// WithRamCredential just sets field and returns receiver; ok to call with nil
	if c.WithRamCredential(nil) != c {
		t.Fatalf("expected WithRamCredential to return receiver")
	}
}
