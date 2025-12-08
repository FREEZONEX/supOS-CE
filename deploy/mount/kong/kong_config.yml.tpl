_transform: false
_format_version: '3.0'
consumers:
- created_at: 1734329209
  updated_at: 1758346462
  username: open-api
  tags: []
  id: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  custom_id: ~
parameters:
- created_at: 1742432166
  value: ce8f9346-2b44-46af-af55-ab141518e1bb
  key: cluster_id
services:
- name: wenhao
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1733117004
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 1f8b86e3-6352-47a3-ac19-01537a11bec9
  read_timeout: 60000
  retries: 5
  port: 3000
  host: wenhao
- name: grafana
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740320
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 228308d0-521c-45f8-a97d-fbe6966efa3c
  read_timeout: 60000
  retries: 5
  port: 3000
  host: grafana
- name: gitea
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729852876
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  read_timeout: 60000
  retries: 5
  port: 3000
  host: gitea
- enabled: true
  id: 26effc32-8aa5-4d21-b86a-f14add8afc87
  tags: []
  ca_certificates: ~
  write_timeout: 60000
  connect_timeout: 60000
  client_certificate: ~
  read_timeout: 60000
  path: ~
  host: marimo
  port: 8080
  tls_verify: ~
  name: marimo
  tls_verify_depth: ~
  created_at: 1760594820
  updated_at: 1760594820
  retries: 5
  protocol: http
- name: hasura
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740396
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 29bf0a09-9e84-449e-b660-d0c9bf24cab3
  read_timeout: 60000
  retries: 5
  port: 8080
  host: hasura
- name: portainer
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740912
  updated_at: 1758346462
  enabled: true
  protocol: https
  connect_timeout: 60000
  id: 2cd97a86-3853-4001-9f0d-7769dc40d508
  read_timeout: 60000
  retries: 5
  port: 9443
  host: portainer
- name: konga
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729736772
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 2df2f7c3-792a-4a30-b30c-6011e2f3f916
  read_timeout: 60000
  retries: 5
  port: 1337
  host: konga
- name: plugin-frontend
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1749280964
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 33c02e25-db50-4332-a044-094fce5f50b0
  read_timeout: 60000
  retries: 5
  port: 3001
  host: plugin
- name: backend-open-api
  tags: []
  ca_certificates: ~
  path: /open-api/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1734330146
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 4007e6c1-6ccc-4747-9b96-7bb3f5f78b32
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: event-flow-proxy
  tags: []
  ca_certificates: ~
  path: /service-api/supos/proxy/event/flows
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1742976133
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 42684a67-ac8f-48d6-ae2b-62e1cd26f9d8
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: elastic
  tags: []
  ca_certificates: ~
  path: /elastic/home/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740779
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 43a4afb4-cabf-4fd8-a3f4-d8b31408e6dc
  read_timeout: 60000
  retries: 5
  port: 5601
  host: kibana
- name: backend-service-api
  tags: []
  ca_certificates: ~
  path: /service-api/supos/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1733290847
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 4d210886-a694-416d-847d-95fc597f5921
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: gateway
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1732610867
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 5e8458a8-7323-4a25-9efa-4d56dbf5fa5b
  read_timeout: 60000
  retries: 5
  port: 8070
  host: gateway
- name: backend
  tags: []
  ca_certificates: ~
  path: /inter-api/supos/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 300000
  tls_verify_depth: ~
  created_at: 1729740054
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 300000
  id: 5f70fd49-e3d7-4ba8-b152-62eca6ec4455
  read_timeout: 300000
  retries: 5
  port: 8080
  host: backend
- name: minio
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1731396402
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 647747b1-8efe-45b6-b37f-430f5e5139d6
  read_timeout: 60000
  retries: 5
  port: 9001
  host: minio
- name: EventFlow
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1742536198
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  read_timeout: 60000
  retries: 5
  port: 1889
  host: eventflow
- name: backend-files
  tags: []
  ca_certificates: ~
  path: /files
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1741670682
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 75a7373b-4dda-4a49-b1f8-de6ebee4d4c8
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: TDengineUI
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1742538082
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 776a9d80-b936-4e00-b892-3861b19dd954
  read_timeout: 60000
  retries: 5
  port: 3000
  host: TDengineUI
- name: emqx
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740676
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 89c6a77c-e3e6-4149-a0e7-6ace7fd0413b
  read_timeout: 60000
  retries: 5
  port: 18083
  host: emqx
- name: frontend
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
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729738059
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  read_timeout: 60000
  retries: 5
  port: 3000
  host: frontend
- name: eventflow-backend
  tags: []
  ca_certificates: ~
  path: /eventflow-api/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1755572726
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: 94797f49-35e2-428f-b46a-84788fa06fa8
  read_timeout: 60000
  retries: 5
  port: 1888
  host: eventflow
- name: swagger
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1757603297
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: fuxa
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1733732987
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: abdb1947-dbbf-4a85-ac55-50969ac0ed85
  read_timeout: 60000
  retries: 5
  port: 1881
  host: fuxa
