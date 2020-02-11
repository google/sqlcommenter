# sqlcommenter

sqlcommenter is a suite of plugins/middleware/wrappers to augment SQL statements from ORMs/Querybuilders
with comments that can be used later to correlate user code with SQL statements.

It supports Node v6 and above to use ES6 features.

### Supported frameworks:

- Sequelize
- Knex.js

### Installation

Go into either of the packages [packages/knex](./packages/knex) or [packages/sequelize](./packages/sequelize)
and then you can run respectively

Middleware|Command|URL
---|---|---
Knex.js|`npm install @google-cloud/sqlcommenter-knex`|https://www.npmjs.com/package/@google-cloud/sqlcommenter-knex
Sequelize.js|`npm install @google-cloud/sqlcommenter-sequelize`|https://www.npmjs.com/package/@google-cloud/sqlcommenter-sequelize
