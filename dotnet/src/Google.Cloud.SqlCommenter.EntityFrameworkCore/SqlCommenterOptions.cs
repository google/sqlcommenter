using System.Collections.Generic;

namespace Google.Cloud.SqlCommenter.EntityFrameworkCore
{
    /// <summary>
    /// Options to customize the sql commenter.
    /// </summary>
    public class SqlCommenterOptions
    {
        /// <summary>
        /// Gets or sets whether the sql commenter is enabled by default.
        /// </summary>
        public bool DefaultEnable { get; set; } = true;

        /// <summary>
        /// Adds additional key value pairs or overwrites existing key value pairs of the sql commenter.
        /// </summary>
        public IDictionary<string, string?> ExtraValues { get; } = new Dictionary<string, string?>();
    }
}
