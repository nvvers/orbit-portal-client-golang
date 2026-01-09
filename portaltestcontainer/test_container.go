package portaltestcontainer

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nvvers/orbit-portal-client-golang/portal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const imageName = "cr.nv-online.dev/nv/orbit/portal"

type ContainerOption func(*TestContainer) error

type TestContainer struct {
	container    testcontainers.Container
	imageTag     string
	internalPort uint16
	token        string
}

func New(opts ...ContainerOption) (*TestContainer, error) {
	c := &TestContainer{
		imageTag:     "latest",
		internalPort: 8080,
		token:        "secret",
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithImageTag(tag string) ContainerOption {
	return func(p *TestContainer) error {
		p.imageTag = tag
		return nil
	}
}

func WithToken(token string) ContainerOption {
	return func(p *TestContainer) error {
		p.token = token
		return nil
	}
}

func WithPort(port uint16) ContainerOption {
	return func(p *TestContainer) error {
		p.internalPort = port
		return nil
	}
}

func (p *TestContainer) Start(ctx context.Context) error {
	config := fmt.Sprintf(`logFormat: "text"
apiEndpoint: ":%d"

ingressOnly: false
attachmentFolder: "/orbit/attachments"

authentication:
  allowAnonymous: false
  clients:
    - clientID: "testclient"
      token: "%s"
      rights:
        ingestEvents: true
        retrieveEvents: true
        managePools: true
        manageSchemas: true
        showInsights: true
      allowedPools: ['*']

eventStore: "orbstore"
orbStore:
  storeDir: "/orbit/orbstore-data"
`, p.internalPort, p.token)

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", imageName, p.imageTag),
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", p.internalPort)},
		Env: map[string]string{
			"ORBIT_PORTAL_CONFIG_PATH": "/orbit/portal-config.yaml",
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(config),
				ContainerFilePath: "/orbit/portal-config.yaml",
				FileMode:          0o644,
			},
		},
		WaitingFor: wait.ForHTTP("/api/v1/ping").
			WithPort(fmt.Sprintf("%d/tcp", p.internalPort)).
			WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start portal container: %w", err)
	}

	p.container = container
	return nil
}

func (p *TestContainer) GetHost(ctx context.Context) (string, error) {
	if p.container == nil {
		return "", errors.New("container must be running")
	}

	host, err := p.container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	return host, nil
}

func (p *TestContainer) GetMappedPort(ctx context.Context) (uint16, error) {
	if p.container == nil {
		return 0, errors.New("container must be running")
	}

	mappedPort, err := p.container.MappedPort(ctx, fmt.Sprintf("%d/tcp", p.internalPort))
	if err != nil {
		return 0, fmt.Errorf("failed to get mapped port: %w", err)
	}

	mapped, err := strconv.ParseUint(mappedPort.Port(), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to parse mapped port: %w", err)
	}

	return uint16(mapped), nil
}

func (p *TestContainer) GetBaseURL(ctx context.Context) (*url.URL, error) {
	host, err := p.GetHost(ctx)
	if err != nil {
		return nil, err
	}

	port, err := p.GetMappedPort(ctx)
	if err != nil {
		return nil, err
	}

	return url.Parse(fmt.Sprintf("http://%s:%d", host, port))
}

func (p *TestContainer) Token() string {
	return p.token
}

func (p *TestContainer) IsRunning(ctx context.Context) bool {
	if p.container == nil {
		return false
	}

	state, err := p.container.State(ctx)
	if err != nil {
		return false
	}

	return state.Running
}

func (p *TestContainer) Stop(ctx context.Context) error {
	if p.container == nil {
		return nil
	}

	if err := p.container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	p.container = nil
	return nil
}

func (p *TestContainer) GetClient(ctx context.Context) (*portal.Client, error) {
	if !p.IsRunning(ctx) {
		return nil, errors.New("container must be running to get a client")
	}

	baseURL, err := p.GetBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	client, err := portal.NewClient(baseURL.String(), portal.WithToken(p.Token()))
	if err != nil {
		return nil, fmt.Errorf("failed to create portal client: %w", err)
	}

	return client, nil
}
