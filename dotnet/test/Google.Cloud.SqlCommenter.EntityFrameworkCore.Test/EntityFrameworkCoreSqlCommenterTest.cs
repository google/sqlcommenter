using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Xunit;

namespace Google.Cloud.SqlCommenter.EntityFrameworkCore.Test
{
    // cant use ToQueryString since the interceptor is not called
    public class EntityFrameworkCoreSqlCommenterTest : IDisposable
    {
        private readonly SqliteConnection _connection;

        public EntityFrameworkCoreSqlCommenterTest()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
        }

        public void Dispose() => _connection.Dispose();

        [Fact]
        public async Task EnabledByDefault()
        {
            var logs = new List<string>();
            await using var context = await NewContext(builder => builder.UseSqlCommenter(opts => opts.ExtraValues.Add("application", "my-app")), logs.Add);
            SetState();
            await context.TestEntities.ToListAsync();

            var expectedLog =
@"      SELECT ""t"".""Id"", ""t"".""Name""
      FROM ""TestEntities"" AS ""t"" /*action='my-action',application='my-app',controller='my-controller',framework='my-framework+ef-core',route='my-route'*/";
            Assert.EndsWith(expectedLog, logs.Last());
        }

        [Fact]
        public async Task EnabledByDefault_ExplicitDisable()
        {
            var logs = new List<string>();
            await using var context = await NewContext(builder => builder.UseSqlCommenter(opts => opts.ExtraValues.Add("application", "my-app")), logs.Add);
            SetState();
            await context.TestEntities.DisableSqlCommenter().ToListAsync();

            var expectedLog =
                @"      SELECT ""t"".""Id"", ""t"".""Name""
      FROM ""TestEntities"" AS ""t""";
            Assert.EndsWith(expectedLog, logs.Last());
        }

        [Fact]
        public async Task DisabledByDefault()
        {
            var logs = new List<string>();
            await using var context = await NewContext(builder => builder.UseSqlCommenter(opts =>
            {
                opts.ExtraValues.Add("application", "my-app");
                opts.DefaultEnable = false;
            }), logs.Add);
            SetState();
            await context.TestEntities.ToListAsync();

            var expectedLog =
                @"      SELECT ""t"".""Id"", ""t"".""Name""
      FROM ""TestEntities"" AS ""t""";
            Assert.EndsWith(expectedLog, logs.Last());
        }

        [Fact]
        public async Task DisabledByDefault_ExplicitEnable()
        {
            var logs = new List<string>();
            await using var context = await NewContext(builder => builder.UseSqlCommenter(opts =>
            {
                opts.ExtraValues.Add("application", "my-app");
                opts.DefaultEnable = false;
            }), logs.Add);
            SetState();
            await context.TestEntities.EnableSqlCommenter().ToListAsync();

            var expectedLog =
                @"      SELECT ""t"".""Id"", ""t"".""Name""
      FROM ""TestEntities"" AS ""t"" /*action='my-action',application='my-app',controller='my-controller',framework='my-framework+ef-core',route='my-route'*/";
            Assert.EndsWith(expectedLog, logs.Last());
        }

        [Fact]
        public async Task DisabledWithExistingComment()
        {
            var logs = new List<string>();
            await using var context = await NewContext(builder => builder.UseSqlCommenter(), logs.Add);
            SetState();
            await context.TestEntities.TagWith("my comment").ToListAsync();

            var expectedLog =
                @"      -- my comment
      
      SELECT ""t"".""Id"", ""t"".""Name""
      FROM ""TestEntities"" AS ""t""";
            Assert.EndsWith(expectedLog, logs.Last());
        }

        private void SetState()
        {
            SqlCommenterStateHolder.State = new SqlCommenterState
            {
                Route = "my-route",
                ActionName = "my-action",
                AppFramework = "my-framework",
                ControllerName = "my-controller",
            };
        }

        private async Task<TestDbContext> NewContext(Action<DbContextOptionsBuilder<TestDbContext>> builder, Action<string> onLog)
        {
            var optionsBuilder = new DbContextOptionsBuilder<TestDbContext>()
                .UseSqlite(_connection)
                .LogTo(onLog, LogLevel.Information);

            builder(optionsBuilder);

            var context = new TestDbContext(optionsBuilder.Options);
            await context.Database.EnsureCreatedAsync();
            return context;
        }
    }
}
