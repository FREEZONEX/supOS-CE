_format_version: '3.0'
_transform: false
consumers:
- tags: []
  created_at: 1734329209
  id: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  updated_at: 1764811192
  custom_id: ~
  username: open-api
parameters:
- key: cluster_id
  value: ce8f9346-2b44-46af-af55-ab141518e1bb
  created_at: 1742432166
services:
- name: grafana
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729740320
  port: 3000
  updated_at: 1764811192
  enabled: true
  host: grafana
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 228308d0-521c-45f8-a97d-fbe6966efa3c
  client_certificate: ~
  read_timeout: 60000
- name: gitea
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729852876
  port: 3000
  updated_at: 1764811192
  enabled: true
  host: gitea
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  client_certificate: ~
  read_timeout: 60000
- name: marimo
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1760594820
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: marimo
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 26effc32-8aa5-4d21-b86a-f14add8afc87
  client_certificate: ~
  read_timeout: 60000
- name: portainer
  path: ~
  protocol: https
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729740912
  port: 9443
  updated_at: 1764811192
  enabled: true
  host: portainer
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 2cd97a86-3853-4001-9f0d-7769dc40d508
  client_certificate: ~
  read_timeout: 60000
- name: konga
  path: /
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729736772
  port: 1337
  updated_at: 1764811192
  enabled: true
  host: konga
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 2df2f7c3-792a-4a30-b30c-6011e2f3f916
  client_certificate: ~
  read_timeout: 60000
- name: plugin-frontend
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1749280964
  port: 3001
  updated_at: 1764811192
  enabled: true
  host: plugin
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 33c02e25-db50-4332-a044-094fce5f50b0
  client_certificate: ~
  read_timeout: 60000
- name: backend-open-api
  path: /open-api/
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1734330146
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 4007e6c1-6ccc-4747-9b96-7bb3f5f78b32
  client_certificate: ~
  read_timeout: 60000
- name: event-flow-proxy
  path: /service-api/supos/proxy/event/flows
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1742976133
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 42684a67-ac8f-48d6-ae2b-62e1cd26f9d8
  client_certificate: ~
  read_timeout: 60000
- name: backend-service-api
  path: /service-api/supos/
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1733290847
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 4d210886-a694-416d-847d-95fc597f5921
  client_certificate: ~
  read_timeout: 60000
- name: gateway
  path: /
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1732610867
  port: 8070
  updated_at: 1764811192
  enabled: true
  host: gateway
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 5e8458a8-7323-4a25-9efa-4d56dbf5fa5b
  client_certificate: ~
  read_timeout: 60000
- name: backend
  path: /inter-api/supos/
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729740054
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 300000
  connect_timeout: 300000
  tls_verify_depth: ~
  id: 5f70fd49-e3d7-4ba8-b152-62eca6ec4455
  client_certificate: ~
  read_timeout: 300000
- name: minio
  path: /
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1731396402
  port: 9001
  updated_at: 1764811192
  enabled: true
  host: minio
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 647747b1-8efe-45b6-b37f-430f5e5139d6
  client_certificate: ~
  read_timeout: 60000
- name: EventFlow
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1742536198
  port: 1889
  updated_at: 1764811192
  enabled: true
  host: eventflow
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  client_certificate: ~
  read_timeout: 60000
- name: backend-files
  path: /files
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1741670682
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 75a7373b-4dda-4a49-b1f8-de6ebee4d4c8
  client_certificate: ~
  read_timeout: 60000
- name: emqx
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729740676
  port: 18083
  updated_at: 1764811192
  enabled: true
  host: emqx
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 89c6a77c-e3e6-4149-a0e7-6ace7fd0413b
  client_certificate: ~
  read_timeout: 60000
- name: frontend
  path: ~
  protocol: http
  tags:
  - root:frontend
  - Home:1
  - SourceFlow:2
  - Namespace:3
  - EventFlow:4
  - CollectionGatewayManagement:5
  - menu.tag.devtools:6
  - menu.tag.uns:1
  - menu.tag.appspace:7
  - menu.tag.system:8
  ca_certificates: ~
  retries: 5
  created_at: 1729738059
  port: 3000
  updated_at: 1764811192
  enabled: true
  host: frontend
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  client_certificate: ~
  read_timeout: 60000
- name: eventflow-backend
  path: /eventflow-api/
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1755572726
  port: 1888
  updated_at: 1764811192
  enabled: true
  host: eventflow
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: 94797f49-35e2-428f-b46a-84788fa06fa8
  client_certificate: ~
  read_timeout: 60000
- name: swagger
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1757603297
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  client_certificate: ~
  read_timeout: 60000