- name: node-red-proxy
  tags: []
  ca_certificates: ~
  path: /service-api/supos/proxy/flows
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1730685511
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: b096bcf5-2984-4acc-9bd5-a570a7653fcd
  read_timeout: 60000
  retries: 5
  port: 8080
  host: backend
- name: keycloak
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729740513
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  read_timeout: 60000
  retries: 5
  port: 8080
  host: keycloak
- name: nodered
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1729739239
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: bba8a174-4679-49df-8bf9-ae9285f1e77e
  read_timeout: 60000
  retries: 5
  port: 1880
  host: nodered
- name: chat2db-backend
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1741747300
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: c53d8d7a-962e-4561-accc-ddf6447336df
  read_timeout: 60000
  retries: 5
  port: 10824
  host: chat2db
- name: minio-inter
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1731460027
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: cc3e8a53-82e1-4f6e-bde6-75207ca2f6d3
  read_timeout: 60000
  retries: 5
  port: 9000
  host: minio-inter
- name: chat2db
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1741747189
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: d9ac42d8-9d41-41e3-a80f-4b64cc215bcf
  read_timeout: 60000
  retries: 5
  port: 3001
  host: platform
- name: nodered-backend
  tags: []
  ca_certificates: ~
  path: /nodered-api/
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1744972980
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: df4453f0-063a-4346-a074-3c2f9a388ca7
  read_timeout: 60000
  retries: 5
  port: 1880
  host: nodered
- name: login
  tags: []
  ca_certificates: ~
  path: /realms/supos/protocol/openid-connect/auth
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1732105495
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: e3e88607-311a-4c23-a9c7-bb879efc463e
  read_timeout: 60000
  retries: 5
  port: 8080
  host: keycloak
- name: apm
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1730264819
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  read_timeout: 60000
  retries: 5
  port: 8080
  host: apm
- name: GenerativeUI
  tags: []
  ca_certificates: ~
  path: /
  client_certificate: ~
  tls_verify: ~
  write_timeout: 600000
  tls_verify_depth: ~
  created_at: 1729748350
  updated_at: 1758351345
  enabled: true
  protocol: http
  connect_timeout: 600000
  id: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  read_timeout: 600000
  retries: 5
  port: 4000
  host: copilotkit
- name: mcpclient
  tags: []
  ca_certificates: ~
  path: ~
  client_certificate: ~
  tls_verify: ~
  write_timeout: 60000
  tls_verify_depth: ~
  created_at: 1742543141
  updated_at: 1758346462
  enabled: true
  protocol: http
  connect_timeout: 60000
  id: fbc56015-44d7-42ba-ac9f-abe725bc2478
  read_timeout: 60000
  retries: 5
  port: 3000
  host: mcpclient
routes:
- snis: ~
  service: 43a4afb4-cabf-4fd8-a3f4-d8b31408e6dc
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730253324
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /elastic/home/app/discover#/?_g=(filters:!(),refreshInterval:(pause:!t,value:60000),time:(from:now-15m,to:now))&_a=(columns:!(container.name,message),dataSource:(type:dataView),filters:!(),grid:(columns:(container.name:(width:301))),index:'31d5c920-8e8a-11ef-80e3-d78ebb901f76',interval:auto,query:(language:kuery,query:''),sort:!(!('@timestamp',desc)))
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 050c90d4-b22d-4d44-a159-0d35c949a111
  name: Logs
  response_buffering: true
  tags:
  - ${ENABLE_ELK_MENU}
  - description:menu.desc.logs
  - sort:7
  - parentName:menu.tag.system
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1757902554
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /WebHooks
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: 0d271f38-f4a5-4008-8bb0-33c706460be3
  name: WebHooks
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 776a9d80-b936-4e00-b892-3861b19dd954
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1742538180
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 0e566d44-48f2-442e-a7e9-b288926f3e88
  name: TDengineUI
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1757604045
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /v3/api-docs
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 0f376c28-8c9d-4a4a-8c6b-5258d9d8edef
  name: swagger-config
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 228308d0-521c-45f8-a97d-fbe6966efa3c
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730270842
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /grafana/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 10c662d8-8304-439b-85a0-398ae09d09e8
  name: grafana-inter
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730254514
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /gitea/home/user/login?redirect_to=/gitea/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 15023439-25a9-481d-90a5-894d4d87d3bf
  name: CICD
  response_buffering: true
  tags:
  - description:menu.desc.cicd
  - sort:1
  - parentName:menu.tag.devtools
  - menu
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731311256
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /app-display
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 1dfa8b71-07e8-41a5-a07d-536fa0e295a2
  name: GenApps
  response_buffering: true
  tags:
  - description:menu.desc.genApps
  - sort:2
  - parentName:menu.tag.appspace
  destinations: ~
