using System.Linq;
using System.Web;

namespace System.Diagnostics
{
    internal static class ActivityExtensions
    {
        internal static string GetStateString(this Activity activity)
            => string.Join(",", activity.Baggage.Select(x => HttpUtility.UrlEncode(x.Key) + "=" + HttpUtility.UrlEncode(x.Value)));
    }
}
