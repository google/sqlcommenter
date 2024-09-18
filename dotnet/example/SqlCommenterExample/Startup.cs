using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using SqlCommenterExample.Services;

namespace SqlCommenterExample
{
    public class Startup
    {
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers();
            services.AddGrpc();
            services.AddDbContext<SampleDbContext>(o => o
                .UseSqlite("Filename=Sample.sqlite")
                .UseSqlCommenter());
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseRouting();
            app.UseSqlCommenter();
            app.UseGrpcSqlCommenter();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapGet("/", context => context.Response.WriteAsync("Hello World!"));
                endpoints.MapControllers();
                endpoints.MapGrpcService<GreeterService>();
            });
        }
    }
}
