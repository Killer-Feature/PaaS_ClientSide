package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/strvals"
)

var (
	prometheusArgs = map[string]string{
		// comma seperated values to set
		"set": "global.storageClass=local-storage,primary.persistence.size=4Gi,auth.postgresPassword=pgpass",
	}
	grafanaArgs = map[string]string{
		// comma seperated values to set
		"set": "admin.password=admin",
	}
)

type HelmInstaller struct {
	namespace string
	repoUrl   string
	repoName  string

	config *genericclioptions.ConfigFlags

	l        *zap.Logger
	settings *cli.EnvSettings
}

func NewHelmInstaller(namespace, repoUrl, repoName string, logger *zap.Logger) (*HelmInstaller, error) {
	hi := &HelmInstaller{
		namespace: namespace,
		repoUrl:   repoUrl,
		repoName:  repoName,
		l:         logger,
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	hi.settings = cli.New()

	hi.config = kube.GetConfig("./config", "", namespace)

	// Add helm repo
	_ = hi.RepoAdd(hi.repoName, hi.repoUrl)

	// Add meatllb repo
	_ = hi.RepoAdd("metallb", "https://metallb.github.io/metallb")

	// Update charts from the helm repo
	_ = hi.RepoUpdate()

	return hi, nil
}

func (hi *HelmInstaller) Install(releaseName string, rType internal.ResourceType) error {
	// Install charts
	switch rType {
	case internal.Postgres:
		return hi.InstallChart(releaseName, hi.repoName, "postgresql", prometheusArgs)
	case internal.Redis:
		return hi.InstallChart(releaseName, hi.repoName, "redis", nil)
	case internal.Prometheus:
		return hi.InstallChart(releaseName, hi.repoName, "kube-prometheus", nil)
	case internal.Grafana:
		return hi.InstallChart(releaseName, hi.repoName, "grafana", grafanaArgs)
	case internal.NginxIngressController:
		return hi.InstallChart(releaseName, hi.repoName, "nginx-ingress-controller", nil)
	case internal.MetalLB:
		return hi.InstallChart(releaseName, "metallb", "metallb", nil)
	}
	return errors.New("resource type not provided")
}

// RepoAdd adds repo with given name and url
func (hi *HelmInstaller) RepoAdd(name, url string) error {
	repoFile := hi.settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		hi.l.Error("failed creating directory for helm repositories")
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := os.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {
		hi.l.Warn("repository name already exists", zap.String("repo_name", name))
		return nil
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(hi.settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	hi.l.Info("new repositories has been added", zap.String("repo_name", name))
	return nil
}

// RepoUpdate updates charts for all helm repos
func (hi *HelmInstaller) RepoUpdate() error {
	repoFile := hi.settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(hi.settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	hi.l.Info("Grabbing the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				hi.l.Error("Unable to get an update from the chart repository", zap.String("repo_name", re.Config.Name), zap.String("repo_url", re.Config.URL), zap.Error(err))
			} else {
				hi.l.Info("Successfully got an update from the chart repository", zap.String("repo_name", re.Config.Name))
			}
		}(re)
	}
	wg.Wait()

	hi.l.Info("Helm update complete")
	return nil
}

// InstallChart installs chart
func (hi *HelmInstaller) InstallChart(name, repo, chart string, args map[string]string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(hi.config, hi.settings.Namespace(), os.Getenv("HELM_DRIVER"), hi.debug); err != nil {
		return err
	}
	client := action.NewInstall(actionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = name
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), hi.settings)
	if err != nil {
		return err
	}

	hi.l.Debug("chart path", zap.String("chart_path", cp))

	p := getter.All(hi.settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}

	// Add args
	if err := strvals.ParseInto(args["set"], vals); err != nil {
		return errors.Wrap(err, "failed parsing --set data")
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: hi.settings.RepositoryConfig,
					RepositoryCache:  hi.settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	client.Namespace = hi.settings.Namespace()
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		return err
	}
	//fmt.Println(release.Manifest)
	return nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func (hi *HelmInstaller) debug(format string, v ...interface{}) {
	hi.l.Debug(fmt.Sprintf(format, v...))
}

// UninstallChart uninstalls chart
func (hi *HelmInstaller) UninstallChart(name string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(hi.config, hi.settings.Namespace(), os.Getenv("HELM_DRIVER"), hi.debug); err != nil {
		return err
	}
	client := action.NewUninstall(actionConfig)

	_, err := client.Run(name)
	//fmt.Println(resp)
	return err
}

type Resourse struct {
	Name          string
	Status        string
	FirstDeployed string
	LastDeployed  string
	AppVersion    string
	ApiVersion    string
	Description   string
	ChartVersion  string
	Type          string
	ChartURL      string
}

func (hi *HelmInstaller) GetResourcesList() ([]Resourse, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(hi.config, hi.settings.Namespace(), os.Getenv("HELM_DRIVER"), hi.debug); err != nil {
		return nil, err
	}
	client := action.NewList(actionConfig)

	resources, err := client.Run()
	if err != nil {
		hi.l.Error("failed listing helm resources", zap.String("err", err.Error()))
		return nil, err
	}

	resourceList := make([]Resourse, 0, len(resources))
	for _, res := range resources {
		resourceList = append(resourceList, Resourse{
			Name:          res.Name,
			Status:        res.Info.Status.String(),
			FirstDeployed: res.Info.FirstDeployed.String(),
			LastDeployed:  res.Info.LastDeployed.String(),
			AppVersion:    res.Chart.Metadata.AppVersion,
			Description:   res.Chart.Metadata.Description,
			ChartVersion:  res.Chart.Metadata.Version,
			ApiVersion:    res.Chart.Metadata.APIVersion,
			Type:          res.Chart.Metadata.Name,
			ChartURL:      res.Chart.Metadata.Home,
		})
	}

	return resourceList, err
}