- snis: ~
  service: 29bf0a09-9e84-449e-b660-d0c9bf24cab3
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1734328715
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /open-api/graphql/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 1fe717e7-804d-473b-b12b-7f5f68b5ec48
  name: open-hasura-graphql
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730265781
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /apps/freezonex-aps/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 209f6169-707c-49f1-86cc-602bc4674b2e
  name: apm-backend-iner
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 89c6a77c-e3e6-4149-a0e7-6ace7fd0413b
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740717
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /emqx/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 215cdb87-6e13-4e2e-804f-9ea1ac2ff417
  name: MqttBroker
  response_buffering: true
  tags:
  - description:menu.desc.mqttBroker
  - parentName:menu.tag.system
  - sort:3
  destinations: ~
- snis: ~
  service: 94797f49-35e2-428f-b46a-84788fa06fa8
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1755572773
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - '/eventflow-api/ '
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 2695543d-3af1-4ca5-b515-82e0fe9b69c7
  name: eventflow-backend
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729748378
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /copilotkit
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 29027b71-49ce-41d1-96fc-14bcb3a2da00
  name: copilotkit
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1748310135
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /plugin-management
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 366d2cfa-8a0e-4a6b-943e-3b5ecd424f73
  name: PluginManagement
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.PluginManagement
  - parentName:menu.tag.system
  - sort:14
  destinations: ~
- snis: ~
  service: 5f70fd49-e3d7-4ba8-b152-62eca6ec4455
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740083
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /inter-api/supos/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 3794799e-0c23-4065-a88d-7a08c46fbaf4
  name: backend
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: b096bcf5-2984-4acc-9bd5-a570a7653fcd
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730685544
  updated_at: 1758346462
  sources: ~
  methods:
  - GET
  hosts: ~
  paths:
  - /nodered/home/flows
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 3f8f8bf7-4d53-4a11-bc6d-d017beda8695
  name: node-red-flows
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1742536226
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /eventflow/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 45a9fec1-6eaa-49c5-ae00-7aa180f1efde
  name: EventFlowBackend
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731311044
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /uns
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 4d89ed56-90bf-490d-ac8c-30def8be2e2c
  name: Namespace
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.dataModeling
  - homeParentName:menu.tag.uns
  - homeIconUrl:homeNamespace
  - sort:2
  destinations: ~
- snis: ~
  service: d9ac42d8-9d41-41e3-a80f-4b64cc215bcf
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741747272
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /chat2db/home/connections
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 4e5336db-053c-4c9f-a773-bb83a2557ffb
  name: DBConnect
  response_buffering: true
  tags:
  - menu
  - parentName:menu.tag.devtools
  - description:menu.desc.dbconnect
  - sort:2
  destinations: ~
- snis: ~
  service: 29bf0a09-9e84-449e-b660-d0c9bf24cab3
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740461
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /hasura/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 4f4b9999-c921-4ca8-9eef-63db46e6b2b5
  name: graphQL
  response_buffering: true
  tags:
  - description:menu.desc.graphQL
  destinations: ~
- sources: ~
  destinations: ~
  snis: ~
  service: 26effc32-8aa5-4d21-b86a-f14add8afc87
  methods: ~
  hosts: ~
  regex_priority: 0
  headers: ~
  paths:
  - /marimo/home/
  id: c8fd090f-2ff9-499c-807b-139f6e5fb976
  https_redirect_status_code: 426
  path_handling: v1
  strip_path: true
  protocols:
  - http
  - https
  preserve_host: false
  name: marimo
  response_buffering: true
  created_at: 1760594842
  updated_at: 1760594842
  request_buffering: true
  tags: ~
- snis: ~
  service: 4d210886-a694-416d-847d-95fc597f5921
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1733290917
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /service-api/supos/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 5a0cdda1-0ac2-4255-a4bc-a11a8b4a00d5
  name: backend-service-api
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: df4453f0-063a-4346-a074-3c2f9a388ca7
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1744972998
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /nodered-api/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 5a9793a9-5085-4830-bba7-a5c053a055a1
  name: nodered-backend
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: fbc56015-44d7-42ba-ac9f-abe725bc2478
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1742543248
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /mcpclient/home
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 6239b86d-4ef1-48d9-b512-7b05306ab705
  name: McpClient
  response_buffering: true
  tags:
  - ${ENABLE_MCP}
  - parentName:menu.tag.appspace
  - description:menu.desc.mcpclient
  - sort:3
  destinations: ~
