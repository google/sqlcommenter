using Google.Cloud.SqlCommenter.AspNetCore.Grpc.Middleware;

namespace Microsoft.AspNetCore.Builder
{
    /// <summary>
    /// Contains extensions for configuring sql commenter for grpc services on an <see cref="IApplicationBuilder"/>.
    /// </summary>
    public static class ApplicationBuilderExtensions
    {
        /// <summary>
        /// Adds a <see cref="SqlCommenterGrpcMiddleware"/> middleware to the specified <see cref="IApplicationBuilder"/>.
        /// </summary>
        /// <param name="builder">The <see cref="IApplicationBuilder"/> to add the middleware to.</param>
        /// <returns>A reference to the <see cref="IApplicationBuilder"/>.</returns>
        /// <remarks>
        /// <para>
        /// A call to <see cref="UseGrpcSqlCommenter(IApplicationBuilder)"/> must be preceded by a call to
        /// <see cref="EndpointRoutingApplicationBuilderExtensions.UseRouting(IApplicationBuilder)"/> for the same <see cref="IApplicationBuilder"/>.
        /// </para>
        /// </remarks>
        public static IApplicationBuilder UseGrpcSqlCommenter(this IApplicationBuilder builder)
        {
            return builder.UseMiddleware<SqlCommenterGrpcMiddleware>();
        }
    }
}
