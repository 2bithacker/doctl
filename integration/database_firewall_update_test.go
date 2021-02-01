package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os/exec"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
)

var _ = suite("database/firewalls", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/databases/d168d635-1c88-4616-b9b4-793b7c573927/firewall":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method == http.MethodPut {
					reqBody, err := ioutil.ReadAll(req.Body)
					expect.NoError(err)

					expect.JSONEq(databasesUpdateFirewallUpdateRequest, string(reqBody))

					w.Write([]byte(databasesUpdateFirewallRuleResponse))
				} else if req.Method == http.MethodGet {
					w.Write([]byte(databasesUpdateFirewallRuleResponse))
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))
	})

	when("command is update", func() {
		it("update a database cluster's firewall rules", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"databases",
				"firewalls",
				"update",
				"d168d635-1c88-4616-b9b4-793b7c573927",
				"--rule", "ip_addr:192.168.1.1",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received error output: %s", output))
<<<<<<< HEAD
			expect.Equal(strings.TrimSpace(databasesUpdateFirewallRuleOutput), strings.TrimSpace(string(output)))
=======
>>>>>>> 7b764d45b13864c3c93a26e57f98a6fa77657cb1

			expected := strings.TrimSpace(databasesUpdateFirewallRuleOutput)
			actual := strings.TrimSpace(string(output))

			if expected != actual {
				t.Errorf("expected\n\n%s\n\nbut got\n\n%s\n\n", expected, actual)
			}
		})
	})

})

const (
	databasesUpdateFirewallUpdateRequest = `{"rules": [{"type": "ip_addr","value": "192.168.1.1"}]}`

	databasesUpdateFirewallRuleOutput = `
UUID                                    ClusterUUID                             Type       Value          Created At
82ebbbd4-437c-4e11-bfd2-644ccb555de0    d168d635-1c88-4616-b9b4-793b7c573927    ip_addr    192.168.1.1    2021-01-29 19:59:35 +0000 UTC`

	databasesUpdateFirewallRuleResponse = `{
		"rules":[
		   {
			  "uuid":"82ebbbd4-437c-4e11-bfd2-644ccb555de0",
			  "cluster_uuid":"d168d635-1c88-4616-b9b4-793b7c573927",
			  "type":"ip_addr",
			  "value":"192.168.1.1",
			  "created_at":"2021-01-29T19:59:35Z"
		   }
		]
	 }`
)
