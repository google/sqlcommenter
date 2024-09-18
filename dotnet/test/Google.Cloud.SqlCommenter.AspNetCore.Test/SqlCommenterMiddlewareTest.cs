using System.Net.Http.Json;
using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Xunit;

namespace Google.Cloud.SqlCommenter.AspNetCore.Test
{
    public class SqlCommenterMiddlewareTest
    {
        [Fact]
        public async Task Middlware_ShouldSetStateWhenControllerIsCalled()
        {
            using var host = await StartHost();

            var client = host.GetTestClient();
            var state = await client.GetFromJsonAsync<SqlCommenterState>("/test");

            Assert.NotNull(state);
            Assert.Equal("/test", state!.Route);
            Assert.Equal("Get", state.ActionName);
            Assert.Equal("Test", state.ControllerName);
            Assert.Equal("AspNetCore", state.AppFramework);
        }

        [Fact]
        public async Task Middlware_ShouldSetStateWithoutController()
        {
            using var host = await StartHost();

            var client = host.GetTestClient();
            var state = await client.GetFromJsonAsync<SqlCommenterState>("/test-raw");

            Assert.NotNull(state);
            Assert.Equal("/test-raw", state!.Route);
            Assert.Equal("/test-raw HTTP: GET", state.ActionName);
            Assert.Null(state.ControllerName);
            Assert.Equal("AspNetCore", state.AppFramework);
        }

        private Task<IHost> StartHost()
        {
            return new HostBuilder()
                .ConfigureWebHost(webBuilder => webBuilder
                    .UseTestServer()
                    .ConfigureServices(services => services.AddControllers())
                    .Configure(app => app
                        .UseRouting()
                        .UseSqlCommenter()
                        .UseEndpoints(endpoints =>
                        {
                            endpoints.MapControllers();
                            endpoints.MapGet("/test-raw", ctx => ctx.Response.WriteAsJsonAsync(SqlCommenterStateHolder.State));
                        })))
                .StartAsync();
        }
    }
}
