namespace Google.Cloud.SqlCommenter.Abstractions
{
    /// <summary>
    /// A very simple implementation of <see cref="ISqlCommenterState"/>.
    /// </summary>
    public class SqlCommenterState : ISqlCommenterState
    {
        /// <inheritdoc />
        public string? ActionName { get; set; }

        /// <inheritdoc />
        public string? ControllerName { get; set; }

        /// <inheritdoc />
        public string? AppFramework { get; set; }

        /// <inheritdoc />
        public string? Route { get; set; }
    }
}
