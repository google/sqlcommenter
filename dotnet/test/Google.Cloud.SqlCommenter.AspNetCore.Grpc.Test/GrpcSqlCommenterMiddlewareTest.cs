using System.Threading.Tasks;
using Google.Cloud.SqlCommetner.AspnetCore.Grpc.Test;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Xunit;

namespace Google.Cloud.SqlCommenter.AspNetCore.Grpc.Test
{
    public class GrpcSqlCommenterMiddlewareTest
    {
        [Fact]
        public async Task Middlware_ShouldSetState()
        {
            using var host = await StartHost();

            var httpClient = host.GetTestClient();
            var client = new Tester.TesterClient(GrpcChannel.ForAddress(httpClient.BaseAddress!, new GrpcChannelOptions
            {
                HttpClient = httpClient,
            }));

            var state = await client.SayHelloAsync(new HelloRequest());

            Assert.NotNull(state);
            Assert.Equal("/test.Tester/SayHello", state!.Route);
            Assert.Equal("SayHello", state.ActionName);
            Assert.Equal("test.Tester", state.ControllerName);
            Assert.Equal("AspNetCore gRPC", state.AppFramework);
        }

        private Task<IHost> StartHost()
        {
            return new HostBuilder()
                .ConfigureWebHost(webBuilder => webBuilder
                    .UseTestServer()
                    .ConfigureServices(services => services.AddGrpc())
                    .Configure(app => app
                        .UseRouting()
                        .UseGrpcSqlCommenter()
                        .UseEndpoints(endpoints => endpoints.MapGrpcService<TesterService>())))
                .StartAsync();
        }
    }
}