- name: node-red-proxy
  path: /service-api/supos/proxy/flows
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1730685511
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: uns
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: b096bcf5-2984-4acc-9bd5-a570a7653fcd
  client_certificate: ~
  read_timeout: 60000
- name: keycloak
  path: /
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729740513
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: keycloak
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  client_certificate: ~
  read_timeout: 60000
- name: nodered
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729739239
  port: 1880
  updated_at: 1764811192
  enabled: true
  host: nodered
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: bba8a174-4679-49df-8bf9-ae9285f1e77e
  client_certificate: ~
  read_timeout: 60000
- name: minio-inter
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1731460027
  port: 9000
  updated_at: 1764811192
  enabled: true
  host: minio-inter
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: cc3e8a53-82e1-4f6e-bde6-75207ca2f6d3
  client_certificate: ~
  read_timeout: 60000
- name: nodered-backend
  path: /nodered-api/
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1744972980
  port: 1880
  updated_at: 1764811192
  enabled: true
  host: nodered
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: df4453f0-063a-4346-a074-3c2f9a388ca7
  client_certificate: ~
  read_timeout: 60000
- name: login
  path: /realms/tier0/protocol/openid-connect/auth
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1732105495
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: keycloak
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: e3e88607-311a-4c23-a9c7-bb879efc463e
  client_certificate: ~
  read_timeout: 60000
- name: apm
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1730264819
  port: 8080
  updated_at: 1764811192
  enabled: true
  host: apm
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  client_certificate: ~
  read_timeout: 60000
- name: GenerativeUI
  path: /
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1729748350
  port: 4000
  updated_at: 1764811192
  enabled: true
  host: copilotkit
  tls_verify: ~
  write_timeout: 600000
  connect_timeout: 600000
  tls_verify_depth: ~
  id: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  client_certificate: ~
  read_timeout: 600000
- name: mcpclient
  path: ~
  protocol: http
  tags: []
  ca_certificates: ~
  retries: 5
  created_at: 1742543141
  port: 3000
  updated_at: 1764811192
  enabled: true
  host: mcpclient
  tls_verify: ~
  write_timeout: 60000
  connect_timeout: 60000
  tls_verify_depth: ~
  id: fbc56015-44d7-42ba-ac9f-abe725bc2478
  client_certificate: ~
  read_timeout: 60000
routes:
- name: WebHooks
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /WebHooks
  created_at: 1757902554
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 0d271f38-f4a5-4008-8bb0-33c706460be3
  tags: ~
  protocols:
  - http
  - https
- name: swagger-config
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  strip_path: false
  paths:
  - /v3/api-docs
  created_at: 1757604045
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 0f376c28-8c9d-4a4a-8c6b-5258d9d8edef
  tags: ~
  protocols:
  - http
  - https
- name: grafana-inter
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 228308d0-521c-45f8-a97d-fbe6966efa3c
  strip_path: true
  paths:
  - /grafana/home/
  created_at: 1730270842
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 10c662d8-8304-439b-85a0-398ae09d09e8
  tags: ~
  protocols:
  - http
  - https
- name: CICD
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  strip_path: true
  paths:
  - /gitea/home/user/login?redirect_to=/gitea/home/
  created_at: 1730254514
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 15023439-25a9-481d-90a5-894d4d87d3bf
  tags:
  - description:menu.desc.cicd
  - sort:1
  - parentName:menu.tag.devtools
  - menu
  protocols:
  - http
  - https
- name: GenApps
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /app-display
  created_at: 1731311256
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 1dfa8b71-07e8-41a5-a07d-536fa0e295a2
  tags:
  - description:menu.desc.genApps
  - sort:2
  - parentName:menu.tag.appspace
  protocols:
  - http
  - https
- name: apm-backend-iner
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  strip_path: true
  paths:
  - /apps/freezonex-aps/
  created_at: 1730265781
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 209f6169-707c-49f1-86cc-602bc4674b2e
  tags: ~
  protocols:
  - http
  - https
- name: MqttBroker
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 89c6a77c-e3e6-4149-a0e7-6ace7fd0413b
  strip_path: true
  paths:
  - /emqx/home/
  created_at: 1729740717
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 215cdb87-6e13-4e2e-804f-9ea1ac2ff417
  tags:
  - description:menu.desc.mqttBroker
  - parentName:menu.tag.system
  - sort:3
  protocols:
  - http
  - https
- name: eventflow-backend
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 94797f49-35e2-428f-b46a-84788fa06fa8
  strip_path: true
  paths:
  - '/eventflow-api/ '
  created_at: 1755572773
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 2695543d-3af1-4ca5-b515-82e0fe9b69c7
  tags: ~
  protocols:
  - http
  - https
