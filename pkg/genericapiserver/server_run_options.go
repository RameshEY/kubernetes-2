/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package genericapiserver

import (
	"net"

	"k8s.io/kubernetes/pkg/util/config"
	utilnet "k8s.io/kubernetes/pkg/util/net"

	"github.com/spf13/pflag"
)

const (
	// TODO: This can be tightened up. It still matches objects named watch or proxy.
	defaultLongRunningRequestRE = "(/|^)((watch|proxy)(/|$)|(logs?|portforward|exec|attach)/?$)"
)

// ServerRunOptions contains the options while running a generic api server.
type ServerRunOptions struct {
	APIGroupPrefix            string
	APIPrefix                 string
	BindAddress               net.IP
	CertDirectory             string
	AdvertiseAddress          net.IP
	ClientCAFile              string
	CorsAllowedOriginList     []string
	EnableLogsSupport         bool
	EnableProfiling           bool
	EnableSwaggerUI           bool
	EnableWatchCache          bool
	ExternalHost              string
	InsecureBindAddress       net.IP
	InsecurePort              int
	KubernetesServiceNodePort int
	LongRunningRequestRE      string
	MasterCount               int
	MaxRequestsInFlight       int
	MinRequestTimeout         int
	RuntimeConfig             config.ConfigurationMap
	SecurePort                int
	ServiceClusterIPRange     net.IPNet // TODO: make this a list
	ServiceNodePortRange      utilnet.PortRange
	TLSCertFile               string
	TLSPrivateKeyFile         string
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		APIGroupPrefix:       "/apis",
		APIPrefix:            "/api",
		BindAddress:          net.ParseIP("0.0.0.0"),
		CertDirectory:        "/var/run/kubernetes",
		EnableLogsSupport:    true,
		EnableProfiling:      true,
		EnableWatchCache:     true,
		InsecureBindAddress:  net.ParseIP("127.0.0.1"),
		InsecurePort:         8080,
		LongRunningRequestRE: defaultLongRunningRequestRE,
		MasterCount:          1,
		MaxRequestsInFlight:  400,
		MinRequestTimeout:    1800,
		RuntimeConfig:        make(config.ConfigurationMap),
		SecurePort:           6443,
	}
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs.IPVar(&s.AdvertiseAddress, "advertise-address", s.AdvertiseAddress, ""+
		"The IP address on which to advertise the apiserver to members of the cluster. This "+
		"address must be reachable by the rest of the cluster. If blank, the --bind-address "+
		"will be used. If --bind-address is unspecified, the host's default interface will "+
		"be used.")

	fs.IPVar(&s.BindAddress, "public-address-override", s.BindAddress, "DEPRECATED: see --bind-address instead")
	fs.MarkDeprecated("public-address-override", "see --bind-address instead")
	fs.IPVar(&s.BindAddress, "bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure-port port. The "+
		"associated interface(s) must be reachable by the rest of the cluster, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0).")

	fs.StringVar(&s.CertDirectory, "cert-dir", s.CertDirectory, "The directory where the TLS certs are located (by default /var/run/kubernetes). "+
		"If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored.")

	fs.StringVar(&s.ClientCAFile, "client-ca-file", s.ClientCAFile, "If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.")
	fs.StringSliceVar(&s.CorsAllowedOriginList, "cors-allowed-origins", s.CorsAllowedOriginList, "List of allowed origins for CORS, comma separated.  An allowed origin can be a regular expression to support subdomain matching.  If this list is empty CORS will not be enabled.")

	fs.BoolVar(&s.EnableProfiling, "profiling", s.EnableProfiling, "Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&s.EnableSwaggerUI, "enable-swagger-ui", s.EnableSwaggerUI, "Enables swagger ui on the apiserver at /swagger-ui")

	// TODO: enable cache in integration tests.
	fs.BoolVar(&s.EnableWatchCache, "watch-cache", s.EnableWatchCache, "Enable watch caching in the apiserver")

	fs.StringVar(&s.ExternalHost, "external-hostname", s.ExternalHost, "The hostname to use when generating externalized URLs for this master (e.g. Swagger API Docs.)")

	fs.IPVar(&s.InsecureBindAddress, "insecure-bind-address", s.InsecureBindAddress, ""+
		"The IP address on which to serve the --insecure-port (set to 0.0.0.0 for all interfaces). "+
		"Defaults to localhost.")
	fs.IPVar(&s.InsecureBindAddress, "address", s.InsecureBindAddress, "DEPRECATED: see --insecure-bind-address instead")
	fs.MarkDeprecated("address", "see --insecure-bind-address instead")

	fs.IntVar(&s.InsecurePort, "insecure-port", s.InsecurePort, ""+
		"The port on which to serve unsecured, unauthenticated access. Default 8080. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the cluster and that port 443 on the cluster's public address is proxied to this "+
		"port. This is performed by nginx in the default setup.")
	fs.IntVar(&s.InsecurePort, "port", s.InsecurePort, "DEPRECATED: see --insecure-port instead")
	fs.MarkDeprecated("port", "see --insecure-port instead")

	// See #14282 for details on how to test/try this option out.  TODO remove this comment once this option is tested in CI.
	fs.IntVar(&s.KubernetesServiceNodePort, "kubernetes-service-node-port", s.KubernetesServiceNodePort, "If non-zero, the Kubernetes master service (which apiserver creates/maintains) will be of type NodePort, using this as the value of the port. If zero, the Kubernetes master service will be of type ClusterIP.")

	fs.StringVar(&s.LongRunningRequestRE, "long-running-request-regexp", s.LongRunningRequestRE, "A regular expression matching long running requests which should be excluded from maximum inflight request handling.")

	fs.IntVar(&s.MasterCount, "apiserver-count", s.MasterCount, "The number of apiservers running in the cluster")

	fs.IntVar(&s.MaxRequestsInFlight, "max-requests-inflight", s.MaxRequestsInFlight, "The maximum number of requests in flight at a given time.  When the server exceeds this, it rejects requests.  Zero for no limit.")

	fs.IntVar(&s.MinRequestTimeout, "min-request-timeout", s.MinRequestTimeout, "An optional field indicating the minimum number of seconds a handler must keep a request open before timing it out. Currently only honored by the watch request handler, which picks a randomized value above this number as the connection timeout, to spread out load.")

	fs.Var(&s.RuntimeConfig, "runtime-config", "A set of key=value pairs that describe runtime configuration that may be passed to apiserver. apis/<groupVersion> key can be used to turn on/off specific api versions. apis/<groupVersion>/<resource> can be used to turn on/off specific resources. api/all and api/legacy are special keys to control all and legacy api versions respectively.")

	fs.IntVar(&s.SecurePort, "secure-port", s.SecurePort, ""+
		"The port on which to serve HTTPS with authentication and authorization. If 0, "+
		"don't serve HTTPS at all.")

	fs.IPNetVar(&s.ServiceClusterIPRange, "service-cluster-ip-range", s.ServiceClusterIPRange, "A CIDR notation IP range from which to assign service cluster IPs. This must not overlap with any IP ranges assigned to nodes for pods.")
	fs.IPNetVar(&s.ServiceClusterIPRange, "portal-net", s.ServiceClusterIPRange, "Deprecated: see --service-cluster-ip-range instead.")
	fs.MarkDeprecated("portal-net", "see --service-cluster-ip-range instead.")

	fs.Var(&s.ServiceNodePortRange, "service-node-port-range", "A port range to reserve for services with NodePort visibility.  Example: '30000-32767'.  Inclusive at both ends of the range.")
	fs.Var(&s.ServiceNodePortRange, "service-node-ports", "Deprecated: see --service-node-port-range instead.")
	fs.MarkDeprecated("service-node-ports", "see --service-node-port-range instead.")

	fs.StringVar(&s.TLSCertFile, "tls-cert-file", s.TLSCertFile, ""+
		"File containing x509 Certificate for HTTPS.  (CA cert, if any, concatenated after server cert). "+
		"If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, "+
		"a self-signed certificate and key are generated for the public address and saved to /var/run/kubernetes.")

	fs.StringVar(&s.TLSPrivateKeyFile, "tls-private-key-file", s.TLSPrivateKeyFile, "File containing x509 private key matching --tls-cert-file.")
}
