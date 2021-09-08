---
title: "Google Cloud Platform"
date: 2019-05-29T17:06:21-07:00
draft: false
weight: 2
tags: ["python", "django", "appengine", "gce", "gcp", "google", "compute"]
---

![](/images/gcp-logo.png)

- [Introduction](#introduction)
- [Requirements](#requirements)
- [Addition to your code](#addition-to-your-code)
- [References](#references)


### Introduction

This guide will help you add [sqlcommenter](/introduction) to your Django applications running on [Google Cloud Platform (GCP)](https://cloud.google.com)

### Requirements

Steps|Resource
---|---
Django on GCP|https://cloud.google.com/python/django/
google-cloud-sqlcommenter|https://pypi.org/project/google-cloud-sqlcommenter
Django 2.X|https://docs.djangoproject.com/en/stable/faq/install
Python 3.X|https://www.python.org/downloads/

### Addition to your code

Firstly, please install [google-cloud-sqlcommenter](/python/django#installation).

For any Django deployment, we can just edit our settings.py file and update the `MIDDLEWARE` section as per:

```python
MIDDLEWARE = [
  'google.cloud.sqlcommenter.django.middleware.SqlCommenter',
  ...
]
```

{{% notice tip %}}
If any middleware execute database queries (that you'd like commented by SqlCommenter), those middleware MUST appear after
'google.cloud.sqlcommenter.django.middleware.SqlCommenter'
{{%/ notice %}}

### References

Resource|URL
---|---
Running Django on GCP|https://cloud.google.com/python/django/
Installing Django middleware|[/python/django#installation](/python/django#installation)
