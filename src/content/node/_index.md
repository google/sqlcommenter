---
title: "Node.js"
date: 2019-05-31T18:16:11-06:00
draft: false
weight: 40
---

- [Introduction](#introduction)
- [Integrations](#integrations)
- [Installation](#install)
    - [From source](#install-from-source)
    - [Verify installation](#verify-installation)

### Introduction

sqlcommenter is a suite of plugins/middleware/wrappers to augment SQL statements from ORMs/Querybuilders with comments that can be used later to correlate user code with SQL statements.

### Integrations

sqlcommenter-nodejs provides support for the following:

{{<card-vendor href="/node/knex" src="/images/knex-logo.png">}}
{{<card-vendor href="/node/sequelize" src="/images/sequelize-logo.png">}}
{{<card-vendor href="/node/express" src="/images/express_js-logo.png">}}

### <a name="install"></a> Installing it
sqlcommenter-nodejs can installed in a couple of ways:

#### <a name="install-from-source"></a> From source

The first step is to clone the repository. This can be done with [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) by running:
{{<highlight shell>}}
git clone https://github.com/google/sqlcommenter.git
{{</highlight>}}
Inspect the source code and note the path to the package you want installed.

```shell 
sqlcommenter/nodejs/sqlcommenter-nodejs
└── packages
    ├── knex
    │   ├── index.js
    │   ├── package.json
    │   ├── test
    │   └── ...
    └── sequelize
        ├── index.js
        ├── package.json
        ├── test
        └── ...
```
Each folder in the `packages` directory can be installed by running 

{{<highlight shell>}}
npm install <path/to/package>
{{</highlight>}}

for example to install `@google-cloud/sqlcommenter-knex` in a given location, run `npm install /path/to/sqlcommenter-nodejs/packages/knex`. Same for every package(folder) in the `packages` directory.
```shell
# install 
> npm install /path/to/sqlcommenter-nodejs/packages/knex

+ @google-cloud/sqlcommenter-knex@0.0.1
```

#### <a name="verify-installation"></a> Verify Installation
If package is properly installed, running `npm list <package-name>` will output details of the package. Let's verify the installation of `@google-cloud/sqlcommenter-knex` below:
```shell
# verify
> npm list @google-cloud/sqlcommenter-knex

project@0.0.0 path/to/project
└── @google-cloud/sqlcommenter-knex@0.0.1  -> /path/to/sqlcommenter-nodejs/packages/knex
```
Inspecting the `package.json` file after installation should also show the installed pacakge.
