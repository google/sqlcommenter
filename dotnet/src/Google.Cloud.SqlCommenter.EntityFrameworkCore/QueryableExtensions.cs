using System.Linq;
using Google.Cloud.SqlCommenter.EntityFrameworkCore;

namespace Microsoft.EntityFrameworkCore
{
    /// <summary>
    /// Contains extensions for configuring sql commenter on queries.
    /// </summary>
    public static class QueryableExtensions
    {
        /// <summary>
        /// The query tag used to disable the sql commenter.
        /// </summary>
        internal const string DisableSqlCommenterTag = "SQL-Commenter-Enabled:false";

        /// <summary>
        /// The query tag used to enable the sql commenter.
        /// </summary>
        internal const string EnableSqlCommenterTag = "SQL-Commenter-Enabled:true";

        /// <summary>
        /// Disables the sql commenter on this query.
        /// Only takes effect if <see cref="SqlCommenterOptions.DefaultEnable"/> is set to <c>true</c>.
        /// </summary>
        /// <param name="queryable">The queryable.</param>
        /// <typeparam name="T">The type of the entity.</typeparam>
        /// <returns>An updated <see cref="IQueryable{T}"/> instance.</returns>
        public static IQueryable<T> DisableSqlCommenter<T>(this IQueryable<T> queryable) => queryable.TagWith(DisableSqlCommenterTag);

        /// <summary>
        /// Enables the sql commenter on this query.
        /// Only takes effect if <see cref="SqlCommenterOptions.DefaultEnable"/> is set to <c>false</c>.
        /// </summary>
        /// <param name="queryable">The queryable.</param>
        /// <typeparam name="T">The type of the entity.</typeparam>
        /// <returns>An updated <see cref="IQueryable{T}"/> instance.</returns>
        public static IQueryable<T> EnableSqlCommenter<T>(this IQueryable<T> queryable) => queryable.TagWith(EnableSqlCommenterTag);
    }
}
