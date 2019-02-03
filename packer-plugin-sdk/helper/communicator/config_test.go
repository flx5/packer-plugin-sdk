package communicator

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/hashicorp/packer/template/interpolate"
	"github.com/masterzen/winrm"
)

func testConfig() *Config {
	return &Config{
		SSHUsername: "root",
	}
}

func TestConfigType(t *testing.T) {
	c := testConfig()
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.Type != "ssh" {
		t.Fatalf("bad: %#v", c)
	}
}

func TestConfig_none(t *testing.T) {
	c := &Config{Type: "none"}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}
}

func TestConfig_badtype(t *testing.T) {
	c := &Config{Type: "foo"}
	if err := c.Prepare(testContext(t)); len(err) != 1 {
		t.Fatalf("bad: %#v", err)
	}
}

func TestConfig_winrm_noport(t *testing.T) {
	c := &Config{
		Type:      "winrm",
		WinRMUser: "admin",
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5985 {
		t.Fatalf("WinRMPort doesn't match default port 5985 when SSL is not enabled and no port is specified.")
	}

}

func TestConfig_winrm_noport_ssl(t *testing.T) {
	c := &Config{
		Type:        "winrm",
		WinRMUser:   "admin",
		WinRMUseSSL: true,
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5986 {
		t.Fatalf("WinRMPort doesn't match default port 5986 when SSL is enabled and no port is specified.")
	}

}

func TestConfig_winrm_port(t *testing.T) {
	c := &Config{
		Type:      "winrm",
		WinRMUser: "admin",
		WinRMPort: 5509,
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5509 {
		t.Fatalf("WinRMPort doesn't match custom port 5509 when SSL is not enabled.")
	}

}

func TestConfig_winrm_port_ssl(t *testing.T) {
	c := &Config{
		Type:        "winrm",
		WinRMUser:   "admin",
		WinRMPort:   5510,
		WinRMUseSSL: true,
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5510 {
		t.Fatalf("WinRMPort doesn't match custom port 5510 when SSL is enabled.")
	}

}

func TestConfig_winrm_use_ntlm(t *testing.T) {
	c := &Config{
		Type:         "winrm",
		WinRMUser:    "admin",
		WinRMUseNTLM: true,
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMTransportDecorator == nil {
		t.Fatalf("WinRMTransportDecorator not set.")
	}

	expected := &winrm.ClientNTLM{}
	actual := c.WinRMTransportDecorator()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("WinRMTransportDecorator isn't ClientNTLM.")
	}

}

func TestConfig_winrm(t *testing.T) {
	c := &Config{
		Type:      "winrm",
		WinRMUser: "admin",
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}
}

func TestConfig_SSHPublicKeyUrlEncoded(t *testing.T) {
	c := &Config{
		SSHPublicKey: []byte("ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBADulbdHCjXhsH8wGtyLhhi3qVvX6M0tGgtousr/DzArwf2KX0L2Zm1OZfqMWFCrSVD743OFY60YL5CGsN9/PVQP7gApll5yTWyaQJu8lReptR5TMnUDn0u3mJN/QRT5Zs8qS5J5Q3WhXwaMF96kSuu+MwXrBnl8sK+bwxOKQtlKJXowcw==\n"),
	}

	encoded := c.SSHPublicKeyUrlEncoded()

	decoded, err := url.PathUnescape(encoded)
	if err != nil {
		t.Fatal(err.Error())
	}

	if decoded != string(c.SSHPublicKey) {
		t.Fatal("resulting public key does not match original public key")
	}
}

func testContext(t *testing.T) *interpolate.Context {
	return nil
}