- name: copilotkit
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  strip_path: false
  paths:
  - /copilotkit
  created_at: 1729748378
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 29027b71-49ce-41d1-96fc-14bcb3a2da00
  tags: ~
  protocols:
  - http
  - https
- name: PluginManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /plugin-management
  created_at: 1748310135
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 366d2cfa-8a0e-4a6b-943e-3b5ecd424f73
  tags:
  - menu
  - description:menu.desc.PluginManagement
  - parentName:menu.tag.system
  - sort:14
  protocols:
  - http
  - https
- name: backend
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 5f70fd49-e3d7-4ba8-b152-62eca6ec4455
  strip_path: true
  paths:
  - /inter-api/supos/
  created_at: 1729740083
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 3794799e-0c23-4065-a88d-7a08c46fbaf4
  tags: ~
  protocols:
  - http
  - https
- name: node-red-flows
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods:
  - GET
  https_redirect_status_code: 426
  service: b096bcf5-2984-4acc-9bd5-a570a7653fcd
  strip_path: true
  paths:
  - /nodered/home/flows
  created_at: 1730685544
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 3f8f8bf7-4d53-4a11-bc6d-d017beda8695
  tags: []
  protocols:
  - http
  - https
- name: EventFlowBackend
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  strip_path: true
  paths:
  - /eventflow/home/
  created_at: 1742536226
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 45a9fec1-6eaa-49c5-ae00-7aa180f1efde
  tags: []
  protocols:
  - http
  - https
- name: Namespace
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /uns
  created_at: 1731311044
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 4d89ed56-90bf-490d-ac8c-30def8be2e2c
  tags:
  - menu
  - description:menu.desc.dataModeling
  - homeParentName:menu.tag.uns
  - homeIconUrl:homeNamespace
  - sort:2
  protocols:
  - http
  - https
- name: backend-service-api
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 4d210886-a694-416d-847d-95fc597f5921
  strip_path: true
  paths:
  - /service-api/supos/
  created_at: 1733290917
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 5a0cdda1-0ac2-4255-a4bc-a11a8b4a00d5
  tags: ~
  protocols:
  - http
  - https
- name: nodered-backend
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: df4453f0-063a-4346-a074-3c2f9a388ca7
  strip_path: true
  paths:
  - /nodered-api/
  created_at: 1744972998
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 5a9793a9-5085-4830-bba7-a5c053a055a1
  tags: ~
  protocols:
  - http
  - https
- name: McpClient
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: fbc56015-44d7-42ba-ac9f-abe725bc2478
  strip_path: false
  paths:
  - /mcpclient/home
  created_at: 1742543248
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 6239b86d-4ef1-48d9-b512-7b05306ab705
  tags:
  - ${ENABLE_MCP}
  - parentName:menu.tag.appspace
  - description:menu.desc.mcpclient
  - sort:3
  protocols:
  - http
  - https
- name: backend-files
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 75a7373b-4dda-4a49-b1f8-de6ebee4d4c8
  strip_path: true
  paths:
  - /files
  created_at: 1741671187
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 670008e2-5811-4d9d-b925-5429d16caa8f
  tags: ~
  protocols:
  - http
  - https
- name: CodeManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /CodeManagement
  created_at: 1757902539
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 6737cbc9-f115-4018-bd0e-90a2342bb002
  tags: ~
  protocols:
  - http
  - https
- name: gateway
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 5e8458a8-7323-4a25-9efa-4d56dbf5fa5b
  strip_path: true
  paths:
  - /gateway
  created_at: 1732611195
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 6ce6b319-0e38-4bf4-ba1d-4d043a926ba7
  tags: ~
  protocols:
  - http
  - https
- name: apm-frontend-inter
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  strip_path: false
  paths:
  - /apps/freezonex-aps/apsfrontend/
  created_at: 1730264872
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 6f2a8d9e-d76f-46b9-b2d4-27286651433d
  tags: ~
  protocols:
  - http
  - https
- name: Alert
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /Alert
  created_at: 1758346489
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 71307652-070e-4ef1-8b30-704d186f3b75
  tags: ~
  protocols:
  - http
  - https
- name: RoutingManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 2df2f7c3-792a-4a30-b30c-6011e2f3f916
  strip_path: true
  paths:
  - /konga/home/
  created_at: 1729736896
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 79d5e57d-340c-4d18-93ee-a6a8f4a0f212
  tags:
  - description:menu.desc.konga
  - sort:1
  - parentName:menu.tag.system
  - menu
  protocols:
  - http
  - https
