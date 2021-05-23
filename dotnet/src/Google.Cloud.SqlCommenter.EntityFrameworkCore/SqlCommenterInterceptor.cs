using System.Data;
using System.Data.Common;
using System.Threading;
using System.Threading.Tasks;
using Google.Cloud.SqlCommenter.Abstractions;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Diagnostics;

namespace Google.Cloud.SqlCommenter.EntityFrameworkCore
{
    /// <summary>
    /// Interceptor which modifies the database commands and adds the sql commenter comments to the command text if needed.
    /// Uses the state of <see cref="SqlCommenterStateHolder"/>.
    /// </summary>
    internal class SqlCommenterInterceptor : DbCommandInterceptor
    {
        private const string DisablePrefix = "-- " + QueryableExtensions.DisableSqlCommenterTag;
        private const string EnablePrefix = "-- " + QueryableExtensions.EnableSqlCommenterTag;
        private const string FrameworkName = "ef-core";

        private readonly SqlCommenterOptions _options;

        public SqlCommenterInterceptor(SqlCommenterOptions options)
        {
            _options = options;
        }

        public override InterceptionResult<DbDataReader> ReaderExecuting(
            DbCommand command,
            CommandEventData eventData,
            InterceptionResult<DbDataReader> result)
        {
            ManipulateCommand(command);
            return base.ReaderExecuting(command, eventData, result);
        }

        public override ValueTask<InterceptionResult<DbDataReader>> ReaderExecutingAsync(
            DbCommand command,
            CommandEventData eventData,
            InterceptionResult<DbDataReader> result,
            CancellationToken cancellationToken = default)
        {
            ManipulateCommand(command);
            return base.ReaderExecutingAsync(command, eventData, result, cancellationToken);
        }

        private void ManipulateCommand(IDbCommand command)
        {
            var enabled = StripEnableStateTags(command);

            if (enabled == false || (!_options.DefaultEnable && enabled != true))
                return;

            // the spec states that if the sql statement has a comment,
            // sql commenter should not add additional comments
            if (HasSqlComment(command.CommandText))
                return;

            var state = SqlCommenterStateHolder.State;
            if (state == null)
                return;

            command.CommandText += " " + state.ToSqlComment(FrameworkName, _options.ExtraValues);
        }

        private static bool? StripEnableStateTags(IDbCommand command)
        {
            if (command.CommandText.StartsWith(DisablePrefix))
            {
                command.CommandText = command.CommandText[DisablePrefix.Length..].TrimStart();
                return false;
            }

            if (command.CommandText.StartsWith(EnablePrefix))
            {
                command.CommandText = command.CommandText[EnablePrefix.Length..].TrimStart();
                return true;
            }

            return null;
        }

        private static bool HasSqlComment(string command) => command.Contains("--") || command.Contains("/*");
    }
}
