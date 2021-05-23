using Google.Cloud.SqlCommenter.Abstractions;
using Microsoft.AspNetCore.Mvc;

namespace Google.Cloud.SqlCommenter.AspNetCore.Test
{
    [ApiController]
    [Route("test")]
    public class TestController
    {
        [HttpGet]
        public ISqlCommenterState? Get() => SqlCommenterStateHolder.State;
    }
}