- name: apm
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  strip_path: false
  paths:
  - /apsfrontend/dashboard
  created_at: 1730265586
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 833d37e3-05bb-4a1d-992c-0a0e5c19b0dd
  tags:
  - description:menu.desc.apm
  protocols:
  - http
  - https
- name: plugin-frontend
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 33c02e25-db50-4332-a044-094fce5f50b0
  strip_path: true
  paths:
  - /plugin/
  created_at: 1749280989
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 84235d0d-bc92-4a86-9413-a66807af7f90
  tags: ~
  protocols:
  - http
  - https
- name: StreamProcessing
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /streams
  created_at: 1733209593
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 8613d94b-0fd7-4b81-8f9d-ac5df267b7e3
  tags:
  - parentName:menu.tag.connections
  protocols:
  - http
  - https
- name: AppManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /AppManagement
  created_at: 1757902544
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 88be7566-488d-4207-b54e-abcb082e2c5e
  tags: ~
  protocols:
  - http
  - https
- name: ThemeManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /ThemeManagement
  created_at: 1764811231
  updated_at: 1764811231
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 89f3d3e9-97fc-48e8-baf0-4751b83992cb
  tags: ~
  protocols:
  - http
  - https
- name: open-backend-api
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 4007e6c1-6ccc-4747-9b96-7bb3f5f78b32
  strip_path: true
  paths:
  - /open-api/
  created_at: 1734330177
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  tags: ~
  protocols:
  - http
  - https
- name: SourceFlow
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /collection-flow
  created_at: 1731311377
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: a13bb597-9740-4dde-929e-d140c286d869
  tags:
  - menu
  - description:menu.desc.nodered.flow
  - sort:1
  - homeParentName:menu.tag.uns
  - homeIconUrl:homeSourceFlow
  protocols:
  - http
  - https
- name: GenerativeUI
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /app-space
  created_at: 1731311359
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: a2060aa0-88e5-4247-9635-f93438bbdd84
  tags:
  - description:menu.desc.generativeUI
  - parentName:menu.tag.appspace
  - sort:2
  protocols:
  - http
  - https
- name: workflow
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /workflow
  created_at: 1741573643
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: a5040934-75dd-40c6-94ea-9497ab2b0579
  tags:
  - parentName:menu.tag.settings
  protocols:
  - http
  - https
- name: objectStorageServer
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 647747b1-8efe-45b6-b37f-430f5e5139d6
  strip_path: true
  paths:
  - /minio/home/
  created_at: 1731396438
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: a6d04fe9-a464-493c-8f3a-4750fdd93a32
  tags:
  - description:menu.desc.objectStorageServer
  - sort:200
  protocols:
  - http
  - https
- name: grafana
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 228308d0-521c-45f8-a97d-fbe6966efa3c
  strip_path: true
  paths:
  - /grafana/home/dashboards/
  created_at: 1730270517
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: aa02bbb9-6459-43bd-9b65-91d89c8854dd
  tags:
  - description:menu.desc.grafana
  protocols:
  - http
  - https
- name: keycloak-auth
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  strip_path: true
  paths:
  - /keycloak/home/auth/
  created_at: 1731473911
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: b610973a-764e-4cef-910e-0794f334e4bd
  tags: ~
  protocols:
  - http
  - https
- name: EventFlow
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /EventFlow
  created_at: 1742968905
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: b8262364-32bf-4422-9d6c-04b97bc8c3a7
  tags:
  - menu
  - homeParentName:menu.tag.uns
  - description:menu.desc.eventflow
  - homeIconUrl:homeEventFlow
  - sort:3
  protocols:
  - http
  - https
- name: login
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
  strip_path: true
  paths:
  - /tier0-login
  created_at: 1732108769
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: ba7b2df9-b0d8-4b6b-844d-43f935f3181f
  tags: ~
  protocols:
  - http
  - https
- name: swagger-ui
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  strip_path: false
  paths:
  - /swagger-ui
  created_at: 1757603385
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: c010b641-a378-4d8f-b377-08a5ac8c79d1
  tags: ~
  protocols:
  - http
  - https
- name: frontend
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /
  created_at: 1729738250
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: c2dececa-99f4-45e1-9859-01e88352bd58
  tags: ~
  protocols:
  - http
  - https
- name: CollectionGatewayManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /CollectionGatewayManagement
  created_at: 1757902564
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: c810eaa9-dbdd-475a-80f6-6cd9a08948a6
  tags: ~
  protocols:
  - http
  - https
- name: marimo
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 26effc32-8aa5-4d21-b86a-f14add8afc87
  strip_path: true
  paths:
  - /marimo/home/
  created_at: 1760594842
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: c8fd090f-2ff9-499c-807b-139f6e5fb976
  tags: ~
  protocols:
  - http
  - https
