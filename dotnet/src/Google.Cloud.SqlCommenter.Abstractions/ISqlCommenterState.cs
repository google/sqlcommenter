namespace Google.Cloud.SqlCommenter.Abstractions
{
    /// <summary>
    /// The sql commenter values of the current context.
    /// </summary>
    public interface ISqlCommenterState
    {
        /// <summary>
        /// Gets the action name.
        /// </summary>
        string? ActionName { get; }

        /// <summary>
        /// Gets the controller name.
        /// </summary>
        string? ControllerName { get; }

        /// <summary>
        /// Gets the name of the application framework.
        /// </summary>
        string? AppFramework { get; }

        /// <summary>
        /// Gets the route.
        /// </summary>
        string? Route { get; }
    }
}