- snis: ~
  service: 75a7373b-4dda-4a49-b1f8-de6ebee4d4c8
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741671187
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /files
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 670008e2-5811-4d9d-b925-5429d16caa8f
  name: backend-files
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1757902539
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /CodeManagement
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: 6737cbc9-f115-4018-bd0e-90a2342bb002
  name: CodeManagement
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 29bf0a09-9e84-449e-b660-d0c9bf24cab3
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1734328933
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /open-api/restapi/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 6856a277-bcae-43a3-992c-81fc5c605baf
  name: open-hasura-restapi
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 5e8458a8-7323-4a25-9efa-4d56dbf5fa5b
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1732611195
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /gateway
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 6ce6b319-0e38-4bf4-ba1d-4d043a926ba7
  name: gateway
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1730264872
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /apps/freezonex-aps/apsfrontend/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 6f2a8d9e-d76f-46b9-b2d4-27286651433d
  name: apm-frontend-inter
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1758346489
  updated_at: 1758346489
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /Alert
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: 71307652-070e-4ef1-8b30-704d186f3b75
  name: Alert
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: c53d8d7a-962e-4561-accc-ddf6447336df
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741747323
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /chat2db/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 78e5de54-af01-4eff-bb5d-bfbde4f35835
  name: chat2dbapi
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 2df2f7c3-792a-4a30-b30c-6011e2f3f916
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729736896
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /konga/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 79d5e57d-340c-4d18-93ee-a6a8f4a0f212
  name: RoutingManagement
  response_buffering: true
  tags:
  - description:menu.desc.konga
  - sort:1
  - parentName:menu.tag.system
  - menu
  destinations: ~
- snis: ~
  service: 43a4afb4-cabf-4fd8-a3f4-d8b31408e6dc
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740838
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /elastic/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 8154564f-57c9-4de4-a605-a204c46bf15a
  name: elastic-inter
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: f145ba5c-e9aa-48a6-8a10-9ee476010f7f
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1730265586
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /apsfrontend/dashboard
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 833d37e3-05bb-4a1d-992c-0a0e5c19b0dd
  name: apm
  response_buffering: true
  tags:
  - description:menu.desc.apm
  destinations: ~
- snis: ~
  service: 33c02e25-db50-4332-a044-094fce5f50b0
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1749280989
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /plugin/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 84235d0d-bc92-4a86-9413-a66807af7f90
  name: plugin-frontend
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1733209593
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /streams
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 8613d94b-0fd7-4b81-8f9d-ac5df267b7e3
  name: StreamProcessing
  response_buffering: true
  tags:
  - parentName:menu.tag.connections
  destinations: ~
- snis: ~
  service: abdb1947-dbbf-4a85-ac55-50969ac0ed85
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1733733028
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /fuxa/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 86fbf4ad-d446-47e3-a9ba-00ead9113c17
  name: LowCodeTool
  response_buffering: true
  tags:
  - parentName:menu.tag.connections
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1757902544
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /AppManagement
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: 88be7566-488d-4207-b54e-abcb082e2c5e
  name: AppManagement
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1758346494
  updated_at: 1758346494
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /ThemeManagement
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: 89f3d3e9-97fc-48e8-baf0-4751b83992cb
  name: ThemeManagement
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 1f8b86e3-6352-47a3-ac19-01537a11bec9
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1733117175
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 95100aa4-a76b-481c-8485-56eb1f71a529
  name: Ollama
  response_buffering: true
  tags:
  - menuPort:13020
  - menuProtocol:http
  - menuHost:${ENTRANCE_DOMAIN}
  - description:The llm webUI which is  easy to use
  destinations: ~
- snis: ~
  service: 4007e6c1-6ccc-4747-9b96-7bb3f5f78b32
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1734330177
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /open-api/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  name: open-backend-api
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731311377
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /collection-flow
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: a13bb597-9740-4dde-929e-d140c286d869
  name: SourceFlow
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.nodered.flow
  - sort:1
  - homeParentName:menu.tag.uns
  - homeIconUrl:homeSourceFlow
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731311359
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /app-space
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: a2060aa0-88e5-4247-9635-f93438bbdd84
  name: GenerativeUI
  response_buffering: true
  tags:
  - description:menu.desc.generativeUI
  - parentName:menu.tag.appspace
  - sort:2
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741573643
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /workflow
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: a5040934-75dd-40c6-94ea-9497ab2b0579
  name: workflow
  response_buffering: true
  tags:
  - parentName:menu.tag.settings
  destinations: ~
- snis: ~
  service: 647747b1-8efe-45b6-b37f-430f5e5139d6
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731396438
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /minio/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: a6d04fe9-a464-493c-8f3a-4750fdd93a32
  name: objectStorageServer
  response_buffering: true
  tags:
  - description:menu.desc.objectStorageServer
  - sort:200
  destinations: ~
- snis: ~
  service: 228308d0-521c-45f8-a97d-fbe6966efa3c
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730270517
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /grafana/home/dashboards/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: aa02bbb9-6459-43bd-9b65-91d89c8854dd
  name: grafana
  response_buffering: true
  tags:
  - description:menu.desc.grafana
  destinations: ~
- snis: ~
  service: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731473911
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /keycloak/home/auth/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: b610973a-764e-4cef-910e-0794f334e4bd
  name: keycloak-auth
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1742968905
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /EventFlow
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: b8262364-32bf-4422-9d6c-04b97bc8c3a7
  name: EventFlow
  response_buffering: true
  tags:
  - menu
  - homeParentName:menu.tag.uns
  - description:menu.desc.eventflow
  - homeIconUrl:homeEventFlow
  - sort:3
  destinations: ~
