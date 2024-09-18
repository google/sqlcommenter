using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc.Controllers;

namespace Google.Cloud.SqlCommenter.AspNetCore.Middleware
{
    /// <summary>
    /// Sets the <see cref="ISqlCommenterState"/> of the current execution context.
    /// </summary>
    public class SqlCommenterMiddleware
    {
        private const string FrameworkName = "AspNetCore";

        private readonly RequestDelegate _next;

        /// <summary>
        /// Initialize the sql commenter grpc middleware.
        /// </summary>
        /// <param name="next"></param>
        public SqlCommenterMiddleware(RequestDelegate next)
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
            var endpoint = context.GetEndpoint();
            var metadata = endpoint?.Metadata.GetMetadata<ControllerActionDescriptor>();
            SqlCommenterStateHolder.State = new SqlCommenterState
            {
                Route = context.Request.Path,
                ActionName = metadata?.ActionName ?? endpoint?.DisplayName,
                ControllerName = metadata?.ControllerName,
                AppFramework = FrameworkName,
            };

            return _next(context);
        }
    }
}