- name: Home
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /home
  created_at: 1731635333
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: c90b4b7a-8a51-4f40-b4e2-6c0a40be1b15
  tags:
  - menu
  - description:menu.desc.home
  - sort:1
  protocols:
  - http
  - https
- name: health
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  strip_path: false
  paths:
  - /open-api/health
  created_at: 1758351166
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: caeab40c-0c09-4465-b54a-2609391a19e8
  tags: ~
  protocols:
  - http
  - https
- name: Authentication
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  strip_path: true
  paths:
  - /keycloak/home/
  created_at: 1729740574
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: d0cea78f-1e0d-4b90-98ea-980a455bf5f5
  tags:
  - description:menu.desc.keycloak
  - menu
  - sort:2
  - parentName:menu.tag.system
  protocols:
  - http
  - https
- name: Dashboards
  sources: ~
  preserve_host: true
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /dashboards
  created_at: 1730770040
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: d2a81d6f-8b2a-4b28-8929-3c51ccd16021
  tags:
  - menu
  - description:menu.desc.dashboards
  - parentName:menu.tag.system
  - sort:5
  protocols:
  - http
  - https
- name: minio-inter
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: cc3e8a53-82e1-4f6e-bde6-75207ca2f6d3
  strip_path: true
  paths:
  - /minio/inter/
  created_at: 1731460050
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: dbb92267-e886-4ee8-b758-a9f9e9af1998
  tags: ~
  protocols:
  - http
  - https
- name: gitea-inter
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  strip_path: true
  paths:
  - /gitea/home/
  created_at: 1729852903
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: e3459f2d-fcb2-412e-87fc-b098d8906b7e
  tags: []
  protocols:
  - http
  - https
- name: event-node-flows
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 42684a67-ac8f-48d6-ae2b-62e1cd26f9d8
  strip_path: true
  paths:
  - /eventflow/home/flows
  created_at: 1742976255
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: e8e7fe7d-16ba-415a-8d19-e2c41b76b365
  tags: []
  protocols:
  - http
  - https
- name: Connection
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /Connection
  created_at: 1764811222
  updated_at: 1764811222
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: e92d4fec-518a-400b-9094-1a91e989769b
  tags: ~
  protocols:
  - http
  - https
- name: AdvancedUse
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /advanced-use
  created_at: 1734056913
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: f0a59836-7eea-45b3-b188-51c45c68f305
  tags:
  - menu
  - description:menu.desc.advanceUse
  - sort:9
  - parentName:menu.tag.system
  protocols:
  - http
  - https
- name: dashboard2
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  strip_path: false
  paths:
  - /dashboard
  - /dashboard/
  created_at: 1760593302
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: f3121c84-e787-4eb6-a054-edbc56df4d0e
  tags: ~
  protocols:
  - http
  - https
- name: UserManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /account-management
  created_at: 1732696995
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: f610002e-8e5c-424a-9fdc-e86d754b3142
  tags:
  - menu
  - description:menu.desc.account
  - parentName:menu.tag.system
  - sort:4
  protocols:
  - http
  - https
- name: NodeRed
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: bba8a174-4679-49df-8bf9-ae9285f1e77e
  strip_path: true
  paths:
  - /nodered/home/
  created_at: 1729739343
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: f8bccdbf-46fd-4dd7-aafe-d1a8e13c3450
  tags: []
  protocols:
  - http
  - https
- name: ContainerManagement
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 2cd97a86-3853-4001-9f0d-7769dc40d508
  strip_path: true
  paths:
  - /portainer/home/
  created_at: 1729740933
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: fb690c3a-e09e-4fcf-aad1-5d357c1938ec
  tags:
  - description:menu.desc.dockerMgmt
  - ${ENABLE_PORTAINER}
  - sort:6
  - parentName:menu.tag.system
  protocols:
  - http
  - https
- name: AboutUs
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /aboutus
  created_at: 1733822464
  updated_at: 1764811192
  response_buffering: true
  path_handling: v1
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: fb8efe62-0030-4540-a521-a3b4215b4f59
  tags:
  - menu
  - description:menu.desc.aboutus
  - parentName:menu.tag.system
  - sort:8
  protocols:
  - http
  - https
- name: OpenData
  sources: ~
  preserve_host: false
  destinations: ~
  headers: ~
  methods: ~
  https_redirect_status_code: 426
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  strip_path: true
  paths:
  - /OpenData
  created_at: 1758346492
  updated_at: 1764811192
  response_buffering: true
  path_handling: v0
  regex_priority: 0
  hosts: ~
  request_buffering: true
  snis: ~
  id: fd23d489-34ad-4c12-9e72-10e43318866d
  tags: ~
  protocols:
  - http
  - https
