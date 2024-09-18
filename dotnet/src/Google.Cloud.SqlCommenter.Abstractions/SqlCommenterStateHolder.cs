using System.Threading;

namespace Google.Cloud.SqlCommenter.Abstractions
{
    /// <summary>
    /// Stores the <see cref="ISqlCommenterState"/> for the current execution context.
    /// <remarks>
    /// Instead of a simple scoped service an async local store is used.
    /// The reason for this is there are database frameworks which use an internal service provider (ex. ef core)
    /// or don't allow to resolve services easily where needed.
    /// Also we don't want to force the user to use the same service provider for the sql framework as for the web framework.
    /// </remarks>
    /// </summary>
    public static class SqlCommenterStateHolder
    {
        private static readonly AsyncLocal<ISqlCommenterState?> SqlCommenterState = new();

        /// <summary>
        /// Gets or sets the <see cref="ISqlCommenterState"/> for the current execution context.
        /// </summary>
        public static ISqlCommenterState? State
        {
            get => SqlCommenterState.Value;
            set => SqlCommenterState.Value = value;
        }
    }
}
