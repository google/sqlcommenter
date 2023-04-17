using System.Threading.Tasks;
using Grpc.Core;
using SqlCommenter.Example;

namespace SqlCommenterExample.Services
{
    public class GreeterService : Greeter.GreeterBase
    {
        private readonly SampleDbContext _db;

        public GreeterService(SampleDbContext db)
        {
            _db = db;
        }

        public override async Task SayHello(HelloRequest request, IServerStreamWriter<HelloReply> responseStream, ServerCallContext context)
        {
            await foreach (var contact in _db.Contacts.AsAsyncEnumerable())
            {
                await responseStream.WriteAsync(new HelloReply
                {
                    Message = "Hello " + contact.Name,
                });
            }
        }
    }
}