plugins:
- tags: ~
  name: supos-auth-checker
  created_at: 1733799127
  updated_at: 1764811192
  config:
    enable_deny_check: true
    login_url: ${BASE_URL}/keycloak/home/auth/realms/tier0/protocol/openid-connect/auth?client_id=tier0&redirect_uri=${BASE_URL}/inter-api/supos/auth/token&response_type=code&scope=openid
    forbidden_url: /403
    enable_resource_check: true
    whitelist_paths:
    - ^/inter-api/supos/auth.*$
    - ^/inter-api/supos/systemConfig.*$
    - ^/inter-api/supos/theme/getConfig.*$
    - ^/$
    - ^/assets.*$
    - ^/locale.*$
    - ^/logo.*$
    - ^/gitea.*git.*$
    - ^/tier0-login.*$
    - ^/403$
    - ^/open-api/.*$
    - ^/keycloak.*$
    - ^/nodered.*$
    - ^/files.*$
    - ^/freeLogin.*$
    - ^/inter-api/supos/dev/logs.*$
    - ^/inter-api/supos/license.*$
    - ^/inter-api/supos/cascade.*$
    - ^/swagger-ui.*$
    - ^/v3/api-docs.*$
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: 1845ee75-d704-40e1-a8b0-aa2baaf9d71b
  route: ~
  service: ~
- tags: ~
  name: key-auth
  created_at: 1734330234
  updated_at: 1764811192
  config:
    realm: ~
    anonymous: ~
    key_in_body: false
    run_on_preflight: true
    hide_credentials: false
    key_names:
    - apikey
    key_in_header: true
    key_in_query: true
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: 2285421f-56e3-4510-be12-69fa1040d810
  route: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  service: ~
- tags: ~
  name: request-transformer
  created_at: 1734073535
  updated_at: 1764811192
  config:
    append:
      querystring:
      - client_id:tier0
      - response_type:code
      - scope:openid
      - redirect_uri:${BASE_URL}/inter-api/supos/auth/token
      body: []
      headers: []
    remove:
      querystring: []
      body: []
      headers: []
    add:
      querystring: []
      body: []
      headers: []
    http_method: ~
    replace:
      uri: ~
      querystring: []
      headers: []
      body: []
    rename:
      querystring: []
      body: []
      headers: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: 46bda5cf-63ea-401f-9f06-b9e024aa5597
  route: ~
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
- tags: ~
  name: response-transformer
  created_at: 1734073759
  updated_at: 1764811192
  config:
    append:
      json_types: []
      json: []
      headers:
      - X-Frame-Options:DENY
      - Content-Security-Policy:frame-ancestors 'none'
    remove:
      json: []
      headers:
      - x-frame-options
      - Content-Security-Policy
    add:
      json_types: []
      json: []
      headers: []
    replace:
      json_types: []
      json: []
      headers: []
    rename:
      json: []
      headers: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: 73285cac-cf4e-4368-bf8f-d3285f9686d4
  route: ~
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
- tags: ~
  name: cors
  created_at: 1757987272
  updated_at: 1764811192
  config:
    origins:
    - '*'
    exposed_headers: ~
    max_age: ~
    headers: ~
    private_network: false
    methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE
    - OPTIONS
    preflight_continue: false
    credentials: false
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: ad24951e-c659-4588-9138-542e4586e790
  route: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  service: ~
- tags: ~
  name: response-transformer
  created_at: 1731400906
  updated_at: 1764811192
  config:
    append:
      json_types: []
      json: []
      headers: []
    remove:
      json: []
      headers:
      - X-Frame-Options
    add:
      json_types: []
      json: []
      headers:
      - X-Frame-Options:SAMEORIGIN
    replace:
      json_types: []
      json: []
      headers: []
    rename:
      json: []
      headers: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: b5722a76-60b9-483d-90ac-e5de173264e2
  route: a6d04fe9-a464-493c-8f3a-4750fdd93a32
  service: ~
- tags: ~
  name: supos-url-transformer
  created_at: 1734014838
  updated_at: 1764811192
  config:
    home_url: /?isLogin=true
  protocols:
  - grpc
  - grpcs
  - http
  - https
  instance_name: ~
  enabled: true
  consumer: ~
  id: f45e7fd0-74e5-4b36-b265-5df962eb1b58
  route: ba7b2df9-b0d8-4b6b-844d-43f935f3181f
  service: ~
