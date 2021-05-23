using Microsoft.EntityFrameworkCore;
using SqlCommenterExample.Models;

namespace SqlCommenterExample
{
    public class SampleDbContext : DbContext
    {
        public DbSet<Contact> Contacts { get; set; } = null!;

        public SampleDbContext(DbContextOptions<SampleDbContext> options)
            : base(options)
        {
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<Contact>()
                .HasData(
                    new Contact { Id = 100_000, Name = "Kristin" },
                    new Contact { Id = 100_001, Name = "Bill" },
                    new Contact { Id = 100_002, Name = "Nicole" },
                    new Contact { Id = 100_003, Name = "James" }
                );
        }
    }
}
