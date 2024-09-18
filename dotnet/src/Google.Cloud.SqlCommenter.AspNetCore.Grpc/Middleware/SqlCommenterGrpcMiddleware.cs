using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Grpc.AspNetCore.Server;
using Microsoft.AspNetCore.Http;

namespace Google.Cloud.SqlCommenter.AspNetCore.Grpc.Middleware
{
    /// <summary>
    /// Sets the <see cref="ISqlCommenterState"/> of the current execution context for grpc methods.
    /// </summary>
    public class SqlCommenterGrpcMiddleware
    {
        private const string FrameworkName = "AspNetCore gRPC";

        private readonly RequestDelegate _next;

        /// <summary>
        /// Initialize the sql commenter grpc middleware.
        /// </summary>
        /// <param name="next"></param>
        public SqlCommenterGrpcMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        /// <summary>
        /// Invoke the middleware.
        /// </summary>
        /// <param name="context">The <see cref="HttpContext"/>.</param>
        /// <returns></returns>
        public Task InvokeAsync(HttpContext context)
        {
            var methodDescriptor = context.GetEndpoint()?.Metadata.GetMetadata<GrpcMethodMetadata>();
            if (methodDescriptor != null)
            {
                SqlCommenterStateHolder.State = new SqlCommenterState
                {
                    Route = context.Request.Path,
                    ActionName = methodDescriptor.Method.Name,
                    ControllerName = methodDescriptor.Method.ServiceName,
                    AppFramework = FrameworkName,
                };
            }

            return _next(context);
        }
    }
}