upstreams:
- name: apm
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1730264799
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 00f62da1-73b3-4308-8ae5-856268388aa9
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: minio-inter
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1731459961
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 09fa87e7-1cf1-4fc3-9ad6-6861572a8ef4
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: platform
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1741754738
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 14bdf793-231d-439e-8d67-3e47a3e7da05
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: keycloak
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739799
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 14ee49e7-f9bf-4234-a48f-7b7df7dda0ea
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: nodered
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739256
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 220c1252-a48b-4e7e-af42-b134f316ed16
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: konga
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729737217
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 2a0fa8a9-98a4-4456-9d2c-faba35b54882
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: grafana
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739617
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 2e9455f1-241d-4a1e-840e-fea66dd9aa09
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: plugin
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1749280889
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 3c840281-aeb3-4969-8ea4-cc12d45423a1
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: portainer
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739903
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 3f51d2ee-7609-4f71-b6ec-8ded27b417a2
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: kibana
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739883
  updated_at: 1764811192
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 420478e2-bdc8-49ec-ba0e-cc4cfd41afc8
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: gitea
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729852910
  updated_at: 1764811192
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 47b5a73e-51c6-4ace-b507-40458125c0a6
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: copilotkit
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739930
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: 747a43cc-42c4-457a-abac-518c2fe537b3
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: hasura
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739734
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: ab6fb6a2-97b9-45e7-b508-46df2df5203e
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: frontend
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729737694
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: c0645f32-b63d-49ba-a723-5cbef9014c99
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: minio
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1731396351
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: c6844f11-b711-4f5f-a2d4-4516995790c5
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: backend
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739594
  updated_at: 1764811193
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: e6729ab9-c894-4963-bbdc-dfae17c88096
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: emqx
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1729739828
  updated_at: 1764811192
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: fa9fc31a-7d71-43bd-a453-c21704b71ac6
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
- name: fuxa
  algorithm: round-robin
  hash_on_cookie_path: /
  client_certificate: ~
  hash_fallback: none
  hash_fallback_header: ~
  slots: 1000
  hash_fallback_uri_capture: ~
  tags: ~
  created_at: 1733536118
  updated_at: 1764811192
  use_srv_name: false
  hash_fallback_query_arg: ~
  hash_on_header: ~
  host_header: ~
  hash_on_query_arg: ~
  hash_on_cookie: ~
  hash_on_uri_capture: ~
  id: fc85cfef-6e90-4592-a022-730f5d84ed99
  healthchecks:
    threshold: 0
    active:
      concurrency: 10
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      healthy:
        http_statuses:
        - 200
        - 302
        successes: 0
        interval: 0
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 404
        - 500
        - 501
        - 502
        - 503
        - 504
        - 505
        interval: 0
      type: http
      headers: ~
      timeout: 1
    passive:
      unhealthy:
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
        http_statuses:
        - 429
        - 500
        - 503
      type: http
      healthy:
        http_statuses:
        - 200
        - 201
        - 202
        - 203
        - 204
        - 205
        - 206
        - 207
        - 208
        - 226
        - 300
        - 301
        - 302
        - 303
        - 304
        - 305
        - 306
        - 307
        - 308
        successes: 0
  hash_on: none
targets:
- tags: ~
  upstream: 00f62da1-73b3-4308-8ae5-856268388aa9
  created_at: 1730264807.829
  weight: 100
  id: 0be879dd-7ba3-42e3-96a8-5be0c8dcadfe
  updated_at: 1764640500.761
  target: apm:8080
- tags: ~
  upstream: 2e9455f1-241d-4a1e-840e-fea66dd9aa09
  created_at: 1729739629.259
  weight: 100
  id: 0c66a055-3e1d-4e77-aedd-6de30939a5a2
  updated_at: 1764640500.814
  target: grafana:3000
- tags: ~
  upstream: 47b5a73e-51c6-4ace-b507-40458125c0a6
  created_at: 1729852918.487
  weight: 100
  id: 29a2c68c-7501-49b1-aef1-56e1be3c2d74
  updated_at: 1764640500.81
  target: gitea:3000
- tags: ~
  upstream: fc85cfef-6e90-4592-a022-730f5d84ed99
  created_at: 1733536131.451
  weight: 100
  id: 2cc80d27-504d-4c1a-afaf-646f6e432fd6
  updated_at: 1764640500.818
  target: fuxa:1881
- tags: ~
  upstream: c6844f11-b711-4f5f-a2d4-4516995790c5
  created_at: 1731397614.344
  weight: 100
  id: 2f312fcf-baa6-4de0-84aa-0779713ae279
  updated_at: 1764640500.805
  target: minio:9001
- tags: ~
  upstream: 2a0fa8a9-98a4-4456-9d2c-faba35b54882
  created_at: 1729737239.883
  weight: 100
  id: 38ae0eeb-fa36-40a1-8c72-aaa43e7b9ceb
  updated_at: 1764640500.823
  target: konga:1337