- snis: ~
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1732108769
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /supos-login
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: ba7b2df9-b0d8-4b6b-844d-43f935f3181f
  name: login
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: a63d9d5e-e6ae-4493-813e-5a7fa92f322b
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1757603385
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /swagger-ui
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: c010b641-a378-4d8f-b377-08a5ac8c79d1
  name: swagger-ui
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729738250
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: c2dececa-99f4-45e1-9859-01e88352bd58
  name: frontend
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1757902564
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /CollectionGatewayManagement
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: c810eaa9-dbdd-475a-80f6-6cd9a08948a6
  name: CollectionGatewayManagement
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731635333
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /home
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: c90b4b7a-8a51-4f40-b4e2-6c0a40be1b15
  name: Home
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.home
  - sort:1
  destinations: ~
- snis: ~
  service: f8fd7fd2-d8f6-47d7-9c54-aa51a24a68ad
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: false
  request_buffering: true
  created_at: 1758351166
  updated_at: 1758351699
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /open-api/health
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: caeab40c-0c09-4465-b54a-2609391a19e8
  name: health
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: b2a70de2-d5db-4755-b8ba-b205d8fbb680
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740574
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /keycloak/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: d0cea78f-1e0d-4b90-98ea-980a455bf5f5
  name: Authentication
  response_buffering: true
  tags:
  - description:menu.desc.keycloak
  - menu
  - sort:2
  - parentName:menu.tag.system
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1730770040
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /dashboards
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: d2a81d6f-8b2a-4b28-8929-3c51ccd16021
  name: Dashboards
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.dashboards
  - parentName:menu.tag.system
  - sort:5
  destinations: ~
- snis: ~
  service: cc3e8a53-82e1-4f6e-bde6-75207ca2f6d3
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1731460050
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /minio/inter/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: dbb92267-e886-4ee8-b758-a9f9e9af1998
  name: minio-inter
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 24a5cd06-6728-48ff-a7c8-1847c836bd59
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729852903
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /gitea/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: e3459f2d-fcb2-412e-87fc-b098d8906b7e
  name: gitea-inter
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: d9ac42d8-9d41-41e3-a80f-4b64cc215bcf
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741747250
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /chat2db/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: e3ff5178-c156-47e6-90f4-ec2a45b05953
  name: chat2db-inter
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: 42684a67-ac8f-48d6-ae2b-62e1cd26f9d8
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1742976255
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /eventflow/home/flows
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: e8e7fe7d-16ba-415a-8d19-e2c41b76b365
  name: event-node-flows
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1758346485
  updated_at: 1758346485
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /Connection
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: e92d4fec-518a-400b-9094-1a91e989769b
  name: Connection
  response_buffering: true
  tags: ~
  destinations: ~
- snis: ~
  service: d9ac42d8-9d41-41e3-a80f-4b64cc215bcf
  regex_priority: 0
  preserve_host: true
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1741747220
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /chat2db/home/workspace
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: ed108474-cc43-4b66-b92a-5de38f061c6d
  name: SQLEditor
  response_buffering: true
  tags:
  - menu
  - parentName:menu.tag.devtools
  - description:menu.desc.sqledit
  - sort:4
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1734056913
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /advanced-use
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: f0a59836-7eea-45b3-b188-51c45c68f305
  name: AdvancedUse
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.advanceUse
  - sort:9
  - parentName:menu.tag.system
  destinations: ~
- https_redirect_status_code: 426
  strip_path: false
  sources: ~
  preserve_host: false
  request_buffering: true
  destinations: ~
  service: 6a3bcaba-7ba3-4c3b-b5c2-1a8ecbc571ac
  methods: ~
  created_at: 1760593302
  updated_at: 1760593714
  headers: ~
  id: f3121c84-e787-4eb6-a054-edbc56df4d0e
  paths:
  - /dashboard
  - /dashboard/
  path_handling: v1
  response_buffering: true
  snis: ~
  hosts: ~
  regex_priority: 0
  tags: ~
  protocols:
  - http
  - https
  name: dashboard2
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1732696995
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /account-management
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: f610002e-8e5c-424a-9fdc-e86d754b3142
  name: UserManagement
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.account
  - parentName:menu.tag.system
  - sort:4
  destinations: ~
- snis: ~
  service: bba8a174-4679-49df-8bf9-ae9285f1e77e
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729739343
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /nodered/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: f8bccdbf-46fd-4dd7-aafe-d1a8e13c3450
  name: NodeRed
  response_buffering: true
  tags: []
  destinations: ~
- snis: ~
  service: 2cd97a86-3853-4001-9f0d-7769dc40d508
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1729740933
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /portainer/home/
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: fb690c3a-e09e-4fcf-aad1-5d357c1938ec
  name: ContainerManagement
  response_buffering: true
  tags:
  - description:menu.desc.dockerMgmt
  - ${ENABLE_PORTAINER}
  - sort:6
  - parentName:menu.tag.system
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1733822464
  updated_at: 1758346462
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /aboutus
  path_handling: v1
  headers: ~
  protocols:
  - http
  - https
  id: fb8efe62-0030-4540-a521-a3b4215b4f59
  name: AboutUs
  response_buffering: true
  tags:
  - menu
  - description:menu.desc.aboutus
  - parentName:menu.tag.system
  - sort:8
  destinations: ~
