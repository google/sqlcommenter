using Google.Cloud.SqlCommenter.AspNetCore.Middleware;

namespace Microsoft.AspNetCore.Builder
{
    /// <summary>
    /// Contains extensions for configuring sql commenter for asp net core services on an <see cref="IApplicationBuilder"/>.
    /// </summary>
    public static class ApplicationBuilderExtensions
    {
        /// <summary>
        /// Adds a <see cref="SqlCommenterMiddleware"/> middleware to the specified <see cref="IApplicationBuilder"/>.
        /// </summary>
        /// <param name="builder">The <see cref="IApplicationBuilder"/> to add the middleware to.</param>
        /// <returns>A reference to the <see cref="IApplicationBuilder"/>.</returns>
        /// <remarks>
        /// <para>
        /// A call to <see cref="UseSqlCommenter(IApplicationBuilder)"/> must be preceded by a call to
        /// <see cref="EndpointRoutingApplicationBuilderExtensions.UseRouting(IApplicationBuilder)"/> for the same <see cref="IApplicationBuilder"/>.
        /// </para>
        /// </remarks>
        public static IApplicationBuilder UseSqlCommenter(this IApplicationBuilder builder)
        {
            return builder.UseMiddleware<SqlCommenterMiddleware>();
        }
    }
}
