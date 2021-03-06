// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"istio.io/istio/istioctl/pkg/clioptions"
	"istio.io/istio/istioctl/pkg/kubernetes"
	"istio.io/istio/istioctl/pkg/util/handlers"

	"istio.io/pkg/log"
)

var (
	listenPort   = 0
	controlZport = 0

	bindAddress = ""

	// label selector
	labelSelector = ""
)

// port-forward to Istio System Prometheus; open browser
func promDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "prometheus",
		Short:   "Open Prometheus web UI",
		Long:    `Open Istio's Prometheus dashboard`,
		Example: `istioctl dashboard prometheus`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			pl, err := client.PodsForSelector(istioNamespace, "app=prometheus")
			if err != nil {
				return fmt.Errorf("not able to locate Prometheus pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Prometheus pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, istioNamespace, "Prometheus",
				"http://localhost:%d", bindAddress, 9090, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to Istio System Grafana; open browser
func grafanaDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "grafana",
		Short:   "Open Grafana web UI",
		Long:    `Open Istio's Grafana dashboard`,
		Example: `istioctl dashboard grafana`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			pl, err := client.PodsForSelector(istioNamespace, "app=grafana")
			if err != nil {
				return fmt.Errorf("not able to locate Grafana pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Grafana pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, istioNamespace, "Grafana",
				"http://localhost:%d", bindAddress, 3000, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to Istio System Kiali; open browser
func kialiDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "kiali",
		Short:   "Open Kiali web UI",
		Long:    `Open Istio's Kiali dashboard`,
		Example: `istioctl dashboard kiali`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			pl, err := client.PodsForSelector(istioNamespace, "app=kiali")
			if err != nil {
				return fmt.Errorf("not able to locate Kiali pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Kiali pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, istioNamespace, "Kiali",
				"http://localhost:%d/kiali", bindAddress, 20001, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to Istio System Jaeger; open browser
func jaegerDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "jaeger",
		Short:   "Open Jaeger web UI",
		Long:    `Open Istio's Jaeger dashboard`,
		Example: `istioctl dashboard jaeger`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			pl, err := client.PodsForSelector(istioNamespace, "app=jaeger")
			if err != nil {
				return fmt.Errorf("not able to locate Jaeger pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Jaeger pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, istioNamespace, "Jaeger",
				"http://localhost:%d", bindAddress, 16686, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to Istio System Zipkin; open browser
func zipkinDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "zipkin",
		Short:   "Open Zipkin web UI",
		Long:    `Open Istio's Zipkin dashboard`,
		Example: `istioctl dashboard zipkin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			pl, err := client.PodsForSelector(istioNamespace, "app=zipkin")
			if err != nil {
				return fmt.Errorf("not able to locate Zipkin pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Zipkin pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, istioNamespace, "Zipkin",
				"http://localhost:%d", bindAddress, 9411, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to sidecar Envoy admin port; open browser
func envoyDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "envoy <pod-name[.namespace]>",
		Short:   "Open Envoy admin web UI",
		Long:    `Open the Envoy admin dashboard for a sidecar`,
		Example: `istioctl dashboard envoy productpage-123-456.default`,
		RunE: func(c *cobra.Command, args []string) error {
			if labelSelector == "" && len(args) < 1 {
				c.Println(c.UsageString())
				return fmt.Errorf("specify a pod or --selector")
			}

			if labelSelector != "" && len(args) > 0 {
				c.Println(c.UsageString())
				return fmt.Errorf("name cannot be provided when a selector is specified")
			}

			client, err := envoyClientFactory(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			var podName, ns string
			if labelSelector != "" {
				pl, err := client.PodsForSelector(handlers.HandleNamespace(namespace, defaultNamespace), labelSelector)
				if err != nil {
					return fmt.Errorf("not able to locate pod with selector %s: %v", labelSelector, err)
				}

				if len(pl.Items) < 1 {
					return errors.New("no pods found")
				}

				if len(pl.Items) > 1 {
					log.Warnf("more than 1 pods fits selector: %s; will use pod: %s", labelSelector, pl.Items[0].Name)
				}

				// only use the first pod in the list
				podName = pl.Items[0].Name
				ns = pl.Items[0].Namespace
			} else {
				podName, ns = handlers.InferPodInfo(args[0], handlers.HandleNamespace(namespace, defaultNamespace))
			}

			return portForward(podName, ns, fmt.Sprintf("Envoy sidecar %s", podName),
				"http://localhost:%d", bindAddress, 15000, client, c.OutOrStdout())
		},
	}

	return cmd
}

// port-forward to sidecar ControlZ port; open browser
func controlZDashCmd() *cobra.Command {
	var opts clioptions.ControlPlaneOptions
	cmd := &cobra.Command{
		Use:     "controlz <pod-name[.namespace]>",
		Short:   "Open ControlZ web UI",
		Long:    `Open the ControlZ web UI for a pod in the Istio control plane`,
		Example: `istioctl dashboard controlz pilot-123-456.istio-system`,
		RunE: func(c *cobra.Command, args []string) error {
			if labelSelector == "" && len(args) < 1 {
				c.Println(c.UsageString())
				return fmt.Errorf("specify a pod or --selector")
			}

			if labelSelector != "" && len(args) > 0 {
				c.Println(c.UsageString())
				return fmt.Errorf("name cannot be provided when a selector is specified")
			}

			client, err := clientExecFactory(kubeconfig, configContext, opts)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			var podName, ns string
			if labelSelector != "" {
				pl, err := client.PodsForSelector(handlers.HandleNamespace(namespace, defaultNamespace), labelSelector)
				if err != nil {
					return fmt.Errorf("not able to locate pod with selector %s: %v", labelSelector, err)
				}

				if len(pl.Items) < 1 {
					return errors.New("no pods found")
				}

				if len(pl.Items) > 1 {
					log.Warnf("more than 1 pods fits selector: %s; will use pod: %s", labelSelector, pl.Items[0].Name)
				}

				// only use the first pod in the list
				podName = pl.Items[0].Name
				ns = pl.Items[0].Namespace
			} else {
				podName, ns = handlers.InferPodInfo(args[0], handlers.HandleNamespace(namespace, defaultNamespace))
			}

			return portForward(podName, ns, fmt.Sprintf("ControlZ %s", podName),
				"http://localhost:%d", bindAddress, controlZport, client, c.OutOrStdout())
		},
	}

	return cmd
}

// portForward first tries to forward localhost:remotePort to podName:remotePort, falls back to dynamic local port
func portForward(podName, namespace, flavor, url, localAddr string, remotePort int, client kubernetes.ExecClient, writer io.Writer) error {
	var err error
	for _, localPort := range []int{listenPort, remotePort} {
		fw, err := client.BuildPortForwarder(podName, namespace, localAddr, localPort, remotePort)
		if err != nil {
			return fmt.Errorf("could not build port forwarder for %s: %v", flavor, err)
		}

		if err = kubernetes.RunPortForwarder(fw, func(fw *kubernetes.PortForward) error {
			log.Debugf(fmt.Sprintf("port-forward to %s pod ready", flavor))
			openBrowser(fmt.Sprintf(url, fw.LocalPort), writer)
			return nil
		}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failure running port forward process: %v", err)
}

func openBrowser(url string, writer io.Writer) {
	var err error

	fmt.Fprintf(writer, "%s\n", url)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Fprintf(writer, "Unsupported platform %q; open %s in your browser.\n", runtime.GOOS, url)
	}

	if err != nil {
		fmt.Fprintf(writer, "Failed to open browser; open %s in your browser.\n", url)
	}

}

func dashboard() *cobra.Command {
	dashboardCmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"dash", "d"},
		Short:   "Access to Istio web UIs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, args)
			if len(args) != 0 {
				return fmt.Errorf("unknown dashboard %q", args[0])
			}

			return nil
		},
	}

	dashboardCmd.PersistentFlags().IntVarP(&listenPort, "port", "p", 0, "Local port to listen to")
	dashboardCmd.PersistentFlags().StringVar(&bindAddress, "address", "localhost",
		"Address to listen on. Only accepts IP address or localhost as a value. "+
			"When localhost is supplied, istioctl will try to bind on both 127.0.0.1 and ::1 "+
			"and will fail if neither of these address are available to bind.")

	dashboardCmd.AddCommand(kialiDashCmd())
	dashboardCmd.AddCommand(promDashCmd())
	dashboardCmd.AddCommand(grafanaDashCmd())
	dashboardCmd.AddCommand(jaegerDashCmd())
	dashboardCmd.AddCommand(zipkinDashCmd())

	envoy := envoyDashCmd()
	envoy.PersistentFlags().StringVarP(&labelSelector, "selector", "l", "", "label selector")
	dashboardCmd.AddCommand(envoy)

	controlz := controlZDashCmd()
	controlz.PersistentFlags().IntVar(&controlZport, "ctrlz_port", 9876, "ControlZ port")
	controlz.PersistentFlags().StringVarP(&labelSelector, "selector", "l", "", "label selector")
	dashboardCmd.AddCommand(controlz)

	return dashboardCmd
}
