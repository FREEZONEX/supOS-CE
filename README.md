# TIER0: The Unified Namespace Software


[![Static Badge](https://img.shields.io/badge/Try%20Tier0-Live%20Demo-blue?style=flat&logo=rocket&logoColor=red)](https://tier0.app/trial)
[![Docs](https://img.shields.io/badge/docs-available-brightgreen?style=flat&logo=readthedocs)](https://tier0edge.vercel.app/)
[![License](https://img.shields.io/badge/License-Apache_2.0-yellow?style=flat&logo=open-source-initiative)](./LICENSE)

**Tier0** is an open-source industrial data integration platform built on the **Unified Namespace (UNS)** methodology and powered by production-grade open-source technologies.

---

## Hardware Requirements

|             | Minimum Requirement                  | Recommended Requirement                       |
|-------------|--------------------------------------|-----------------------------------------------|
| CPU         | 4 cores                              | 8 cores                                       |
| Memory      | 8 GB                                 | 16 GB                                         |
| Disk        | 100 GB, 1000 IOPS (30% random write)      | 1 TB, 2000 IOPS (30% random write)        |
| Browser     | Chrome 89, Edge 89, Firefox 89, Safari 15 | Chrome 89, Edge 89, Firefox 89, Safari 15 |

## Deployment
> For detailed guides and advanced examples, see the <a href="https://tier0edge.tech/">Tier0 Community Docs</a>.
### 1.Linux
#### 1.1 Operating Environment
- **Operating System**: Currently tested on Ubuntu Server 24.04 with Docker. We welcome feedback on other OS distributions.
- **Docker**: We assume you have Docker (with `docker compose` and `buildx`) installed. Our tested versions:
  - Docker Engine - Community: 27.4.0
  - Docker Buildx: v0.19.2
  - Docker Compose: v2.31.0
  - containerd: 1.7.24

#### 1.2 Installing Tier0
1. Clone the project.
   ```bash
   git clone <this repo>
   ```
2. Navigate to the `Tier0` directory and edit environment variables in the `.env` file.
   ```bash
   cd Tier0-Edge/deploy
   vi .env.default
   ```
  - Update `VOLUMES_PATH` (directory for storing project data).
  - Update `ENTRANCE_DOMAIN` (frontend entry domain/IP address).
  - Modify other variables as needed.
3. Install Tier0.
   ```bash
   bash bin/install.sh
   ```
### 2.Windows
#### 2.1 Operating Environment
- Install the latest version of **Docker Desktop** and **Git** on Windows 10 or Windows 11.
- It is recommended to perform all operations in **Git Bash** for better compatibility.
#### 2.2 Installing Tier0
1. Clone the project using **Git Bash**.
   ```bash
   git clone <this repo>
   ```
2. Navigate to the `Tier0` directory and edit environment variables in the `.env` file.
   ```bash
   cd Tier0-Edge/deploy
   vi .env.default
   ```
  - Update `OS_PLATFORM_TYPE` = windows
  - Update `VOLUMES_PATH` (directory for storing project data).
  - Update `ENTRANCE_DOMAIN`  (Do not use 127.0.0.1 or localhost, otherwise login and authentication functions **will NOT** work.)
  - Modify other variables as required by the system.
3. Install Tier0.
   ```bash
   bash bin/install.sh
   ```
### 3. Access the Platform
1. Visit `http://<YOUR-DOMAIN>:<YOUR-PORT>` in your browser (based on ENTRANCE_DOMAIN and ENTRANCE_PORT in `.env`).
2. Sign in to Tier0 with default account and password: `tier0/tier0`.
---

## Important Startup Operations
### 1. UNS Data Model Creation
1. Send the template JSON from **Import** to an LLM, and use a similar prompt.
    ```
    Generate a UNS model used for xx in xx plant, including xx equipment and data sources based on the template.
    ```
2. Import the generated result in UNS.

### 2. Model Data Source Connection
> Connect real data to make models alive.
1. Use nodes based on the data source type to build a flow, and end it with an `mqtt out` node.
  > Install nodes from Node-RED community for your requirements. Official nodes often start with `node-red-contrib`.

2. Make sure the **Server** of the `mqtt out` node is set to the UNS broker, and topic is a model from **UNS**.
  > The UNS broker has the same name as that of the flow.
---
## License
This project is licensed under the [Apache 2.0 License](./LICENSE).

## Support & Contact
- 📖 [Documentation](https://suposcommunity.vercel.app)
- 🐞 [GitHub Issues](./issues)

## Contributors
We gratefully acknowledge the following individuals for their contributions to Tier0:

**Wenhao Yu**, **Liebo**, **Weipeng Dong**, **Kangxi**, **Lifang Sun**, **Minghe Zhuang**,  
**Wangji Xin**, **Fayue Zheng & Yue Yang**, **Yanqiu Liu**, **Dongdong An**, **Jianan Zhu**
