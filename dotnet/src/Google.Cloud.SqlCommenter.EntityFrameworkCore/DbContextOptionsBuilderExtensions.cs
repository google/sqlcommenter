using System;
using Google.Cloud.SqlCommenter.EntityFrameworkCore;

namespace Microsoft.EntityFrameworkCore
{
    /// <summary>
    /// Contains extensions for configuring sql commenter on an <see cref="DbContextOptionsBuilder"/>.
    /// </summary>
    public static class DbContextOptionsBuilderExtensions
    {
        /// <summary>
        /// Configures the context to add sql statement comments.
        /// </summary>
        /// <seealso href="https://google.github.io/sqlcommenter/spec/"/>
        /// <param name="builder">A builder for setting options on the context.</param>
        /// <param name="optionsAction">An optional action to allow additional configuration.</param>
        /// <returns>
        /// The <see cref="DbContextOptionsBuilder"/> so that further configuration can be chained.
        /// </returns>
        public static DbContextOptionsBuilder UseSqlCommenter(
            this DbContextOptionsBuilder builder,
            Action<SqlCommenterOptions>? optionsAction = null)
        {
            var options = new SqlCommenterOptions();
            optionsAction?.Invoke(options);
            return builder.AddInterceptors(new SqlCommenterInterceptor(options));
        }
    }
}