- snis: ~
  service: 8e081976-8223-4494-9b4c-0aa5a441bdd5
  regex_priority: 0
  preserve_host: false
  https_redirect_status_code: 426
  strip_path: true
  request_buffering: true
  created_at: 1758346492
  updated_at: 1758346492
  sources: ~
  methods: ~
  hosts: ~
  paths:
  - /OpenData
  path_handling: v0
  headers: ~
  protocols:
  - http
  - https
  id: fd23d489-34ad-4c12-9e72-10e43318866d
  name: OpenData
  response_buffering: true
  tags: ~
  destinations: ~
plugins:
- created_at: 1733799127
  name: supos-auth-checker
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 1845ee75-d704-40e1-a8b0-aa2baaf9d71b
  enabled: ${KONG_AUTH_ENABLED}
  route: ~
  config:
    forbidden_url: /403
    whitelist_paths:
    - ^/inter-api/supos/auth.*$
    - ^/inter-api/supos/systemConfig.*$
    - ^/inter-api/supos/theme/getConfig.*$
    - ^/$
    - ^/assets.*$
    - ^/locale.*$
    - ^/logo.*$
    - ^/gitea.*git.*$
    - ^/supos-login.*$
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
    enable_deny_check: true
    login_url: ${BASE_URL}/keycloak/home/auth/realms/supos/protocol/openid-connect/auth?client_id=supos&redirect_uri=${BASE_URL}/inter-api/supos/auth/token&response_type=code&scope=openid
    enable_resource_check: true
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734330234
  name: key-auth
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 2285421f-56e3-4510-be12-69fa1040d810
  enabled: true
  route: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  config:
    key_in_query: true
    key_in_body: false
    realm: ~
    anonymous: ~
    key_names:
    - apikey
    hide_credentials: false
    run_on_preflight: true
    key_in_header: true
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734073535
  name: request-transformer
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 46bda5cf-63ea-401f-9f06-b9e024aa5597
  enabled: true
  route: ~
  config:
    remove:
      body: []
      querystring: []
      headers: []
    append:
      body: []
      querystring:
      - client_id:supos
      - response_type:code
      - scope:openid
      - redirect_uri:${BASE_URL}/inter-api/supos/auth/token
      headers: []
    add:
      body: []
      querystring: []
      headers: []
    http_method: ~
    replace:
      body: []
      headers: []
      uri: ~
      querystring: []
    rename:
      body: []
      querystring: []
      headers: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734329742
  name: key-auth
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 4e339a48-e193-4bc7-99ba-94f2e1e18a8f
  enabled: true
  route: 1fe717e7-804d-473b-b12b-7f5f68b5ec48
  config:
    key_in_query: true
    key_in_body: false
    realm: ~
    anonymous: ~
    key_names:
    - apikey
    hide_credentials: false
    run_on_preflight: true
    key_in_header: true
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734073759
  name: response-transformer
  service: e3e88607-311a-4c23-a9c7-bb879efc463e
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 73285cac-cf4e-4368-bf8f-d3285f9686d4
  enabled: true
  route: ~
  config:
    remove:
      headers:
      - x-frame-options
      - Content-Security-Policy
      json: []
    append:
      headers:
      - X-Frame-Options:DENY
      - Content-Security-Policy:frame-ancestors 'none'
      json_types: []
      json: []
    add:
      headers: []
      json_types: []
      json: []
    replace:
      headers: []
      json_types: []
      json: []
    rename:
      headers: []
      json: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734329166
  name: key-auth
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: 9e4c8423-76d1-498d-8883-b9b18bdbf161
  enabled: true
  route: 6856a277-bcae-43a3-992c-81fc5c605baf
  config:
    key_in_query: true
    key_in_body: false
    realm: ~
    anonymous: ~
    key_names:
    - apikey
    hide_credentials: false
    run_on_preflight: true
    key_in_header: true
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1757987272
  name: cors
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: ad24951e-c659-4588-9138-542e4586e790
  enabled: true
  route: 9df937e7-2ffb-49f4-b60b-4bb5b551419a
  config:
    headers: ~
    credentials: false
    preflight_continue: false
    origins:
    - '*'
    private_network: false
    exposed_headers: ~
    methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE
    - OPTIONS
    max_age: ~
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1731400906
  name: response-transformer
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: b5722a76-60b9-483d-90ac-e5de173264e2
  enabled: true
  route: a6d04fe9-a464-493c-8f3a-4750fdd93a32
  config:
    remove:
      headers:
      - X-Frame-Options
      json: []
    append:
      headers: []
      json_types: []
      json: []
    add:
      headers:
      - X-Frame-Options:SAMEORIGIN
      json_types: []
      json: []
    replace:
      headers: []
      json_types: []
      json: []
    rename:
      headers: []
      json: []
  protocols:
  - grpc
  - grpcs
  - http
  - https
