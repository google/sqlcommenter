using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Reflection;
using System.Web;

namespace Google.Cloud.SqlCommenter.Abstractions
{
    /// <summary>
    /// Contains extensions for <see cref="ISqlCommenterState"/>.
    /// </summary>
    public static class SqlCommenterStateExtensions
    {
        private static readonly string? ApplicationName = BuildApplicationTagValue();

        /// <summary>
        /// Serializes the <see cref="ISqlCommenterState"/> according to the sql commenter specification.
        /// </summary>
        /// <seealso href="https://google.github.io/sqlcommenter/spec/"/>
        /// <param name="state"></param>
        /// <param name="dbFramework"></param>
        /// <param name="extraValues"></param>
        /// <returns></returns>
        public static string ToSqlComment(
            this ISqlCommenterState state,
            string dbFramework,
            IEnumerable<KeyValuePair<string, string?>>? extraValues = null)
        {
            // an alphabetical order of the keys is required by the spec.
            var values = new SortedDictionary<string, string?>
            {
                { "action", state.ActionName },
                { "application", ApplicationName },
                { "controller", state.ControllerName },
                { "framework", $"{state.AppFramework} {dbFramework}" },
                { "route", state.Route },
            };

            if (Activity.Current is { IdFormat: ActivityIdFormat.W3C })
            {
                values.Add("traceparent", Activity.Current.Id);
                values.Add("tracestate", Activity.Current.GetStateString());
            }

            if (extraValues != null)
            {
                foreach (var (key, value) in extraValues)
                {
                    values[key] = value;
                }
            }

            var keyValuePairs = values
                .Where(x => x.Value != null)
                .Select(x => $"{Encode(x.Key)}='{Encode(x.Value)}'");

            return $"/*{string.Join(",", keyValuePairs)}*/";
        }

        private static string Encode(string? value)
            => HttpUtility.UrlEncode(value ?? string.Empty);

        private static string? BuildApplicationTagValue()
        {
            var assemblyName = Assembly.GetEntryAssembly()?.GetName();
            return assemblyName == null
                ? null
                : $"{assemblyName.Name}@{assemblyName.Version}";
        }
    }
}
