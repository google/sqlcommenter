using System.Diagnostics;
using Xunit;

namespace Google.Cloud.SqlCommenter.Abstractions.Test
{
    public class ActivityExtensionsTest
    {
        [Fact]
        public void GetStateString()
        {
            using var activity = new Activity("testing").Start();
            Assert.Equal(string.Empty, activity.GetStateString());

            activity.AddBaggage("my-baggage", "my-baggage-value");
            activity.AddBaggage("my-baggage-2", "my-baggage-value 2");

            Assert.Equal("my-baggage-2=my-baggage-value+2,my-baggage=my-baggage-value", activity.GetStateString());
        }
    }
}