- created_at: 1734014838
  name: supos-url-transformer
  service: ~
  tags: ~
  updated_at: 1758346462
  instance_name: ~
  consumer: ~
  id: f45e7fd0-74e5-4b36-b265-5df962eb1b58
  enabled: true
  route: ba7b2df9-b0d8-4b6b-844d-43f935f3181f
  config:
    home_url: /?isLogin=true
  protocols:
  - grpc
  - grpcs
  - http
  - https
upstreams:
- name: apm
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 00f62da1-73b3-4308-8ae5-856268388aa9
  created_at: 1730264799
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: minio-inter
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 09fa87e7-1cf1-4fc3-9ad6-6861572a8ef4
  created_at: 1731459961
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: platform
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 14bdf793-231d-439e-8d67-3e47a3e7da05
  created_at: 1741754738
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: keycloak
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 14ee49e7-f9bf-4234-a48f-7b7df7dda0ea
  created_at: 1729739799
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: nodered
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 220c1252-a48b-4e7e-af42-b134f316ed16
  created_at: 1729739256
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: konga
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 2a0fa8a9-98a4-4456-9d2c-faba35b54882
  created_at: 1729737217
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: grafana
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 2e9455f1-241d-4a1e-840e-fea66dd9aa09
  created_at: 1729739617
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: plugin
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 3c840281-aeb3-4969-8ea4-cc12d45423a1
  created_at: 1749280889
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: portainer
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 3f51d2ee-7609-4f71-b6ec-8ded27b417a2
  created_at: 1729739903
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: kibana
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 420478e2-bdc8-49ec-ba0e-cc4cfd41afc8
  created_at: 1729739883
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: gitea
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 47b5a73e-51c6-4ace-b507-40458125c0a6
  created_at: 1729852910
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: copilotkit
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: 747a43cc-42c4-457a-abac-518c2fe537b3
  created_at: 1729739930
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: hasura
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: ab6fb6a2-97b9-45e7-b508-46df2df5203e
  created_at: 1729739734
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: frontend
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: c0645f32-b63d-49ba-a723-5cbef9014c99
  created_at: 1729737694
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: minio
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: c6844f11-b711-4f5f-a2d4-4516995790c5
  created_at: 1731396351
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: backend
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: e6729ab9-c894-4963-bbdc-dfae17c88096
  created_at: 1729739594
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: emqx
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: fa9fc31a-7d71-43bd-a453-c21704b71ac6
  created_at: 1729739828
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
- name: fuxa
  hash_on: none
  healthchecks:
    passive:
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
      unhealthy:
        http_statuses:
        - 429
        - 500
        - 503
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      type: http
    active:
      headers: ~
      unhealthy:
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
        tcp_failures: 0
        timeouts: 0
        http_failures: 0
      healthy:
        interval: 0
        http_statuses:
        - 200
        - 302
        successes: 0
      type: http
      timeout: 1
      http_path: /
      https_sni: ~
      https_verify_certificate: true
      concurrency: 10
    threshold: 0
  tags: ~
  hash_on_header: ~
  slots: 1000
  hash_on_query_arg: ~
  hash_fallback: none
  hash_on_uri_capture: ~
  hash_fallback_query_arg: ~
  hash_fallback_uri_capture: ~
  id: fc85cfef-6e90-4592-a022-730f5d84ed99
  created_at: 1733536118
  hash_on_cookie: ~
  algorithm: round-robin
  host_header: ~
  client_certificate: ~
  use_srv_name: false
  updated_at: 1758346462
  hash_fallback_header: ~
  hash_on_cookie_path: /
targets:
- created_at: 1730264807.829
  updated_at: 1758346462.943
  upstream: 00f62da1-73b3-4308-8ae5-856268388aa9
  tags: ~
  id: 0be879dd-7ba3-42e3-96a8-5be0c8dcadfe
  weight: 100
  target: apm:8080
- created_at: 1729739629.259
  updated_at: 1758346462.957
  upstream: 2e9455f1-241d-4a1e-840e-fea66dd9aa09
  tags: ~
  id: 0c66a055-3e1d-4e77-aedd-6de30939a5a2
  weight: 100
  target: grafana:3000
- created_at: 1729852918.487
  updated_at: 1758346462.96
  upstream: 47b5a73e-51c6-4ace-b507-40458125c0a6
  tags: ~
  id: 29a2c68c-7501-49b1-aef1-56e1be3c2d74
  weight: 100
  target: gitea:3000
- created_at: 1733536131.451
  updated_at: 1758346462.953
  upstream: fc85cfef-6e90-4592-a022-730f5d84ed99
  tags: ~
  id: 2cc80d27-504d-4c1a-afaf-646f6e432fd6
  weight: 100
  target: fuxa:1881
- created_at: 1731397614.344
  updated_at: 1758346462.954
  upstream: c6844f11-b711-4f5f-a2d4-4516995790c5
  tags: ~
  id: 2f312fcf-baa6-4de0-84aa-0779713ae279
  weight: 100
  target: minio:9001