- tags: ~
  upstream: fa9fc31a-7d71-43bd-a453-c21704b71ac6
  created_at: 1729739844.316
  weight: 100
  id: 3ec4626c-673a-4cc1-86fa-10f3b1bdc537
  updated_at: 1764640500.694
  target: emqx:18083
- tags: ~
  upstream: 220c1252-a48b-4e7e-af42-b134f316ed16
  created_at: 1729739274.174
  weight: 100
  id: 5a923a6d-0661-43ad-a668-e89f4980e5a6
  updated_at: 1764640500.745
  target: nodered:1880
- tags: ~
  upstream: 3f51d2ee-7609-4f71-b6ec-8ded27b417a2
  created_at: 1729739915.095
  weight: 100
  id: 78410225-c503-4d2c-885d-c38a72244aa0
  updated_at: 1764640500.8
  target: portainer:9443
- tags: ~
  upstream: 14bdf793-231d-439e-8d67-3e47a3e7da05
  created_at: 1741754798.218
  weight: 100
  id: 7c5e8a1d-6022-42b6-83b7-b1f39739b9df
  updated_at: 1764640500.769
  target: frontend:3001
- tags: ~
  upstream: ab6fb6a2-97b9-45e7-b508-46df2df5203e
  created_at: 1729739753.711
  weight: 100
  id: 8d09a927-4def-479d-8ed9-0473be2a281a
  updated_at: 1764640500.78
  target: hasura:8080
- tags: ~
  upstream: 14ee49e7-f9bf-4234-a48f-7b7df7dda0ea
  created_at: 1729832116.752
  weight: 100
  id: a4cdac00-4421-4a49-b2ba-b2f720f4998d
  updated_at: 1764640500.791
  target: keycloak:8080
- tags: ~
  upstream: 420478e2-bdc8-49ec-ba0e-cc4cfd41afc8
  created_at: 1729739894.022
  weight: 100
  id: afe138dc-2d5e-4329-80ff-a0154b634f42
  updated_at: 1764640500.736
  target: kibana:5601
- tags: ~
  upstream: e6729ab9-c894-4963-bbdc-dfae17c88096
  created_at: 1729739606.748
  weight: 100
  id: e401fd6b-6a87-41be-a501-6456afeb14d3
  updated_at: 1764640500.756
  target: uns:8080
- tags: ~
  upstream: 09fa87e7-1cf1-4fc3-9ad6-6861572a8ef4
  created_at: 1731459990.339
  weight: 100
  id: e5f0e5f7-8b6a-476b-ac72-2e4a32bab571
  updated_at: 1764640500.764
  target: minio:9000
- tags: ~
  upstream: c0645f32-b63d-49ba-a723-5cbef9014c99
  created_at: 1729914708.061
  weight: 100
  id: ef3e3b44-a490-40e8-afac-2b9e2d202643
  updated_at: 1764640500.775
  target: frontend:3000
- tags: ~
  upstream: 747a43cc-42c4-457a-abac-518c2fe537b3
  created_at: 1729914781.603
  weight: 100
  id: f4e24691-4181-42a6-b32c-6c7e8cd92077
  updated_at: 1764640500.787
  target: frontend:4000
- tags: ~
  upstream: 3c840281-aeb3-4969-8ea4-cc12d45423a1
  created_at: 1749280900.4
  weight: 100
  id: f6a73dfa-96f1-4bd4-969f-a46592a5ac49
  updated_at: 1764640500.796
  target: frontend:3002
keyauth_credentials:
- tags: ~
  created_at: 1758346492
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 0b7dc033e36f4a1492ac8562885cac27
  id: 5bfc410b-ec4b-4371-a6dd-6775e38c5dc1
  ttl: ~
- tags: ~
  created_at: 1734329245
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 4174348a-9222-4e81-b33e-5d72d2fd7f1e
  id: 6b9443ae-73f0-4db6-af00-4f1e3a415dbb
  ttl: ~
- tags: ~
  created_at: 1749280821
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: d763a82bd8154b58bd29d9cd141dcab0
  id: 6b979102-3576-4a98-b06e-ea72d25249d0
  ttl: ~
- tags: ~
  created_at: 1757578203
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: d3672e250b634a4b8602d88f1ae55a81
  id: 7e12e6c7-e28b-4e5f-876b-3f49341ae403
  ttl: ~
- tags: ~
  created_at: 1749260555
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 33d911ac009240279e01b4a95655a0ad
  id: 90d62a81-2634-45a4-8384-6e155876db84
  ttl: ~
