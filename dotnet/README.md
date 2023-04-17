# sqlcommenter for dotnet

- [Introduction](#Introduction)
- [Integrations](#Integrations)
  - [ORMs](#Integrations)
  - [Frameworks](#Integrations)
- [Usage](#Usage)

## Introduction

NuGet packages that add meta info to your SQL queries as comments.

## Integrations

### ORMs

- [X] EntityFrameworkCore

### Frameworks

- [X] AspNetCore
- [X] AspNetCore gRPC

## Usage

### AspNetCore with EntityFramework Core

Add the NuGet PackageReferences:
```c#
<PackageReference Include="Google.Cloud.SqlCommenter.AspNetCore" Version="X.X.X"
<PackageReference Include="Google.Cloud.SqlCommenter.EntityFrameworkCore" Version="X.X.X"
```

Use the sql commenter middleware in the Startup:
```c#
public void Configure(IApplicationBuilder app)
{
  // other initializations

  app.UseRouting();
  
  app.UseSqlCommenter(); // needs to be after the UseRouting call, but before the UseEndpoints call.
  
  app.UseEndpoints(endpoints => ...);
}
```

Use sql commenter on the DbContext:
```c#
public void ConfigureServices(IserviceCollection services)
{
  // Other DI initializations

  services.AddDbContext<MyDb>(options => options
    .UseSqlite() // use your database configuration
    .UseSqlCommenter());
}
```

#### gRPC

To use sql commenter for grpc services add the following NuGet PackageReference:
```c#
<PackageReference Include="Google.Cloud.SqlCommenter.AspNetCore.Grpc" Version="X.X.X"
```

And enable the middleware in the Startup:
```c#
public void Configure(IApplicationBuilder app)
{
  // other initializations

  app.UseRouting();
  
  // needs to be after the UseRouting call, but before the UseEndpoints call.
  app.UseGrpcSqlCommenter();
  
  app.UseEndpoints(endpoints => ...);
}
```