- created_at: 1729737239.883
  updated_at: 1758346462.966
  upstream: 2a0fa8a9-98a4-4456-9d2c-faba35b54882
  tags: ~
  id: 38ae0eeb-fa36-40a1-8c72-aaa43e7b9ceb
  weight: 100
  target: konga:1337
- created_at: 1729739844.316
  updated_at: 1758346462.965
  upstream: fa9fc31a-7d71-43bd-a453-c21704b71ac6
  tags: ~
  id: 3ec4626c-673a-4cc1-86fa-10f3b1bdc537
  weight: 100
  target: emqx:18083
- created_at: 1729739274.174
  updated_at: 1758346462.958
  upstream: 220c1252-a48b-4e7e-af42-b134f316ed16
  tags: ~
  id: 5a923a6d-0661-43ad-a668-e89f4980e5a6
  weight: 100
  target: nodered:1880
- created_at: 1729739915.095
  updated_at: 1758346462.963
  upstream: 3f51d2ee-7609-4f71-b6ec-8ded27b417a2
  tags: ~
  id: 78410225-c503-4d2c-885d-c38a72244aa0
  weight: 100
  target: portainer:9443
- created_at: 1741754798.218
  updated_at: 1758346462.945
  upstream: 14bdf793-231d-439e-8d67-3e47a3e7da05
  tags: ~
  id: 7c5e8a1d-6022-42b6-83b7-b1f39739b9df
  weight: 100
  target: frontend:3001
- created_at: 1729739753.711
  updated_at: 1758346462.946
  upstream: ab6fb6a2-97b9-45e7-b508-46df2df5203e
  tags: ~
  id: 8d09a927-4def-479d-8ed9-0473be2a281a
  weight: 100
  target: hasura:8080
- created_at: 1729832116.752
  updated_at: 1758346462.947
  upstream: 14ee49e7-f9bf-4234-a48f-7b7df7dda0ea
  tags: ~
  id: a4cdac00-4421-4a49-b2ba-b2f720f4998d
  weight: 100
  target: keycloak:8080
- created_at: 1729739894.022
  updated_at: 1758346462.948
  upstream: 420478e2-bdc8-49ec-ba0e-cc4cfd41afc8
  tags: ~
  id: afe138dc-2d5e-4329-80ff-a0154b634f42
  weight: 100
  target: kibana:5601
- created_at: 1729739606.748
  updated_at: 1758346462.95
  upstream: e6729ab9-c894-4963-bbdc-dfae17c88096
  tags: ~
  id: e401fd6b-6a87-41be-a501-6456afeb14d3
  weight: 100
  target: backend:8080
- created_at: 1731459990.339
  updated_at: 1758346462.955
  upstream: 09fa87e7-1cf1-4fc3-9ad6-6861572a8ef4
  tags: ~
  id: e5f0e5f7-8b6a-476b-ac72-2e4a32bab571
  weight: 100
  target: minio:9000
- created_at: 1729914708.061
  updated_at: 1758346462.962
  upstream: c0645f32-b63d-49ba-a723-5cbef9014c99
  tags: ~
  id: ef3e3b44-a490-40e8-afac-2b9e2d202643
  weight: 100
  target: frontend:3000
- created_at: 1729914781.603
  updated_at: 1758346462.951
  upstream: 747a43cc-42c4-457a-abac-518c2fe537b3
  tags: ~
  id: f4e24691-4181-42a6-b32c-6c7e8cd92077
  weight: 100
  target: frontend:4000
- created_at: 1749280900.4
  updated_at: 1758346462.949
  upstream: 3c840281-aeb3-4969-8ea4-cc12d45423a1
  tags: ~
  id: f6a73dfa-96f1-4bd4-969f-a46592a5ac49
  weight: 100
  target: frontend:3002
keyauth_credentials:
- created_at: 1758346492
  id: 5bfc410b-ec4b-4371-a6dd-6775e38c5dc1
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 0b7dc033e36f4a1492ac8562885cac27
  tags: ~
  ttl: ~
- created_at: 1734329245
  id: 6b9443ae-73f0-4db6-af00-4f1e3a415dbb
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 4174348a-9222-4e81-b33e-5d72d2fd7f1e
  tags: ~
  ttl: ~
- created_at: 1749280821
  id: 6b979102-3576-4a98-b06e-ea72d25249d0
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: d763a82bd8154b58bd29d9cd141dcab0
  tags: ~
  ttl: ~
- created_at: 1757578203
  id: 7e12e6c7-e28b-4e5f-876b-3f49341ae403
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: d3672e250b634a4b8602d88f1ae55a81
  tags: ~
  ttl: ~
- created_at: 1749260555
  id: 90d62a81-2634-45a4-8384-6e155876db84
  consumer: 59d1ef15-24a5-4373-b957-e8192c15ff6e
  key: 33d911ac009240279e01b4a95655a0ad
  tags: ~
  ttl: ~
