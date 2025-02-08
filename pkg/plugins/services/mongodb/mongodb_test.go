package mongodb

import (
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/praetorian-inc/fingerprintx/pkg/test"
)

func TestMongoDb(t *testing.T) {
	testcases := []test.Testcase{
		{
			Description: "mongodb",
			Port:        3306,
			Protocol:    plugins.TCP,
			Expected: func(res *plugins.Service) bool {
				return res != nil
			},
			RunConfig: dockertest.RunOptions{
				Repository: "mongo",
				Tag:        "5.7.39",
				Env: []string{
					"MongoDB_ROOT_PASSWORD=my-secret-pw",
				},
			},
		},
	}

	p := &MongoDBPlugin{}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			err := test.RunTest(t, tc, p)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
