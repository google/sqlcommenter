using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Google.Cloud.SqlCommetner.AspnetCore.Grpc.Test;
using Grpc.Core;

namespace Google.Cloud.SqlCommenter.AspNetCore.Grpc.Test
{
    public class TesterService : Tester.TesterBase
    {
        public override Task<HelloReply> SayHello(HelloRequest request, ServerCallContext context)
        {
            return Task.FromResult(new HelloReply
            {
                Route = SqlCommenterStateHolder.State?.Route,
                ActionName = SqlCommenterStateHolder.State?.ActionName,
                ControllerName = SqlCommenterStateHolder.State?.ControllerName,
                AppFramework = SqlCommenterStateHolder.State?.AppFramework,
            });
        }
    }
}
