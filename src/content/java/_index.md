---
title: "Java"
date: 2019-05-29T08:46:39-06:00
draft: false
weight: 40
---

- [Introduction](#introduction)
- [Integrations](#integrations)
- [Installing it](#install)
    - [From source](#install-from-source)
        - [Building it](#building-it)
    - [Verify installation](#verify-installation)
    - [Tests](#tests)

### Introduction
sqlcommenter-java is the implementation of [sqlcommenter](/) in the Java programming language.


### Integrations

sqlcommenter-java provides support for the following plugins/ORMs:

{{<card-vendor href="/java/hibernate" src="/images/hibernate-logo.svg">}}
{{<card-vendor href="/java/spring" src="/images/spring-logo.png">}}

### <a name="install"></a> Installing it
sqlcommenter-java can installed in a couple of ways:

#### <a name="install-from-source"></a> From source

Please visit [source page on Github](https://github.com/google/sqlcommenter/tree/master/java/sqlcommenter-java)

#### Building it

Next, after changing directories into `java/sqlcommenter-java`, run `./gradlew install`
which should produce should output
```shell
$ ./gradlew install

BUILD SUCCESSFUL in 1s
7 actionable tasks: 1 executed, 6 up-to-date
```

#### Verify installation

sqlcommenter-java if properly installed should appear in the directory `$HOME/.m2/integrations/repository/io`.

The following should be your directory structure:
```shell
~/.m2/repository/io
└── com
    └── google
        └── cloud
            └── sqlcommenter
                ├── 0.0.1
                │   ├── sqlcommenter-java-0.0.1-javadoc.jar
                │   ├── sqlcommenter-java-0.0.1-javadoc.jar.asc
                │   ├── sqlcommenter-java-0.0.1-sources.jar
                │   ├── sqlcommenter-java-0.0.1-sources.jar.asc
                │   ├── sqlcommenter-java-0.0.1.jar
                │   ├── sqlcommenter-java-0.0.1.jar.asc
                │   └── sqlcommenter-java-0.0.1.pom
                └── maven-metadata-local.xml
```

and then in your programs that use Maven, when building packages, please do
```shell
mvn install -nsu
```
to use look up local packages.

#### Tests

Tests can be run by
```shell
$ ./gradlew test
```
