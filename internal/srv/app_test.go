package srv

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestNewHelmValues(t *testing.T) {
	type testCase struct {
		name        string
		valuesPath  string
		overrides   []valueSet
		expectError bool
	}

	testCases := []testCase{
		{
			name:        "valid values path",
			expectError: false,
			valuesPath:  "/tmp/values.yaml",
			overrides:   nil,
		},
		{
			name:        "valid overrides",
			expectError: false,
			valuesPath:  "/tmp/values.yaml",
			overrides: []valueSet{
				{
					helmKey: "hello",
					value:   "world",
				},
			},
		},
		{
			name:        "missing values path",
			expectError: true,
			valuesPath:  "",
			overrides:   nil,
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			srv := Server{
				Logger:     setupTestLogger(t, tcase.name),
				ValuesPath: tcase.valuesPath,
			}
			values, err := srv.newHelmValues(tcase.overrides)
			fmt.Println(values)
			if tcase.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, values)
			}
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	type testCase struct {
		name         string
		appNamespace string
		expectError  bool
		kubeclient   *rest.Config
	}

	env := envtest.Environment{}

	cfg, err := env.Start()
	if err != nil {
		panic(err)
	}

	testCases := []testCase{
		{
			name:         "valid yaml",
			expectError:  false,
			appNamespace: "flintlock",
			kubeclient:   cfg,
		},
		{
			name:         "invalid namespace",
			expectError:  true,
			appNamespace: "DarkwingDuck",
			kubeclient:   cfg,
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			srv := Server{
				Context:    context.TODO(),
				Logger:     setupTestLogger(t, tcase.name),
				KubeClient: tcase.kubeclient,
			}

			err := srv.CreateNamespace(tcase.appNamespace)

			if tcase.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	err = env.Stop()
	if err != nil {
		panic(err)
	}
}

func TestCreateApp(t *testing.T) {
	type testCase struct {
		name         string
		appNamespace string
		appName      string
		expectError  bool
		chartPath    string
	}

	env := envtest.Environment{}

	cfg, err := env.Start()
	if err != nil {
		panic(err)
	}

	testCases := []testCase{
		{
			name:         "valid yaml",
			expectError:  false,
			appNamespace: uuid.New().String(),
			appName:      uuid.New().String(),
			chartPath:    "/tmp/chart.tgz",
		},
		{
			name:         "invalid namespace",
			expectError:  true,
			appNamespace: "DarkwingDuck",
			appName:      uuid.New().String(),
		},
		{
			name:         "invalid chart",
			expectError:  true,
			appNamespace: uuid.New().String(),
			appName:      uuid.New().String(),
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			srv := Server{
				Context:    context.TODO(),
				Logger:     setupTestLogger(t, tcase.name),
				KubeClient: cfg,
				ChartPath:  tcase.chartPath,
				ValuesPath: "/tmp/values.yaml",
			}

			_ = srv.CreateNamespace(tcase.appNamespace)
			err = srv.CreateApp(tcase.appName, tcase.appNamespace, nil)

			if tcase.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	err = env.Stop()
	if err != nil {
		panic(err)
	}
}
