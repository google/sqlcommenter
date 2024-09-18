using Microsoft.EntityFrameworkCore;

namespace Google.Cloud.SqlCommenter.EntityFrameworkCore.Test
{
    public class TestDbContext : DbContext
    {
        public DbSet<TestEntity> TestEntities { get; set; } = null!;

        public TestDbContext(DbContextOptions<TestDbContext> options)
            : base(options)
        {
        }
    }
}
