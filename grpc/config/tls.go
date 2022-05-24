// Copyright (c) 2022 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"crypto/tls"
)

func getBaseTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256, tls.X25519},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			/* Not blacklisted cyphers:

			cipher_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256
			cipher_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384
			cipher_TLS_DHE_DSS_WITH_AES_128_GCM_SHA256
			cipher_TLS_DHE_DSS_WITH_AES_256_GCM_SHA384
			cipher_TLS_DHE_PSK_WITH_AES_128_GCM_SHA256
			cipher_TLS_DHE_PSK_WITH_AES_256_GCM_SHA384
			cipher_TLS_FALLBACK_SCSV
			cipher_TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
			cipher_TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
			cipher_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
			cipher_TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
			cipher_TLS_DHE_RSA_WITH_ARIA_128_GCM_SHA256
			cipher_TLS_DHE_RSA_WITH_ARIA_256_GCM_SHA384
			cipher_TLS_DHE_DSS_WITH_ARIA_128_GCM_SHA256
			cipher_TLS_DHE_DSS_WITH_ARIA_256_GCM_SHA384
			cipher_TLS_ECDHE_ECDSA_WITH_ARIA_128_GCM_SHA256
			cipher_TLS_ECDHE_ECDSA_WITH_ARIA_256_GCM_SHA384
			cipher_TLS_ECDHE_RSA_WITH_ARIA_128_GCM_SHA256
			cipher_TLS_ECDHE_RSA_WITH_ARIA_256_GCM_SHA384
			cipher_TLS_DHE_PSK_WITH_ARIA_128_GCM_SHA256
			cipher_TLS_DHE_PSK_WITH_ARIA_256_GCM_SHA384
			cipher_TLS_DHE_RSA_WITH_CAMELLIA_128_GCM_SHA256
			cipher_TLS_DHE_RSA_WITH_CAMELLIA_256_GCM_SHA384
			cipher_TLS_DHE_DSS_WITH_CAMELLIA_128_GCM_SHA256
			cipher_TLS_DHE_DSS_WITH_CAMELLIA_256_GCM_SHA384
			cipher_TLS_ECDHE_ECDSA_WITH_CAMELLIA_128_GCM_SHA256
			cipher_TLS_ECDHE_ECDSA_WITH_CAMELLIA_256_GCM_SHA384
			cipher_TLS_ECDHE_RSA_WITH_CAMELLIA_128_GCM_SHA256
			cipher_TLS_ECDHE_RSA_WITH_CAMELLIA_256_GCM_SHA384
			cipher_TLS_DHE_PSK_WITH_CAMELLIA_128_GCM_SHA256
			cipher_TLS_DHE_PSK_WITH_CAMELLIA_256_GCM_SHA384
			cipher_TLS_DHE_RSA_WITH_AES_128_CCM
			cipher_TLS_DHE_RSA_WITH_AES_256_CCM
			cipher_TLS_DHE_RSA_WITH_AES_128_CCM_8
			cipher_TLS_DHE_RSA_WITH_AES_256_CCM_8
			cipher_TLS_DHE_PSK_WITH_AES_128_CCM
			cipher_TLS_DHE_PSK_WITH_AES_256_CCM
			cipher_TLS_PSK_DHE_WITH_AES_128_CCM_8
			cipher_TLS_PSK_DHE_WITH_AES_256_CCM_8
			cipher_TLS_ECDHE_ECDSA_WITH_AES_128_CCM
			cipher_TLS_ECDHE_ECDSA_WITH_AES_256_CCM
			cipher_TLS_ECDHE_ECDSA_WITH_AES_128_CCM_8
			cipher_TLS_ECDHE_ECDSA_WITH_AES_256_CCM_8
			cipher_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_DHE_RSA_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_PSK_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_ECDHE_PSK_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_DHE_PSK_WITH_CHACHA20_POLY1305_SHA256
			cipher_TLS_RSA_PSK_WITH_CHACHA20_POLY1305_SHA256
			*/
			// v1.2
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_FALLBACK_SCSV,

			// v1.3
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		NextProtos: []string{"h2"},
	}
}

// DefaultServerTLSConfig setting https://wiki.mozilla.org/Security/Server_Side_TLS
func DefaultServerTLSConfig(serverKeyPair *tls.Certificate) (serverConfig *tls.Config) {
	serverConfig = getBaseTLSConfig()
	serverConfig.Certificates = append(serverConfig.Certificates, *serverKeyPair)
	return
}

// ClientTLSConfig returns TLS Config for client
func ClientTLSConfig() (clientTLSConf *tls.Config, err error) {
	clientTLSConf = getBaseTLSConfig()
	return
}
