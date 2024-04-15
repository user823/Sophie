# Sophie

English | [中文](README_cn.md)

[![Static Badge](https://img.shields.io/badge/release-1.0.0-green)](https://github.com/user823/Sophie/releases)
[![Static Badge](https://img.shields.io/badge/website-sophie-green)](https://49.234.183.205/)
[![Static Badge](https://img.shields.io/badge/license-Apache--2.0-green)](https://github.com/user823/Sophie/blob/main/LICENSE)

![sophie](docs/images/sophie.jpg)

Sophie is a front-end and back-end separated permission management system designed based on Hertz + Kitex + Element UI. Individuals and enterprises can quickly develop based on this system.

## Features

- Adopting a front-end and back-end separation model, allowing independent deployment and modification of front-end and back-end applications, with strong flexibility.
- The back-end adopts Hertz and Kitex frameworks from [ByteDance](https://www.cloudwego.io/), which provide numerous extension interfaces, ensuring high performance while maintaining strong scalability.
- The back-end adopts a gateway architecture, with microservices of various components exposed to the outside world through a unified sophie-gateway.
- The back-end follows the RESTful API design specification.
- Utilizing an RBAC-based access control model.
- Equipped with comprehensive subsystems such as caching and log aggregation.
- Providing distributed task scheduling functionality.
- The project includes rich documentation and testing, making it easy to understand.

## Architecture

![architecture](docs/images/architecture.png)

## Built-in Functions

1. User Management: Configure attributes and status of users.
2. Department Management: Manage attributes and status of organization's various hierarchical levels using a tree structure.
3. Position Management: Manage positions defined within the organization and their status.
4. Menu Management: Roles with permissions can edit the system menus.
5. Role Management: Manage user permissions based on roles, organizations can define permission roles internally.
6. Dictionary Management: Dictionaries consist of dictionary names, types, and statuses, and dictionary types and their value ranges are defined internally by the organization.
7. Parameter Management: Parameters are represented by key-value pairs, managing system runtime parameter settings, and modifying system behavior by modifying parameters.
8. Notice Bulletin: Publish and maintain system notice bulletin information.
9. Operation Log: Record and query normal operation logs of the system; record and query system exception information logs.
10. Login Log: Record and query system login logs.
11. Online Users: Monitor the active user status in the current system.
12. Scheduled Tasks: Online (add, modify, delete) task scheduling includes execution result logs.
13. Code Generation: Support generation of front-end and back-end code (java, html, xml, sql), support CRUD download.
14. System Interface: Automatically generate related API interface documents based on business code.
15. Service Monitoring: Administrators can monitor service invocation chains, service status of various components, etc.
16. Form Builder: Users define pages by dragging components and setting component properties.

## Demo

- admin/admin123

Demo: [https://49.234.183.205/](https://49.234.183.205/)

## Project Details

- [Requirement Analysis](docs/devel/requirements_analysis.md)
- [Technical Selection](docs/devel/technology_selection.md)
- [System Architecture](docs/devel/architecture.md)
- [Project Structure](docs/guide/project_structure.md)
- [Deployment Guide](docs/guide/deployment.md)
- [App Startup Configuration and Startup Process](docs/guide/app.md)

## Related Projects

- [Sophie-ui](https://github.com/user823/Sophie-ui)

## Open Source License

Sophie is licensed under the [Apache License 2.0](LICENSE)

## Contact Information

- Email: hq869860837@163.com
- Issues: [Issues](https://github.com/user823/Sophie/issues)
- QQ: 869860837
