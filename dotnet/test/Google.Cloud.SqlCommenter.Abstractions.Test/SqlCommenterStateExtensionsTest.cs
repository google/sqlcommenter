using System.Collections.Generic;
using System.Text.RegularExpressions;
using Xunit;

namespace Google.Cloud.SqlCommenter.Abstractions.Test
{
    public class SqlCommenterStateExtensionsTest
    {
        [Fact]
        public void ToSqlComment()
        {
            var state = new SqlCommenterState
            {
                Route = "my route",
                ActionName = "my action",
                AppFramework = "AspNetCore",
                ControllerName = "my+controller",
            };

            AssertSqlComment("/*action='my+action',controller='my%2bcontroller',framework='AspNetCore+ef-core',route='my+route'*/", state,
                "ef-core");
        }

        [Fact]
        public void ToSqlCommentWithExtraValues()
        {
            var state = new SqlCommenterState
            {
                Route = "my route",
                ActionName = "my action",
                AppFramework = "AspNetCore",
                ControllerName = "my+controller",
            };

            var extraValues = new Dictionary<string, string?>
            {
                { "application", "my-app" },
                { "my-key", "my-value" },
                { "my-key-null", null },
            };

            AssertSqlComment(
                "/*action='my+action',application='my-app',controller='my%2bcontroller',framework='AspNetCore+ef-core',my-key='my-value',route='my+route'*/",
                state,
                "ef-core",
                extraValues,
                false);
        }

        private static void AssertSqlComment(
            string expectedComment,
            ISqlCommenterState state,
            string dbFramework,
            IEnumerable<KeyValuePair<string, string?>>? extraValues = null,
            bool removeApplicationTag = true)
        {
            var sqlComment = state.ToSqlComment(dbFramework, extraValues);

            if (removeApplicationTag)
            {
                // remove the application tag since it is dependent on the test launcher.
                sqlComment = Regex.Replace(sqlComment, "application='.*?',", string.Empty);
            }

            Assert.Equal(expectedComment, sqlComment);
        }
    }
}
